package httpclient

import (
	"context"
	"fmt"
	"github.com/b2wdigital/restQL-golang/v6/internal/domain"
	"github.com/b2wdigital/restQL-golang/v6/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v6/internal/platform/plugins"
	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"github.com/pkg/errors"
	"github.com/rs/dnscache"
	"github.com/valyala/fasthttp"
	"net"
	"sync"
	"time"
)

var errInvalidJSON = errors.New("invalid json")

type httpResult struct {
	target   string
	duration time.Duration
	err      error
	response *fasthttp.Response
}

type fastHTTPClient struct {
	client       *fasthttp.Client
	log          restql.Logger
	lifecycle    plugins.Lifecycle
	responsePool *sync.Pool
}

func newFastHTTPClient(log restql.Logger, pm plugins.Lifecycle, cfg *conf.Config) *fastHTTPClient {
	clientCfg := cfg.HTTP.Client

	r := &dnscache.Resolver{}
	go func() {
		t := time.NewTicker(clientCfg.DnsRefreshInterval)
		defer t.Stop()
		for range t.C {
			r.Refresh(true)
		}
	}()
	dialer := &fasthttp.TCPDialer{
		Resolver: &net.Resolver{
			PreferGo:     true,
			StrictErrors: false,
			Dial: func(ctx context.Context, network, address string) (conn net.Conn, err error) {
				host, port, err := net.SplitHostPort(address)
				if err != nil {
					return nil, err
				}
				ips, err := r.LookupHost(ctx, host)
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
		},
	}

	rp := &sync.Pool{
		New: func() interface{} {
			return make(chan httpResult)
		},
	}

	c := &fasthttp.Client{
		Name:                          "restql",
		NoDefaultUserAgentHeader:      false,
		DisableHeaderNamesNormalizing: true,
		Dial:                          dialer.Dial,
		MaxConnsPerHost:               clientCfg.MaxConnsPerHost,
		MaxIdleConnDuration:           clientCfg.MaxIdleConnDuration,
		MaxConnWaitTimeout:            clientCfg.ConnTimeout,
	}

	return &fastHTTPClient{client: c, log: log, lifecycle: pm, responsePool: rp}
}

func (hc *fastHTTPClient) Do(ctx context.Context, request restql.HTTPRequest) (restql.HTTPResponse, error) {
	requestCtx := hc.lifecycle.BeforeRequest(ctx, request)

	c := hc.responsePool.Get().(chan httpResult)

	go func() {
		req := fasthttp.AcquireRequest()

		err := setupRequest(request, req)
		if err != nil {
			hc.log.Error("failed to setup http client request", err)
			fasthttp.ReleaseRequest(req)
			c <- httpResult{target: request.Host, err: err, duration: 0}
			return
		}

		res := fasthttp.AcquireResponse()
		start := time.Now()
		err = hc.client.DoTimeout(req, res, request.Timeout)
		finish := time.Since(start)

		reqUri := req.URI().String()
		fasthttp.ReleaseRequest(req)

		c <- httpResult{target: reqUri, err: err, duration: finish, response: res}
	}()

	hr := <-c
	hc.responsePool.Put(c)

	switch {
	case hr.err == fasthttp.ErrTimeout || hr.err == fasthttp.ErrDialTimeout || hr.err == fasthttp.ErrTLSHandshakeTimeout:
		hc.log.Info("request timed out", "url", hr.target, "method", request.Method, "duration-ms", hr.duration.Milliseconds())
		response := makeErrorResponse(hr.target, hr.duration, fasthttp.StatusRequestTimeout)

		fasthttp.ReleaseResponse(hr.response)

		err := fmt.Errorf("%w: %s", hr.err, domain.ErrRequestTimeout)
		hc.lifecycle.AfterRequest(requestCtx, request, response, err)

		return response, domain.ErrRequestTimeout
	case hr.err != nil:
		response := makeErrorResponse(hr.target, hr.duration, hr.response.StatusCode())

		if hr.response != nil {
			fasthttp.ReleaseResponse(hr.response)
		}

		hc.lifecycle.AfterRequest(requestCtx, request, response, hr.err)

		return response, errors.Wrap(hr.err, "request execution failed")
	}

	body, err := unmarshalBody(hc.log, hr.response)
	if err != nil {
		hc.log.Error("invalid json as body", err, "url", hr.target, "body", body.Unmarshal(), "statusCode", hr.response.StatusCode())
	}

	response := restql.HTTPResponse{
		URL:        hr.target,
		StatusCode: hr.response.StatusCode(),
		Headers:    readHeaders(hr.response),
		Duration:   hr.duration,
		Body:       body,
	}

	fasthttp.ReleaseResponse(hr.response)

	hc.lifecycle.AfterRequest(requestCtx, request, response, hr.err)

	return response, nil
}
