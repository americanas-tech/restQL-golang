package web

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"log"
	"net/http/pprof"
)

type debug struct {
	index   fasthttp.RequestHandler
	profile fasthttp.RequestHandler
}

func newDebug() debug {
	return debug{
		index:   fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index),
		profile: fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile),
	}
}

func (d debug) Index(ctx *fasthttp.RequestCtx) error {
	log.Printf("[DEBUG] profile requested")
	d.index(ctx)
	return nil
}

func (d debug) Profile(ctx *fasthttp.RequestCtx) error {
	log.Printf("[DEBUG] profile requested")
	d.profile(ctx)
	return nil
}
