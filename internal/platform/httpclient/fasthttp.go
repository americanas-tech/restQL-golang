package httpclient

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/internal/platform/plugins"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"sync"
	"time"
)

const clientPoolSize = 80

type httpResult struct {
	target   string
	duration time.Duration
	err      error
	response *fasthttp.Response
}

type fastHttpClient struct {
	clientPool    *clientPool
	log           *logger.Logger
	pluginManager plugins.Manager
	responsePool  *sync.Pool
}

func newFastHttpClient(log *logger.Logger, pm plugins.Manager, cfg *conf.Config) *fastHttpClient {
	rp := &sync.Pool{
		New: func() interface{} {
			return make(chan httpResult)
		},
	}

	pool := newClientPool(cfg)

	return &fastHttpClient{clientPool: pool, log: log, pluginManager: pm, responsePool: rp}
}

func (hc *fastHttpClient) Do(ctx context.Context, request domain.HttpRequest) (domain.HttpResponse, error) {
	requestCtx := hc.pluginManager.RunBeforeRequest(ctx, request)

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

		client := hc.clientPool.Get(request)

		res := fasthttp.AcquireResponse()
		start := time.Now()
		err = client.DoTimeout(req, res, request.Timeout)
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
		response := makeErrorResponse(hr.target, hr.duration, http.StatusRequestTimeout)

		fasthttp.ReleaseResponse(hr.response)

		hc.pluginManager.RunAfterRequest(requestCtx, request, response, hr.err)

		return response, domain.ErrRequestTimeout
	case hr.err != nil:
		response := makeErrorResponse(hr.target, hr.duration, 0)

		if hr.response != nil {
			fasthttp.ReleaseResponse(hr.response)
		}

		hc.pluginManager.RunAfterRequest(requestCtx, request, response, hr.err)

		return response, errors.Wrap(hr.err, "request execution failed")
	}

	response, err := makeResponse(hr.target, hr.response, hr.duration)
	if err != nil {
		response = makeErrorResponse(hr.target, hr.duration, 0)

		fasthttp.ReleaseResponse(hr.response)

		hc.pluginManager.RunAfterRequest(requestCtx, request, response, err)

		return response, err
	}

	fasthttp.ReleaseResponse(hr.response)

	hc.pluginManager.RunAfterRequest(requestCtx, request, response, err)

	return response, nil
}
