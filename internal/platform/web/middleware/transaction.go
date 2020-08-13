package middleware

import (
	"github.com/b2wdigital/restQL-golang/internal/platform/plugins"
	"github.com/valyala/fasthttp"
)

type Transaction struct {
	lifecycle plugins.Lifecycle
}

func NewTransaction(l plugins.Lifecycle) Middleware {
	return Transaction{lifecycle: l}
}

func (t Transaction) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
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
