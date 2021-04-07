package middleware

import (
	"github.com/b2wdigital/restQL-golang/v6/internal/platform/plugins"
	"github.com/valyala/fasthttp"
)

type transaction struct {
	lifecycle plugins.Lifecycle
}

func newTransaction(l plugins.Lifecycle) Middleware {
	return transaction{lifecycle: l}
}

func (t transaction) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nativeContext := GetNativeContext(ctx)

		transactionCtx := t.lifecycle.BeforeTransaction(nativeContext, ctx)
		WithNativeContext(ctx, transactionCtx)

		defer func() {
			if reason := recover(); reason != nil {
				t.lifecycle.AfterTransaction(transactionCtx, ctx)
				panic(reason)
			}
		}()

		h(ctx)

		t.lifecycle.AfterTransaction(transactionCtx, ctx)
	}
}
