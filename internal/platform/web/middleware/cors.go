package middleware

import (
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"strings"

	"github.com/valyala/fasthttp"
)

var (
	accessControlAllowOriginHeaderName   = []byte("Access-Control-Allow-Origin")
	accessControlAllowMethodsHeaderName  = []byte("Access-Control-Allow-Methods")
	accessControlAllowHeadersHeaderName  = []byte("Access-Control-Allow-Headers")
	accessControlExposeHeadersHeaderName = []byte("Access-Control-Expose-Headers")

	accessControlRequestHeadersHeaderName = []byte("Access-Control-Request-Headers")
	accessControlRequestMethodHeaderName  = []byte("Access-Control-Request-Method")
	originHeaderName                      = []byte("Origin")
)

type option func(c *cors)

func withAllowOrigins(allowedOrigins string) option {
	return func(c *cors) {
		origins := strings.Split(allowedOrigins, ",")
		allowedOriginSet := make(map[string]struct{})

		for _, o := range origins {
			o := strings.TrimSpace(o)
			if o == "*" {
				c.allowedOriginsAll = true
			}
			allowedOriginSet[o] = struct{}{}
		}

		c.allowedOriginSet = allowedOriginSet
	}
}

func withAllowHeaders(allowedHeaders string) option {
	return func(c *cors) {
		headers := strings.Split(allowedHeaders, ",")
		allowedHeadersSet := make(map[string]struct{})

		for _, h := range headers {
			h := strings.TrimSpace(h)
			if h == "*" {
				c.allowedHeadersAll = true
			}
			allowedHeadersSet[h] = struct{}{}
		}

		c.allowedHeadersSet = allowedHeadersSet
	}
}

func withAllowMethods(allowedMethods string) option {
	return func(c *cors) {
		methods := strings.Split(allowedMethods, ",")
		allowedMethodsSet := make(map[string]struct{})

		for _, m := range methods {
			allowedMethodsSet[strings.TrimSpace(m)] = struct{}{}
		}

		c.allowedMethodsSet = allowedMethodsSet
	}
}

func withExposedHeaders(exposedHeaders string) option {
	return func(c *cors) {
		c.exposedHeaders = exposedHeaders
	}
}

type cors struct {
	allowedOriginSet  map[string]struct{}
	allowedOriginsAll bool
	allowedHeadersSet map[string]struct{}
	allowedHeadersAll bool
	allowedMethodsSet map[string]struct{}
	exposedHeaders    string
	logger            restql.Logger
}

func newCors(log restql.Logger, options ...option) Middleware {
	c := cors{logger: log}

	for _, option := range options {
		option(&c)
	}

	return &c
}

func (c *cors) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Method()) == fasthttp.MethodOptions {
			c.handlePreflight(ctx)
			ctx.SetStatusCode(200)
		} else {
			c.handleActual(ctx)
			h(ctx)
		}
	}
}

func (c *cors) handlePreflight(ctx *fasthttp.RequestCtx) {
	originHeader := ctx.Request.Header.PeekBytes(originHeaderName)
	if len(originHeader) == 0 || !c.isAllowedOrigin(originHeader) {
		c.logger.Debug("origin is not allowed", "origin", originHeader)
		return
	}

	method := ctx.Request.Header.PeekBytes(accessControlRequestMethodHeaderName)
	if !c.isAllowedMethod(method) {
		c.logger.Debug("method is not allowed", "method", method)
		return
	}

	headers, allowed := c.extractAndVerifyAccessControlRequestHeaders(ctx)
	if !allowed {
		c.logger.Debug("headers not allowed", "headers", headers)
		return
	}

	if headers != "" {
		ctx.Response.Header.SetBytesK(accessControlAllowHeadersHeaderName, headers)
	}
	ctx.Response.Header.SetBytesKV(accessControlAllowOriginHeaderName, originHeader)
	ctx.Response.Header.SetBytesKV(accessControlAllowMethodsHeaderName, method)
}

func (c *cors) handleActual(ctx *fasthttp.RequestCtx) {
	originHeader := ctx.Request.Header.PeekBytes(originHeaderName)
	if len(originHeader) == 0 || !c.isAllowedOrigin(originHeader) {
		c.logger.Debug("origin is not allowed", "origin", originHeader)
		return
	}

	ctx.Response.Header.SetBytesKV(accessControlAllowOriginHeaderName, originHeader)
	if c.exposedHeaders != "" {
		ctx.Response.Header.SetBytesK(accessControlExposeHeadersHeaderName, c.exposedHeaders)
	}
}

func (c *cors) isAllowedOrigin(originHeader []byte) bool {
	if c.allowedOriginsAll {
		return true
	}

	origin := string(originHeader)
	_, found := c.allowedOriginSet[origin]

	return found
}

func (c *cors) isAllowedMethod(methodHeader []byte) bool {
	if len(c.allowedMethodsSet) == 0 {
		return false
	}

	method := string(methodHeader)

	if method == "OPTIONS" {
		return true
	}

	_, found := c.allowedMethodsSet[method]
	return found
}

func (c *cors) areHeadersAllowed(headers []string) bool {
	if c.allowedHeadersAll || len(headers) == 0 {
		return true
	}

	for _, header := range headers {
		_, found := c.allowedHeadersSet[header]

		if !found {
			return false
		}
	}

	return true
}

func (c *cors) extractAndVerifyAccessControlRequestHeaders(ctx *fasthttp.RequestCtx) (string, bool) {
	if len(ctx.Request.Header.PeekBytes(accessControlRequestHeadersHeaderName)) == 0 {
		return "", true
	}

	accessControlRequestHeaders := string(ctx.Request.Header.PeekBytes(accessControlRequestHeadersHeaderName))
	headers := strings.Split(accessControlRequestHeaders, ",")

	return accessControlRequestHeaders, c.areHeadersAllowed(headers)
}
