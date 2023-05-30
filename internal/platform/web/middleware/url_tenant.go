package middleware

import (
	"strings"

	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"

	"github.com/valyala/fasthttp"
)

type urlTenant struct {
	log restql.Logger
}

func newUrlTenant(log restql.Logger) Middleware {
	return urlTenant{log: log}
}

func (r urlTenant) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if strings.Contains(string(ctx.Request.Host()), "americanas.teste") {
			ctx.QueryArgs().Add("tenant", "acom-npf")
		}
		h(ctx)
	}
}
