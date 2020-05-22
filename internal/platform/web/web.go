package web

import (
	"github.com/b2wdigital/restQL-golang/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/internal/platform/plugins"
	"github.com/b2wdigital/restQL-golang/internal/platform/web/middleware"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Handler func(ctx *fasthttp.RequestCtx) error

type App struct {
	config         *conf.Config
	router         *fasthttprouter.Router
	log            *logger.Logger
	pluginsManager plugins.Manager
}

func NewApp(log *logger.Logger, config *conf.Config, pm plugins.Manager) App {
	r := fasthttprouter.New()
	r.NotFound = func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBodyString("There is nothing here. =/") }

	return App{router: r, config: config, log: log, pluginsManager: pm}
}

func (a App) Handle(method, url string, handler Handler) {
	fn := func(ctx *fasthttp.RequestCtx) {
		err := handler(ctx)

		if err != nil {
			a.log.Error("handler has an error", err)

			if err := RespondError(ctx, err, a.log); err != nil {
				a.log.Error("failed to send error response", err)
			}
		}
	}

	a.router.Handle(method, url, fn)
}

func (a App) RequestHandler() fasthttp.RequestHandler {
	mws := middleware.FetchEnabled(a.log, a.config, a.pluginsManager)
	h := middleware.Apply(a.router.Handler, mws, a.log)
	return h
}

func (a App) RequestHandlerWithoutMiddlewares() fasthttp.RequestHandler {
	return a.router.Handler
}
