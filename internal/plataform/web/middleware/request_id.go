package middleware

import (
	"bytes"
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
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

type requestIdConf struct {
	Web struct {
		Middlewares struct {
			RequestId struct {
				Header   string `yaml:"header"`
				Strategy string `yaml:"strategy"`
			} `yaml:"requestId"`
		} `yaml:"middlewares"`
	} `yaml:"web"`
}

func NewRequestId(config conf.Config) Middleware {
	var rc requestIdConf
	err := config.File().Unmarshal(&rc)
	if err != nil {
		return NoopMiddleware{}
	}

	return RequestId{
		header:    rc.Web.Middlewares.RequestId.Header,
		generator: strategyToGenerator[rc.Web.Middlewares.RequestId.Strategy],
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
