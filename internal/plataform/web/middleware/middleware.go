package middleware

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/b2wdigital/restQL-golang/internal/plataform/logger"
	"github.com/valyala/fasthttp"
)

type Middleware interface {
	Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler
}

type NoopMiddleware struct{}

func (nm NoopMiddleware) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return h
}

func Apply(h fasthttp.RequestHandler, mws []Middleware, log *logger.Logger) fasthttp.RequestHandler {
	handler := h

	for i := len(mws) - 1; i >= 0; i-- {
		m := mws[i]
		log.Debug(fmt.Sprintf("applying middleware %T", m))
		handler = m.Apply(handler)
	}

	return handler
}

func FetchEnabled(cfg *conf.Config, log *logger.Logger) []Middleware {
	mws := []Middleware{NewRecover(log), NewNativeContext()}

	mwCfg := cfg.Web.Server.Middlewares
	if mwCfg.Timeout != nil {
		mws = append(mws, NewTimeout(mwCfg.Timeout.Duration, log))
	}

	if mwCfg.RequestId != nil {
		mws = append(mws, NewRequestId(mwCfg.RequestId.Header, mwCfg.RequestId.Strategy, log))
	}

	return mws
}
