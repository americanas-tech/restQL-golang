package middleware

import (
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/valyala/fasthttp"
	"log"
)

type Middleware interface {
	Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler
}

type NoopMiddleware struct{}

func (nm NoopMiddleware) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return h
}

func Apply(h fasthttp.RequestHandler, mws []Middleware) fasthttp.RequestHandler {
	handler := h

	for i := len(mws) - 1; i >= 0; i-- {
		m := mws[i]
		log.Printf("[DEBUG] applying middleware %T", m)
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

func FetchEnabled(config conf.Config) []Middleware {
	mws := []Middleware{NewRecover(), NewNativeContext()}

	var mc middlewareConf
	err := config.File().Unmarshal(&mc)
	if err != nil {
		log.Printf("[WARN] failed to unmarshal middleware configuration : %s", err)
		return mws
	}

	tc := mc.Web.Middlewares.Timeout
	if tc != nil {
		mws = append(mws, NewTimeout(tc.Duration))
	}

	rc := mc.Web.Middlewares.RequestId
	if rc != nil {
		mws = append(mws, NewRequestId(rc.Header, rc.Strategy))
	}

	return mws
}
