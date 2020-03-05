package middleware

import (
	"context"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

type Timeout struct {
	duration time.Duration
}

func NewTimeout(duration string) Middleware {
	d, parseErr := time.ParseDuration(duration)
	if parseErr != nil {
		log.Printf("[WARN] failed to initialize timeout middleware : invalid duration : %s", duration)
		return NoopMiddleware{}
	}

	return Timeout{d}
}

func (t Timeout) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nativeCtx := GetNativeContext(ctx)
		timeout, cancel := context.WithTimeout(nativeCtx, t.duration)
		defer cancel()

		WithNativeContext(ctx, timeout)

		h(ctx)
	}
}
