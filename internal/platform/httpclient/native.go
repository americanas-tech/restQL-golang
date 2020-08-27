package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/plugins"
	"github.com/pkg/errors"
	"github.com/rs/dnscache"
)

const defaultStatusCode = 0

type nativeHTTPClient struct {
	client    *http.Client
	log       restql.Logger
	lifecycle plugins.Lifecycle
}

func newNativeHTTPClient(log restql.Logger, l plugins.Lifecycle, cfg *conf.Config) *nativeHTTPClient {
	clientCfg := cfg.HTTP.Client

	r := &dnscache.Resolver{}
	go func() {
		t := time.NewTicker(10 * time.Minute)
		defer t.Stop()
		for range t.C {
			r.Refresh(true)
		}
	}()

	dialer := net.Dialer{
		Timeout: clientCfg.ConnTimeout,
	}

	t := &http.Transport{
		MaxIdleConns:        clientCfg.MaxIdleConns,
		MaxIdleConnsPerHost: clientCfg.MaxIdleConnsPerHost,
		MaxConnsPerHost:     clientCfg.MaxConnsPerHost,
		IdleConnTimeout:     clientCfg.MaxIdleConnDuration,
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			ips, err := r.LookupHost(ctx, host)
			if err != nil {
				return nil, err
			}
			for _, ip := range ips {
				conn, err = dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
				if err == nil {
					break
				}
			}
			return
		},
	}

	c := &http.Client{
		Timeout:   clientCfg.MaxRequestTimeout,
		Transport: t,
	}

	return &nativeHTTPClient{
		client:    c,
		log:       log,
		lifecycle: l,
	}
}

func (nc *nativeHTTPClient) Do(ctx context.Context, request domain.HTTPRequest) (domain.HTTPResponse, error) {
	ctx = nc.lifecycle.BeforeRequest(ctx, request)
	log := restql.GetLogger(ctx)

	req, err := nc.makeRequest(request)
	if err != nil {
		return domain.HTTPResponse{}, err
	}
	requestURL := req.URL.String()
	target := req.URL.Host

	timeout, cancel := context.WithTimeout(ctx, request.Timeout)
	defer cancel()

	req = req.WithContext(timeout)

	log.Debug("request created", "request-url", req.URL.String())

	start := time.Now()
	response, err := nc.client.Do(req)
	duration := time.Since(start)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			errorResponse := makeErrorResponse(requestURL, duration, http.StatusRequestTimeout)
			log.Warn("request timed out", "url", requestURL, "target", target, "method", request.Method, "duration-ms", duration.Milliseconds())

			nc.lifecycle.AfterRequest(ctx, request, errorResponse, domain.ErrRequestTimeout)

			return errorResponse, domain.ErrRequestTimeout
		}

		errorResponse := makeErrorResponse(requestURL, duration, defaultStatusCode)
		log.Error("request finished with error", err, "url", requestURL, "target", target, "method", request.Method, "duration-ms", duration.Milliseconds())

		nc.lifecycle.AfterRequest(ctx, request, errorResponse, err)
		return errorResponse, err
	}

	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			log.Error("failed to close response body", closeErr)
		}
	}()

	body, err := nc.unmarshalBody(log, response)
	if err != nil {
		errorResponse := makeErrorResponse(requestURL, duration, defaultStatusCode)
		nc.lifecycle.AfterRequest(ctx, request, errorResponse, err)

		return errorResponse, err
	}

	hr := make(map[string]string)
	for k, s := range response.Header {
		hr[k] = s[0]
	}

	httpResponse := domain.HTTPResponse{
		URL:        requestURL,
		StatusCode: response.StatusCode,
		Body:       body,
		Headers:    hr,
		Duration:   duration,
	}

	nc.lifecycle.AfterRequest(ctx, request, httpResponse, err)

	return httpResponse, nil
}

func (nc *nativeHTTPClient) makeRequest(request domain.HTTPRequest) (*http.Request, error) {
	req := http.Request{
		Method: request.Method,
		URL:    makeURL(request),
	}

	if request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodPatch {
		body, err := nc.makeBody(request)
		if err != nil {
			return nil, err
		}

		req.Body = body
	}

	if len(request.Headers) > 0 {
		req.Header = make(http.Header)
		for k, v := range request.Headers {
			req.Header.Add(k, v)
		}
	}

	return &req, nil
}

func (nc *nativeHTTPClient) makeBody(request domain.HTTPRequest) (io.ReadCloser, error) {
	body := request.Body

	if body, ok := body.(string); ok {
		return ioutil.NopCloser(strings.NewReader(body)), nil
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request body")
	}

	r := ioutil.NopCloser(bytes.NewReader(data))
	return r, nil
}

func (nc *nativeHTTPClient) unmarshalBody(log restql.Logger, response *http.Response) (interface{}, error) {
	target := response.Request.URL.Host
	requestURL := response.Request.URL.String()
	statusCode := response.StatusCode

	bodyByte, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Error("failed to read response body", readErr, "url", requestURL, "target", target, "statusCode", statusCode)

		return nil, readErr
	}

	if len(bodyByte) == 0 || !json.Valid(bodyByte) {
		body := string(bodyByte)
		err := errors.New("invalid json")

		log.Error("invalid json as body", err, "body", body, "url", requestURL, "target", target, "statusCode", statusCode)

		return body, nil
	}

	var responseBody interface{}
	err := json.Unmarshal(bodyByte, &responseBody)
	if err != nil {
		body := string(bodyByte)
		log.Error("failed to unmarshal response body", err, "body", body, "url", requestURL, "target", target, "statusCode", statusCode)

		return body, nil
	}

	return responseBody, nil
}

func makeURL(request domain.HTTPRequest) *url.URL {
	u := &url.URL{
		Host:   request.Host,
		Scheme: request.Schema,
		Path:   request.Path,
	}
	queryParams := u.Query()

	for k, v := range request.Query {
		switch v := v.(type) {
		case []interface{}:
			for _, vv := range v {
				parsed := parseQueryValue(vv)
				if parsed != "" {
					queryParams.Add(k, parsed)
				}
			}
		default:
			parsed := parseQueryValue(v)
			if parsed != "" {
				queryParams.Add(k, parsed)
			}
		}
	}

	u.RawQuery = queryParams.Encode()

	return u
}

func parseQueryValue(value interface{}) string {
	switch value := value.(type) {
	case string:
		return value
	case bool:
		return strconv.FormatBool(value)
	case int:
		return strconv.Itoa(value)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case map[string]interface{}:
		return parseMapParam(value)
	default:
		return ""
	}
}

func parseMapParam(value map[string]interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return string(data)
}

func makeErrorResponse(requestURL string, responseTime time.Duration, statusCode int) domain.HTTPResponse {
	return domain.HTTPResponse{
		URL:        requestURL,
		StatusCode: statusCode,
		Duration:   responseTime,
	}
}
