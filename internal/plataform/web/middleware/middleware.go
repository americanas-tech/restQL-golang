package middleware

import (
	"github.com/valyala/fasthttp"
)

type Middleware interface {
	Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler
}

type NoopMiddleware struct{}

func (nm NoopMiddleware) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return h
}

func Apply(h fasthttp.RequestHandler, mws ...Middleware) fasthttp.RequestHandler {
	handler := h
	for _, m := range mws {
		handler = m.Apply(handler)
	}

	return handler
}
