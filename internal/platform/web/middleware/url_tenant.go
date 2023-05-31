package middleware

import (
	"strings"

	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"

	"github.com/valyala/fasthttp"
)

type urlTenant struct {
	log            restql.Logger
	tenantsByHosts map[string]string
	defaultTenant  string
}

func newUrlTenant(log restql.Logger, defaultTenant string, tenantsByHost map[string]string) Middleware {
	return urlTenant{
		log:            log,
		tenantsByHosts: tenantsByHost,
		defaultTenant:  defaultTenant,
	}
}

func (r urlTenant) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		r.setTenant(ctx)
		h(ctx)
	}
}

func (u urlTenant) setTenant(ctx *fasthttp.RequestCtx) {
	if string(ctx.QueryArgs().Peek("tenant")) != "" {
		u.log.Debug("tenant already set", "tenant", string(ctx.QueryArgs().Peek("tenant")))
		return
	}

	for k, v := range u.tenantsByHosts {
		if strings.Contains(string(ctx.Request.Host()), k) {
			u.log.Debug("setting tenant", "tenant", v, "host", string(ctx.Request.Host()))
			ctx.QueryArgs().Set("tenant", v)
			return
		}
	}
	u.log.Debug("setting default tenant", "tenant", u.defaultTenant, "host", string(ctx.Request.Host()))
	ctx.QueryArgs().Set("tenant", u.defaultTenant)

}
