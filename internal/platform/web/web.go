package web

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/plugins"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/web/middleware"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type handler func(ctx *fasthttp.RequestCtx) error

type app struct {
	config    *conf.Config
	router    *router.Router
	log       restql.Logger
	lifecycle plugins.Lifecycle
}

func newApp(log restql.Logger, config *conf.Config, pm plugins.Lifecycle) app {
	r := router.New()
	r.NotFound = func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBodyString("There is nothing here. =/") }

	return app{router: r, config: config, log: log, lifecycle: pm}
}

func (a app) Handle(method, url string, handler handler) {
	fn := func(ctx *fasthttp.RequestCtx) {
		err := handler(ctx)

		if err != nil {
			a.log.Error("handler has an error", err)

			if err := RespondError(ctx, err); err != nil {
				a.log.Error("failed to send error response", err)
			}
		}
	}

	a.router.Handle(method, url, fn)
}

func (a app) RequestHandler() fasthttp.RequestHandler {
	mws := middleware.FetchEnabled(a.log, a.config, a.lifecycle)
	h := middleware.Apply(a.log, a.router.Handler, mws)
	return h
}

func (a app) RequestHandlerWithoutMiddlewares() fasthttp.RequestHandler {
	return a.router.Handler
}
