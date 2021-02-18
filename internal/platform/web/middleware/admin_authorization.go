package middleware

import (
	"bytes"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/valyala/fasthttp"
)

var adminPrefix = []byte("/admin")

type adminAuthorization struct {
	log  restql.Logger
	code []byte
}

func newAdminAuthorization(log restql.Logger, code string) adminAuthorization {
	return adminAuthorization{log: log, code: []byte(code)}
}

func (a adminAuthorization) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		isAdminRequest := bytes.Contains(ctx.Request.URI().Path(), adminPrefix)
		method := string(ctx.Method())

		if isAdminRequest && (method != fasthttp.MethodOptions && method != fasthttp.MethodGet) {
			bearerCode := getBearerToken(ctx)
			if len(bearerCode) == 0 {
				ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
				return
			}

			bearerCode = bytes.TrimPrefix(bearerCode, []byte("Bearer"))
			bearerCode = bytes.TrimPrefix(bearerCode, []byte("bearer"))
			bearerCode = bytes.TrimSpace(bearerCode)

			if !bytes.Equal(bearerCode, a.code) {
				ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
				return
			}
		}

		h(ctx)
	}
}

func getBearerToken(ctx *fasthttp.RequestCtx) []byte {
	bearerCode := ctx.Request.Header.Peek("Authorization")
	if len(bearerCode) > 0 {
		return bearerCode
	}

	bearerCode = ctx.Request.Header.Peek("authorization")

	return bearerCode
}
