package middleware

import (
	"context"
	"github.com/valyala/fasthttp"
	"time"
)

type Timeout struct {
	duration time.Duration
}

func NewTimeout(duration string) Middleware {
	d, parseErr := time.ParseDuration(duration)
	if parseErr != nil {
		return NoopMiddleware{}
	}

	return Timeout{d}
}

func (t Timeout) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nativeCtx := NativeContext(ctx)
		timeout, cancel := context.WithTimeout(nativeCtx, t.duration)
		defer cancel()

		WithNativeContext(ctx, timeout)

		h(ctx)
	}
}
