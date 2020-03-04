package middleware

import (
	"bytes"
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

func NewRequestId(header string, strategy string) Middleware {
	return RequestId{
		header:    header,
		generator: strategyToGenerator[strategy],
	}
}

func (r RequestId) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var requestId string

		currentRequestId := ctx.Request.Header.Peek(r.header)

		if len(currentRequestId) == 0 {
			requestId = r.generator.Run()
		} else {
			requestId = bytes.NewBuffer(currentRequestId).String()
		}

		ctx.SetUserValueBytes(requestIdKey, requestId)

		h(ctx)

		ctx.Response.Header.Set(r.header, requestId)
	}
}
