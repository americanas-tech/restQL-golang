package web

import (
	"github.com/valyala/fasthttp"
	"log"
	"net/http/pprof"
)
import "github.com/valyala/fasthttp/fasthttpadaptor"

type Pprof struct {
	index   fasthttp.RequestHandler
	profile fasthttp.RequestHandler
}

func NewPprof() Pprof {
	return Pprof{
		index:   fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index),
		profile: fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile),
	}
}

func (d Pprof) Index(ctx *fasthttp.RequestCtx) error {
	log.Printf("[DEBUG] profile requested")
	d.index(ctx)
	return nil
}

func (d Pprof) Profile(ctx *fasthttp.RequestCtx) error {
	log.Printf("[DEBUG] profile requested")
	d.profile(ctx)
	return nil
}
