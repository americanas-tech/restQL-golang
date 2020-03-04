package middleware

import (
	"context"
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

	return contextMiddleware(handler)
}

func contextMiddleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		WithNativeContext(ctx, context.Background())

		h(ctx)
	}
}
