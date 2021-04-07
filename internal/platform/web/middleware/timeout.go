package middleware

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"time"

	"github.com/valyala/fasthttp"
)

type timeout struct {
	duration time.Duration
}

func newTimeout(duration string, log restql.Logger) Middleware {
	d, parseErr := time.ParseDuration(duration)
	if parseErr != nil {
		log.Warn("failed to initialize timeout middleware : invalid duration", "duration", duration)
		return noopMiddleware{}
	}

	return timeout{d}
}

func (t timeout) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nativeCtx := GetNativeContext(ctx)
		timeout, cancel := context.WithTimeout(nativeCtx, t.duration)
		defer cancel()

		WithNativeContext(ctx, timeout)

		h(ctx)
	}
}
