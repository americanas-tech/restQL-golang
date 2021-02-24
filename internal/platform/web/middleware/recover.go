package middleware

import (
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type recoverer struct {
	log restql.Logger
}

func newRecoverer(log restql.Logger) Middleware {
	return recoverer{log: log}
}

func (r recoverer) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if reason := recover(); reason != nil {
				err := errors.Errorf("reason : %v", reason)
				r.log.Error("application recovered from panic", err, "stack", string(debug.Stack()))

				ctx.SetStatusCode(http.StatusInternalServerError)
			}
		}()

		h(ctx)
	}
}
