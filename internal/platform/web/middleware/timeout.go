package middleware

import (
	"context"
	"time"

	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/valyala/fasthttp"
)

type Timeout struct {
	duration time.Duration
}

func NewTimeout(duration string, log *logger.Logger) Middleware {
	d, parseErr := time.ParseDuration(duration)
	if parseErr != nil {
		log.Warn("failed to initialize timeout middleware : invalid duration", "duration", duration)
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
