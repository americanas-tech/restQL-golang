package web

import (
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/valyala/fasthttp"
	"net/http"
)

func API(config conf.Config) fasthttp.RequestHandler {
	app := NewApp(config)
	restQl := NewRestQl(config)

	app.Handle(http.MethodPost, "/validate-query", restQl.validateQuery)

	return app.RequestHandler()
}

func Health(config conf.Config) fasthttp.RequestHandler {
	app := NewApp(config)
	check := NewCheck()

	app.Handle(http.MethodGet, "/health", check.health)
	app.Handle(http.MethodGet, "/resource-status", check.resourceStatus)

	return app.RequestHandler()
}
