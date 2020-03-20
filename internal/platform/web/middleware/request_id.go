package middleware

import (
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/valyala/fasthttp"
)

var requestIdKey = []byte{0x1}

type IdGenerator interface {
	Run() string
}

type RequestId struct {
	header    string
	generator IdGenerator
}

var strategyToGenerator = map[string]IdGenerator{
	"base64": NewBase64IdGenerator(),
	"uuid":   NewUuidIdGenerator(),
}

func NewRequestId(header string, strategy string, log *logger.Logger) Middleware {
	if header == "" {
		log.Warn("failed to initialize request id middleware : empty header name")
		return NoopMiddleware{}
	}

	generator, ok := strategyToGenerator[strategy]
	if !ok {
		log.Warn("failed to initialize request id middleware : unknow strategy", "strategy", strategy)
		return NoopMiddleware{}
	}

	return RequestId{
		header:    header,
		generator: generator,
	}
}

func (r RequestId) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var requestId string

		currentRequestId := ctx.Request.Header.Peek(r.header)

		if len(currentRequestId) == 0 {
			requestId = r.generator.Run()
		} else {
			requestId = string(currentRequestId)
		}

		ctx.SetUserValueBytes(requestIdKey, requestId)

		h(ctx)

		ctx.Response.Header.Set(r.header, requestId)
	}
}
