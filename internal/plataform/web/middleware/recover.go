package middleware

import (
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

type Recover struct{}

func NewRecover() Middleware {
	return Recover{}
}

func (r Recover) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[ERROR] application recovered from panic : %+v", err)
				ctx.SetStatusCode(http.StatusInternalServerError)
			}
		}()

		h(ctx)
	}
}
