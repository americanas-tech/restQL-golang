package web

import (
	"github.com/b2wdigital/restQL-golang/v5/internal/platform/web/middleware"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type appOptions struct {
	MiddlewareDecorator *middleware.Decorator
}

type handler func(ctx *fasthttp.RequestCtx) error

type app struct {
	router  *router.Router
	log     restql.Logger
	options appOptions
}

func newApp(log restql.Logger, o appOptions) app {
	r := router.New()
	r.NotFound = func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBodyString("There is nothing here. =/") }

	return app{router: r, log: log, options: o}
}

func (a app) Handle(method, url string, handler handler) {
	fn := func(ctx *fasthttp.RequestCtx) {
		err := handler(ctx)

		if err != nil {
			a.log.Error("handler has an error", err)

			if err := RespondError(ctx, err, errToStatusCode); err != nil {
				a.log.Error("failed to send error response", err)
			}
		}
	}

	a.router.Handle(method, url, fn)
}

func (a app) RequestHandler() fasthttp.RequestHandler {
	h := a.router.Handler
	md := a.options.MiddlewareDecorator
	if md != nil {
		h = md.Apply(h)
	}

	return h
}
