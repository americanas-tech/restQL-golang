package middleware

import (
	"context"
	"github.com/valyala/fasthttp"
)

type nativeContext struct{}

func newNativeContext() nativeContext {
	return nativeContext{}
}

func (n nativeContext) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		WithNativeContext(ctx, context.Background())

		h(ctx)
	}
}

// GetNativeContext retrieve a standard library context
// from FastHTTP request context.
func GetNativeContext(ctx *fasthttp.RequestCtx) context.Context {
	return ctx.UserValue("context").(context.Context)
}

// WithNativeContext stores a standard library context
//into FastHTTP request context.
func WithNativeContext(ctx *fasthttp.RequestCtx, nativeCtx context.Context) {
	ctx.SetUserValue("context", nativeCtx)
}
