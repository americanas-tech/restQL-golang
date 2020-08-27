package middleware

import (
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/valyala/fasthttp"
)

type idGenerator interface {
	Run() string
}

type requestID struct {
	header    string
	generator idGenerator
}

var strategyToGenerator = map[string]idGenerator{
	"base64": newBase64IdGenerator(),
	"uuid":   newUUIDGenerator(),
}

func newRequestID(header string, strategy string, log restql.Logger) Middleware {
	if header == "" {
		log.Warn("failed to initialize request id middleware : empty header name")
		return noopMiddleware{}
	}

	generator, ok := strategyToGenerator[strategy]
	if !ok {
		log.Warn("failed to initialize request id middleware : unknow strategy", "strategy", strategy)
		return noopMiddleware{}
	}

	return requestID{
		header:    header,
		generator: generator,
	}
}

func (r requestID) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var requestID string

		currentRequestID := ctx.Request.Header.Peek(r.header)

		if len(currentRequestID) == 0 {
			requestID = r.generator.Run()
			ctx.Request.Header.Set(r.header, requestID)
		} else {
			requestID = string(currentRequestID)
		}

		h(ctx)

		ctx.Response.Header.Set(r.header, requestID)
	}
}
