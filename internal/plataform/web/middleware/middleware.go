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

func Apply(h fasthttp.RequestHandler, mws []Middleware, log logger.Logger) fasthttp.RequestHandler {
	handler := h

	for i := len(mws) - 1; i >= 0; i-- {
		m := mws[i]
		log.Debug(fmt.Sprintf("applying middleware %T", m))
		handler = m.Apply(handler)
	}

	return handler
}

type requestIdConf struct {
	Header   string `yaml:"header"`
	Strategy string `yaml:"strategy"`
}

type timeoutConf struct {
	Duration string `yaml:"duration"`
}

type middlewareConf struct {
	Web struct {
		Middlewares struct {
			RequestId *requestIdConf `yaml:"requestId"`
			Timeout   *timeoutConf   `yaml:"timeout"`
		} `yaml:"middlewares"`
	} `yaml:"web"`
}

func FetchEnabled(config conf.Config, log logger.Logger) []Middleware {
	mws := []Middleware{NewRecover(log), NewNativeContext()}

	var mc middlewareConf
	err := config.File().Unmarshal(&mc)
	if err != nil {
		log.Warn("failed to unmarshal middleware configuration", "error", err)
		return mws
	}

	tc := mc.Web.Middlewares.Timeout
	if tc != nil {
		mws = append(mws, NewTimeout(tc.Duration, log))
	}

	rc := mc.Web.Middlewares.RequestId
	if rc != nil {
		mws = append(mws, NewRequestId(rc.Header, rc.Strategy, log))
	}

	return mws
}
