package middleware

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/valyala/fasthttp"
	"time"
)

type timeoutConf struct {
	Web struct {
		Middlewares struct {
			Timeout struct {
				Duration string `yaml:"duration"`
			} `yaml:"timeout"`
		} `yaml:"middlewares"`
	} `yaml:"web"`
}

type Timeout struct {
	duration time.Duration
}

func NewTimeout(config conf.Config) Middleware {
	var tc timeoutConf
	err := config.File().Unmarshal(&tc)
	if err != nil {
		return NoopMiddleware{}
	}

	d, parseErr := time.ParseDuration(tc.Web.Middlewares.Timeout.Duration)
	if parseErr != nil {
		return NoopMiddleware{}
	}

	return Timeout{d}
}

func (t Timeout) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nativeCtx := NativeContext(ctx)
		timeout, cancel := context.WithTimeout(nativeCtx, t.duration)
		defer cancel()

		WithNativeContext(ctx, timeout)

		h(ctx)
	}
}
