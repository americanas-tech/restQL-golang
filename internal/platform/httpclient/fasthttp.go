package httpclient

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/plugins"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
	"github.com/rs/dnscache"
	"github.com/valyala/fasthttp"
	"net"
	"sync"
	"time"
)

type httpResult struct {
	target   string
	duration time.Duration
	err      error
	response *fasthttp.Response
}

type fastHttpClient struct {
	client        *fasthttp.Client
	log           restql.Logger
	pluginManager plugins.Lifecycle
	responsePool  *sync.Pool
}

func newFastHttpClient(log restql.Logger, pm plugins.Lifecycle, cfg *conf.Config) *fastHttpClient {
	clientCfg := cfg.HTTP.Client

	r := &dnscache.Resolver{}
	go func() {
		t := time.NewTicker(10 * time.Minute)
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
		//ReadTimeout:                   clientCfg.ReadTimeout,
		//WriteTimeout:                  clientCfg.WriteTimeout,
		MaxConnsPerHost:     clientCfg.MaxConnsPerHost,
		MaxIdleConnDuration: clientCfg.MaxIdleConnDuration,
		//MaxConnDuration:               clientCfg.MaxConnDuration,
		MaxConnWaitTimeout: clientCfg.ConnTimeout,
	}

	return &fastHttpClient{client: c, log: log, pluginManager: pm, responsePool: rp}
}

func (hc *fastHttpClient) Do(ctx context.Context, request restql.HTTPRequest) (restql.HTTPResponse, error) {
	requestCtx := hc.pluginManager.BeforeRequest(ctx, request)

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
	case hr.err == fasthttp.ErrTimeout:
		hc.log.Info("request timed out", "url", hr.target, "method", request.Method, "duration-ms", hr.duration.Milliseconds())
		response := makeErrorResponse(hr.target, hr.duration, fasthttp.StatusRequestTimeout)

		fasthttp.ReleaseResponse(hr.response)

		hc.pluginManager.AfterRequest(requestCtx, request, response, hr.err)

		return response, domain.ErrRequestTimeout
	case hr.err != nil:
		response := makeErrorResponse(hr.target, hr.duration, hr.response.StatusCode())

		if hr.response != nil {
			fasthttp.ReleaseResponse(hr.response)
		}

		hc.pluginManager.AfterRequest(requestCtx, request, response, hr.err)

		return response, errors.Wrap(hr.err, "request execution failed")
	}

	response := restql.HTTPResponse{
		URL:        hr.target,
		StatusCode: hr.response.StatusCode(),
		Headers:    readHeaders(hr.response),
		Duration:   hr.duration,
		Body:       hc.unmarshalBody(hc.log, hr.response),
	}

	fasthttp.ReleaseResponse(hr.response)

	hc.pluginManager.AfterRequest(requestCtx, request, response, hr.err)

	return response, nil
}

func (hc *fastHttpClient) unmarshalBody(log restql.Logger, response *fasthttp.Response) *restql.ResponseBody {
	//target := response.Request.URL.Host
	//requestURL := response.Header.
	statusCode := response.StatusCode
	//
	////response := restql.HTTPResponse{
	////	URL:        requestURL,
	////	StatusCode: res.StatusCode(),
	////	Headers:    readHeaders(res),
	////	Duration:   responseTime,
	////}

	bodyByte := response.Body()
	bb := make([]byte, len(bodyByte))
	copy(bb, bodyByte)

	rb := restql.NewResponseBodyFromBytes(log, bb)
	if !rb.Valid() {
		log.Error("invalid json as body", errInvalidJson, "body", rb.Unmarshal(), "statusCode", statusCode)
	}

	return rb
}

func readHeaders(res *fasthttp.Response) restql.Headers {
	h := make(restql.Headers)
	res.Header.VisitAll(func(key, value []byte) {
		h[string(key)] = string(value)
	})

	return h
}
