package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/internal/platform/plugins"
	"github.com/pkg/errors"
	"github.com/rs/dnscache"
)

const poolSize = 20
const defaultStatusCode = 0

type nativeHttpClient struct {
	clients       []*http.Client
	log           *logger.Logger
	pluginManager plugins.Manager
}

func newNativeHttpClient(log *logger.Logger, pm plugins.Manager, cfg *conf.Config) *nativeHttpClient {
	clientCfg := cfg.Http.Client

	r := &dnscache.Resolver{}
	go func() {
		t := time.NewTicker(10 * time.Minute)
		defer t.Stop()
		for range t.C {
			r.Refresh(true)
		}
	}()

	clients := make([]*http.Client, poolSize)

	for i := 0; i < poolSize; i++ {
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

		clients[i] = c
	}

	return &nativeHttpClient{
		clients:       clients,
		log:           log,
		pluginManager: pm,
	}
}

func (nc *nativeHttpClient) Do(ctx context.Context, request domain.HttpRequest) (domain.HttpResponse, error) {
	ctx = nc.pluginManager.RunBeforeRequest(ctx, request)
	log := restql.GetLogger(ctx)

	req, err := nc.makeRequest(request)
	if err != nil {
		return domain.HttpResponse{}, err
	}
	requestUrl := req.URL.String()

	timeout, cancel := context.WithTimeout(ctx, request.Timeout)
	defer cancel()

	req = req.WithContext(timeout)

	log.Debug("request created", "request-url", req.URL.String())

	client := nc.peekClient()

	start := time.Now()
	response, err := client.Do(req)
	duration := time.Since(start)
	if err != nil {
		errorResponse := makeErrorResponse(requestUrl, duration, http.StatusRequestTimeout)

		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Info("request timed out", "url", requestUrl, "method", request.Method, "duration-ms", duration.Milliseconds())
			nc.pluginManager.RunAfterRequest(ctx, request, errorResponse, domain.ErrRequestTimeout)

			return errorResponse, domain.ErrRequestTimeout
		}

		nc.pluginManager.RunAfterRequest(ctx, request, errorResponse, err)
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
		errorResponse := makeErrorResponse(requestUrl, duration, defaultStatusCode)
		nc.pluginManager.RunAfterRequest(ctx, request, errorResponse, err)

		return errorResponse, fmt.Errorf("%w: %v", domain.ErrInvalidResponseBody, body)
	}

	hr := make(map[string]string)
	for k, s := range response.Header {
		hr[k] = s[0]
	}

	httpResponse := domain.HttpResponse{
		Url:        requestUrl,
		StatusCode: response.StatusCode,
		Body:       body,
		Headers:    hr,
		Duration:   duration,
	}

	nc.pluginManager.RunAfterRequest(ctx, request, httpResponse, err)

	return httpResponse, nil
}

func (nc *nativeHttpClient) peekClient() *http.Client {
	n := rand.Intn(poolSize)
	return nc.clients[n]
}

func (nc *nativeHttpClient) makeRequest(request domain.HttpRequest) (*http.Request, error) {
	req := http.Request{
		Method: request.Method,
		URL:    makeUrl(request),
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

func (nc *nativeHttpClient) makeBody(request domain.HttpRequest) (io.ReadCloser, error) {
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

func (nc *nativeHttpClient) unmarshalBody(log restql.Logger, response *http.Response) (interface{}, error) {
	var responseBody interface{}
	decoder := json.NewDecoder(response.Body)
	if !decoder.More() {
		return nil, nil
	}

	err := decoder.Decode(&responseBody)
	if err != nil {
		bb, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			log.Error("failed to read response body", readErr, "target", response.Request.Host)
			return nil, readErr
		}
		body := string(bb)

		log.Error("failed to unmarshal response body", err, "body", body, "target", response.Request.Host)

		return body, err
	}

	return responseBody, nil
}

func makeUrl(request domain.HttpRequest) *url.URL {
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

func makeErrorResponse(requestUrl string, responseTime time.Duration, statusCode int) domain.HttpResponse {
	return domain.HttpResponse{
		Url:        requestUrl,
		StatusCode: statusCode,
		Duration:   responseTime,
	}
}
