package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/internal/platform/plugins"
	"github.com/pkg/errors"
	"github.com/rs/dnscache"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type nativeHttpClient struct {
	client        *http.Client
	log           *logger.Logger
	pluginManager plugins.Manager
}

func newNativeHttpClient(log *logger.Logger, pm plugins.Manager, cfg *conf.Config) *nativeHttpClient {
	clientCfg := cfg.Web.Client

	r := &dnscache.Resolver{}
	go func() {
		t := time.NewTicker(10 * time.Minute)
		defer t.Stop()
		for range t.C {
			r.Refresh(true)
		}
	}()

	t := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxConnsPerHost:       clientCfg.MaxConnsPerHost,
		MaxIdleConnsPerHost:   1024,
		IdleConnTimeout:       clientCfg.MaxIdleConnDuration,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			ips, err := r.LookupHost(context.Background(), host)
			if err != nil {
				return nil, err
			}
			for _, ip := range ips {
				var dialer net.Dialer
				conn, err = dialer.Dial(network, net.JoinHostPort(ip, port))
				if err == nil {
					break
				}
			}
			return
		},
	}

	c := &http.Client{
		Timeout:   clientCfg.ReadTimeout,
		Transport: t,
	}

	return &nativeHttpClient{
		client:        c,
		log:           log,
		pluginManager: pm,
	}
}

func (nc *nativeHttpClient) Do(ctx context.Context, request domain.HttpRequest) (domain.HttpResponse, error) {
	req, err := nc.makeRequest(request)
	if err != nil {
		return domain.HttpResponse{}, err
	}
	requestUrl := req.URL.String()

	timeout, cancel := context.WithTimeout(ctx, request.Timeout)
	defer cancel()

	req = req.WithContext(timeout)

	nc.log.Debug("request created", "request-url", req.URL.String())

	start := time.Now()
	response, err := nc.client.Do(req)
	duration := time.Since(start)
	if err != nil {
		errorResponse := makeErrorResponse(requestUrl, duration, err)

		if err, ok := err.(net.Error); ok && err.Timeout() {
			nc.log.Info("request timed out", "url", requestUrl, "method", request.Method, "duration-ms", duration.Milliseconds())
			return errorResponse, domain.ErrRequestTimeout
		}

		return errorResponse, err
	}

	defer func() {
		closeErr := response.Body.Close()
		if err != nil {
			nc.log.Error("failed to close response body", closeErr)
		}
	}()

	body, err := unmarshal(response)
	if err != nil {
		return makeErrorResponse(requestUrl, duration, err), err
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

	return httpResponse, nil
}

func (nc *nativeHttpClient) makeRequest(request domain.HttpRequest) (*http.Request, error) {
	req := http.Request{
		Method: request.Method,
		URL:    makeUrl(request),
	}

	if request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodPatch {
		data, err := json.Marshal(request.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request body")
		}

		req.Body = ioutil.NopCloser(bytes.NewReader(data))
	}

	if len(request.Headers) > 0 {
		req.Header = make(http.Header)
		for k, v := range request.Headers {
			req.Header.Add(k, v)
		}
	}

	return &req, nil
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
		return url.QueryEscape(value)
	case int:
		return url.QueryEscape(strconv.Itoa(value))
	case float64:
		return url.QueryEscape(strconv.FormatFloat(value, 'f', -1, 64))
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

	return url.QueryEscape(string(data))
}

func unmarshal(response *http.Response) (interface{}, error) {
	var responseBody interface{}
	err := json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal response. error: %v", err)
	}
	return responseBody, nil
}
