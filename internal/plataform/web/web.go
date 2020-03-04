package web

import (
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/b2wdigital/restQL-golang/internal/plataform/web/middleware"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
)

type Handler func(ctx *fasthttp.RequestCtx) error

type App struct {
	config conf.Config
	router *fasthttprouter.Router
}

func NewApp(config conf.Config) App {
	r := fasthttprouter.New()
	r.NotFound = func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBodyString("There is nothing here. =/") }

	return App{router: r, config: config}
}

func (a App) Handle(method, url string, handler Handler) {
	fn := func(ctx *fasthttp.RequestCtx) {
		err := handler(ctx)

		if err != nil {
			log.Printf("[ERROR] handler has an error : %s", err)

			if err := RespondError(ctx, err); err != nil {
				log.Printf("[ERROR] failed to send error response : %s", err)
			}
		}
	}

	a.router.Handle(method, url, fn)
}

func (a App) RequestHandler() fasthttp.RequestHandler {
	mws := middleware.FetchEnabled(a.config)
	h := middleware.Apply(a.router.Handler, mws)
	return h
}

func (a App) RequestHandlerWithoutMiddlewares() fasthttp.RequestHandler {
	return a.router.Handler
}
