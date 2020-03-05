package middleware

import (
	"context"
	"github.com/valyala/fasthttp"
)

type NativeContext struct{}

func NewNativeContext() NativeContext {
	return NativeContext{}
}

func (n NativeContext) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		WithNativeContext(ctx, context.Background())

		h(ctx)
	}
}

func GetNativeContext(ctx *fasthttp.RequestCtx) context.Context {
	return ctx.UserValue("context").(context.Context)
}

func WithNativeContext(ctx *fasthttp.RequestCtx, nativeCtx context.Context) {
	ctx.SetUserValue("context", nativeCtx)
}
