package middleware

import (
	"bytes"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
	"net/http"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

// corsOptions is a configuration container to setup the CORS middleware.
type corsOptions struct {
	// AllowedOrigins is a comma separated list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is "*"
	AllowedOrigins string
	// AllowedMethods is a comma separated list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods string
	// AllowedHeaders is a comma separated list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders string
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification.
	ExposedHeaders string
	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached by the client.
	MaxAge int
	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool
}

// cors middleware
type cors struct {
	log                    restql.Logger
	allowedOrigins         [][]byte
	allowedWildcardOrigins []wildcard
	allowedHeaders         []byte
	allowedMethods         []byte
	exposedHeaders         []byte
	maxAge                 []byte
	allowedOriginsAll      bool
	allowedHeadersAll      bool
	allowCredentials       bool
}

// newCors creates a new cors middleware with the provided options.
func newCors(log restql.Logger, options corsOptions) *cors {
	c := &cors{
		allowCredentials: options.AllowCredentials,
		log:              log,
	}

	parseOptions(c, options)

	return c
}

func parseOptions(c *cors, options corsOptions) {
	// Normalize options
	// Note: for origins and methods matching, the spec requires a case-sensitive matching.
	// As it may error prone, we chose to ignore the spec here.

	// Allowed Origins
	if options.AllowedOrigins == "" {
		// Default is all origins
		c.allowedOriginsAll = true
	} else {
		origins := strings.Split(options.AllowedOrigins, ",")
		origins = convert(origins, strings.TrimSpace)

		c.allowedOrigins = [][]byte{}
		c.allowedWildcardOrigins = []wildcard{}
		for _, origin := range origins {
			if origin == "*" {
				// If "*" is present in the list, turn the whole list into a match all
				c.allowedOriginsAll = true
				c.allowedOrigins = nil
				c.allowedWildcardOrigins = nil
				break
			} else if i := strings.IndexByte(origin, '*'); i >= 0 {
				// Split the origin in two: start and end string without the *
				w := wildcard{prefix: []byte(origin[0:i]), suffix: []byte(origin[i+1:])}
				c.allowedWildcardOrigins = append(c.allowedWildcardOrigins, w)
			} else {
				c.allowedOrigins = append(c.allowedOrigins, []byte(origin))
			}
		}
	}

	// Allowed Headers
	var headers []string
	if options.AllowedHeaders == "" {
		// Use sensible defaults
		headers = []string{"Origin", "Accept", "Content-Type", "X-Requested-With"}
	} else {
		// Origin is always appended as some browsers will always request for this header at preflight
		headers = strings.Split(options.AllowedHeaders, ",")
		headers = convert(headers, strings.TrimSpace)
		headers = append(headers, "Origin")
		headers = convert(headers, http.CanonicalHeaderKey)

		for _, h := range headers {
			if h == "*" {
				c.allowedHeadersAll = true
				c.allowedHeaders = nil
				break
			}
		}
	}
	if !c.allowedHeadersAll {
		c.allowedHeaders = []byte(strings.Join(headers, ", "))
	}

	if options.ExposedHeaders != "" {
		exposedHeaders := strings.Split(options.ExposedHeaders, ",")
		exposedHeaders = convert(exposedHeaders, strings.TrimSpace)
		exposedHeaders = convert(exposedHeaders, http.CanonicalHeaderKey)

		c.exposedHeaders = []byte(strings.Join(exposedHeaders, ", "))
	}

	// Allowed Methods
	var methods []string
	if options.AllowedMethods == "" {
		// Default is spec's "simple" methods
		methods = []string{http.MethodGet, http.MethodPost, http.MethodHead}
	} else {
		methods = strings.Split(options.AllowedMethods, ",")
		methods = convert(methods, strings.TrimSpace)
		methods = convert(methods, strings.ToUpper)
	}

	c.allowedMethods = []byte(strings.Join(methods, ", "))

	if options.MaxAge > 0 {
		c.maxAge = []byte(strconv.Itoa(options.MaxAge))
	}
}

// Apply wraps a request handler with the CORS middleware.
func (c *cors) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		method := string(ctx.Method())

		if method == fasthttp.MethodOptions {
			c.log.Debug("handling preflight request")
			c.handlePreflight(ctx)
			// Preflight requests are standalone and should stop the chain as some other
			// middleware may not handle OPTIONS requests correctly. One typical example
			// is authentication middleware ; OPTIONS requests won't carry authentication
			// headers (see rs/cors#1)
			ctx.SetStatusCode(fasthttp.StatusOK)
		} else {
			c.log.Debug("handling actual request")
			c.handleActualRequest(ctx)
			h(ctx)
		}
	}
}

var (
	// Response headers names
	accessControlAllowOrigin      = []byte("Access-Control-Allow-Origin")
	accessControlAllowMethods     = []byte("Access-Control-Allow-Methods")
	accessControlAllowHeaders     = []byte("Access-Control-Allow-Headers")
	accessControlExposeHeaders    = []byte("Access-Control-Expose-Headers")
	accessControlAllowCredentials = []byte("Access-Control-Allow-Credentials")
	accessControlMaxAge           = []byte("Access-Control-Max-Age")
	vary                          = []byte("Vary")

	// Vary header values
	varyOrigin                      = []byte("Origin")
	varyAccessControlRequestMethod  = []byte("Access-Control-Request-Method")
	varyAccessControlRequestHeaders = []byte("Access-Control-Request-Headers")
)

// handlePreflight handles pre-flight CORS requests
func (c *cors) handlePreflight(ctx *fasthttp.RequestCtx) {
	headers := &ctx.Response.Header
	origin := ctx.Request.Header.Peek("Origin")

	// Always set Vary headers
	// see https://github.com/rs/cors/issues/10,
	//     https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	headers.AddBytesKV(vary, varyOrigin)
	headers.AddBytesKV(vary, varyAccessControlRequestMethod)
	headers.AddBytesKV(vary, varyAccessControlRequestHeaders)

	if len(origin) == 0 {
		c.log.Debug("preflight request missing origin")
		return
	}

	if c.allowedOriginsAll {
		headers.SetBytesK(accessControlAllowOrigin, "*")
	} else {
		if c.isOriginAllowed(origin) {
			headers.SetBytesKV(accessControlAllowOrigin, origin)
		}
	}

	headers.SetBytesKV(accessControlAllowMethods, c.allowedMethods)

	if c.allowedHeadersAll {
		rh := ctx.Request.Header.Peek("Access-Control-Request-Headers")
		headers.SetBytesKV(accessControlAllowHeaders, rh)
	} else {
		headers.SetBytesKV(accessControlAllowHeaders, c.allowedHeaders)
	}

	if c.allowCredentials {
		headers.SetBytesK(accessControlAllowCredentials, "true")
	}
	if c.maxAge != nil {
		headers.SetBytesKV(accessControlMaxAge, c.maxAge)
	}
}

// handleActualRequest handles simple cross-origin requests, actual request or redirects
func (c *cors) handleActualRequest(ctx *fasthttp.RequestCtx) {
	headers := &ctx.Response.Header
	origin := ctx.Request.Header.Peek("Origin")

	// Always set Vary, see https://github.com/rs/cors/issues/10
	headers.AddBytesKV(vary, varyOrigin)
	if len(origin) == 0 {
		c.log.Debug("actual request missing origin")
		return
	}

	if c.allowedOriginsAll {
		headers.SetBytesK(accessControlAllowOrigin, "*")
	} else {
		if c.isOriginAllowed(origin) {
			headers.SetBytesKV(accessControlAllowOrigin, origin)
		} else {
			c.log.Debug("origin not allowed", "origin", string(origin))
		}
	}
	if len(c.exposedHeaders) > 0 {
		headers.SetBytesKV(accessControlExposeHeaders, c.exposedHeaders)
	}
	if c.allowCredentials {
		headers.SetBytesK(accessControlAllowCredentials, "true")
	}
}

// isOriginAllowed checks if a given origin is allowed to perform cross-domain requests
// on the endpoint
func (c *cors) isOriginAllowed(origin []byte) bool {
	if c.allowedOriginsAll {
		return true
	}
	for _, o := range c.allowedOrigins {
		if bytes.EqualFold(o, origin) {
			return true
		}
	}

	for _, w := range c.allowedWildcardOrigins {
		if w.match(origin) {
			return true
		}
	}

	return false
}

type converter func(string) string

// convert converts a list of string using the passed converter function
func convert(s []string, c converter) []string {
	out := make([]string, len(s))
	for i, si := range s {
		out[i] = c(si)
	}
	return out
}

type wildcard struct {
	prefix []byte
	suffix []byte
}

func (w wildcard) match(s []byte) bool {
	if len(s) < len(w.prefix)+len(w.suffix) {
		return false
	}

	sp := s[:len(w.prefix)]
	if !bytes.EqualFold(sp, w.prefix) {
		return false
	}

	ss := s[len(s)-len(w.suffix):]
	if !bytes.EqualFold(ss, w.suffix) {
		return false
	}

	return true
}
