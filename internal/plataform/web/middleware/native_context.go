package middleware

import (
	"context"
	"github.com/valyala/fasthttp"
)

func NativeContext(ctx *fasthttp.RequestCtx) context.Context {
	return ctx.UserValue("context").(context.Context)
}

func WithNativeContext(ctx *fasthttp.RequestCtx, nativeCtx context.Context) {
	ctx.SetUserValue("context", nativeCtx)
}
