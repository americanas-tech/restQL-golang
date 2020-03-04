package web

import (
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/b2wdigital/restQL-golang/internal/plataform/web/middleware"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func API(config conf.Config) fasthttp.RequestHandler {
	r := fasthttprouter.New()

	restQl := NewRestQl(config)

	r.POST("/validate-query", restQl.validateQuery)

	r.NotFound = func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBodyString("There is nothing here. =/") }

	mws := middleware.FetchEnabled(config)
	handler := middleware.Apply(r.Handler, mws)

	return handler
}

func Health(config conf.Config) fasthttp.RequestHandler {
	r := fasthttprouter.New()
	check := NewCheck()

	r.GET("/health", check.health)
	r.GET("/resource-status", check.resourceStatus)

	r.NotFound = func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBodyString("There is nothing here. =/") }

	mws := middleware.FetchEnabled(config)
	handler := middleware.Apply(r.Handler, mws)

	return handler
}
