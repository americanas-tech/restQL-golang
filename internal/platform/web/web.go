package web

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/plugins"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/web/middleware"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"strings"
)

type Handler func(ctx *fasthttp.RequestCtx) error

type App struct {
	config    *conf.Config
	router    *fasthttprouter.Router
	log       *logger.Logger
	lifecycle plugins.Lifecycle
}

func NewApp(log *logger.Logger, config *conf.Config, pm plugins.Lifecycle) App {
	r := fasthttprouter.New()
	r.NotFound = func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBodyString("There is nothing here. =/") }

	return App{router: r, config: config, log: log, lifecycle: pm}
}

func (a App) Handle(method, url string, handler Handler) {
	fn := func(ctx *fasthttp.RequestCtx) {
		err := handler(ctx)

		if err != nil {
			a.log.Error("handler has an error", err)

			if err := RespondError(ctx, err); err != nil {
				a.log.Error("failed to send error response", err)
			}
		}
	}

	normalizedUrls := []string{url}
	if strings.HasSuffix(url, "/") {
		normalizedUrls = append(normalizedUrls, strings.TrimRight(url, "/"))
	} else {
		normalizedUrls = append(normalizedUrls, fmt.Sprintf("%s/", url))
	}

	for _, u := range normalizedUrls {
		a.router.Handle(method, u, fn)
	}
}

func (a App) RequestHandler() fasthttp.RequestHandler {
	mws := middleware.FetchEnabled(a.log, a.config, a.lifecycle)
	h := middleware.Apply(a.router.Handler, mws, a.log)
	return h
}

func (a App) RequestHandlerWithoutMiddlewares() fasthttp.RequestHandler {
	return a.router.Handler
}
