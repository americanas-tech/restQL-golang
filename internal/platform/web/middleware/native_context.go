package middleware

import (
	"context"
	"github.com/valyala/fasthttp"
)

type nativeContext struct {
	cm *ConnManager
}

func newNativeContext(cm *ConnManager) nativeContext {
	return nativeContext{cm: cm}
}

func (n nativeContext) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(reqCtx *fasthttp.RequestCtx) {
		ctx := n.cm.ContextForConnection(reqCtx.Conn())

		WithNativeContext(reqCtx, ctx)

		h(reqCtx)
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
