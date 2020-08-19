package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Recover struct {
	log *logger.Logger
}

func NewRecover(log *logger.Logger) Middleware {
	return Recover{log: log}
}

func (r Recover) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
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
