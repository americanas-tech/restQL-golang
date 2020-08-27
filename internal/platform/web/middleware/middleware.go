package middleware

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"

	"github.com/b2wdigital/restQL-golang/v4/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/plugins"
	"github.com/valyala/fasthttp"
)

// Middleware defines a generic handler wrapper
// capable of controlling request flow.
type Middleware interface {
	Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler
}

type noopMiddleware struct{}

func (nm noopMiddleware) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return h
}

// Apply takes a base handler and a slice of middlewares,
// applying each one in the order they are given.
func Apply(log restql.Logger, h fasthttp.RequestHandler, mws []Middleware) fasthttp.RequestHandler {
	handler := h

	for i := len(mws) - 1; i >= 0; i-- {
		m := mws[i]
		log.Debug(fmt.Sprintf("applying middleware %T", m))
		handler = m.Apply(handler)
	}

	return handler
}

// FetchEnabled returns all middlewares enabled in configuration
func FetchEnabled(log restql.Logger, cfg *conf.Config, pm plugins.Lifecycle) []Middleware {
	mws := []Middleware{newRecoverer(log), newNativeContext(), newTransaction(pm)}

	mwCfg := cfg.HTTP.Server.Middlewares
	if mwCfg.Timeout != nil {
		mws = append(mws, newTimeout(mwCfg.Timeout.Duration, log))
	}

	if mwCfg.RequestID != nil {
		mws = append(mws, newRequestID(mwCfg.RequestID.Header, mwCfg.RequestID.Strategy, log))
	}

	if mwCfg.Cors != nil {
		cors := newCors(log,
			withAllowOrigins(mwCfg.Cors.AllowOrigin),
			withAllowHeaders(mwCfg.Cors.AllowHeaders),
			withAllowMethods(mwCfg.Cors.AllowMethods),
			withExposedHeaders(mwCfg.Cors.ExposeHeaders),
		)
		mws = append(mws, cors)
	}

	return mws
}
