package middleware

import (
	"github.com/b2wdigital/restQL-golang/internal/platform/plugins"
	"github.com/valyala/fasthttp"
)

type Transaction struct {
	pluginsManager plugins.Manager
}

func NewTransaction(pm plugins.Manager) Middleware {
	return Transaction{pluginsManager: pm}
}

func (t Transaction) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nativeContext := GetNativeContext(ctx)

		transactionCtx := t.pluginsManager.RunBeforeTransaction(nativeContext, ctx)
		WithNativeContext(ctx, transactionCtx)

		defer func() {
			if reason := recover(); reason != nil {
				t.pluginsManager.RunAfterTransaction(transactionCtx, ctx)
				panic(reason)
			}
		}()

		h(ctx)

		t.pluginsManager.RunAfterTransaction(transactionCtx, ctx)
	}
}
