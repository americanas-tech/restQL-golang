package middleware

import (
	"github.com/b2wdigital/restQL-golang/v6/test"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
)

var testHandler = fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
	ctx.SetBodyString("bar")
})

var allHeaders = []string{
	"Vary",
	"Access-Control-Allow-Origin",
	"Access-Control-Allow-Methods",
	"Access-Control-Allow-Headers",
	"Access-Control-Allow-Credentials",
	"Access-Control-Max-Age",
	"Access-Control-Expose-Headers",
}

func assertHeaders(t *testing.T, resHeaders *fasthttp.ResponseHeader, expHeaders map[string]string) {
	headers := make(map[string][]string, resHeaders.Len())
	resHeaders.VisitAll(func(key []byte, value []byte) {
		if arr, found := headers[string(key)]; found {
			headers[string(key)] = append(arr, string(value))
		} else {
			headers[string(key)] = []string{string(value)}
		}
	})

	for _, name := range allHeaders {
		got := strings.Join(headers[name], ", ")
		want := expHeaders[name]
		if got != want {
			t.Errorf("Response header %q = %q, want %q", name, got, want)
		}
	}
}

func TestSpec(t *testing.T) {
	cases := []struct {
		name       string
		options    corsOptions
		method     string
		reqHeaders map[string]string
		resHeaders map[string]string
	}{
		{
			"NoConfig",
			corsOptions{
				// Intentionally left blank.
			},
			"GET",
			map[string]string{},
			map[string]string{
				"Vary": "Origin",
			},
		},
		{
			"MatchAllOrigin",
			corsOptions{
				AllowedOrigins: "*",
			},
			"GET",
			map[string]string{
				"Origin": "http://foobar.com",
			},
			map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "*",
			},
		},
		{
			"MatchAllOriginWithCredentials",
			corsOptions{
				AllowedOrigins:   "*",
				AllowCredentials: true,
			},
			"GET",
			map[string]string{
				"Origin": "http://foobar.com",
			},
			map[string]string{
				"Vary":                             "Origin",
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			"AllowedOrigin",
			corsOptions{
				AllowedOrigins: "http://hero.com, http://foobar.com",
			},
			"GET",
			map[string]string{
				"Origin": "http://foobar.com",
			},
			map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "http://foobar.com",
			},
		},
		{
			"WildcardOrigin",
			corsOptions{
				AllowedOrigins: "http://hero.com, http://*.bar.com",
			},
			"GET",
			map[string]string{
				"Origin": "http://foo.bar.com",
			},
			map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "http://foo.bar.com",
			},
		},
		{
			"DisallowedOrigin",
			corsOptions{
				AllowedOrigins: "http://hero.com, http://foobar.com",
				ExposedHeaders: "X-Header-1, x-Header-2",
			},
			"GET",
			map[string]string{
				"Origin": "http://barbaz.com",
			},
			map[string]string{
				"Vary":                          "Origin",
				"Access-Control-Expose-Headers": "X-Header-1, X-Header-2",
			},
		},
		{
			"DisallowedWildcardOrigin",
			corsOptions{
				AllowedOrigins: "http://hero.com, http://*.bar.com",
			},
			"GET",
			map[string]string{
				"Origin": "http://foo.baz.com",
			},
			map[string]string{
				"Vary": "Origin",
			},
		},
		{
			"PreflightDisallowedOrigin",
			corsOptions{
				AllowedOrigins: "http://hero.com, http://foobar.com",
			},
			"OPTIONS",
			map[string]string{
				"Origin": "http://barbaz.com",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Headers": "Origin, Accept, Content-Type, X-Requested-With",
				"Access-Control-Allow-Methods": "GET, POST, HEAD",
			},
		},
		{
			"PreflightDisallowedWildcardOrigin",
			corsOptions{
				AllowedOrigins: "http://hero.com, http://*.bar.com",
			},
			"OPTIONS",
			map[string]string{
				"Origin": "http://foo.baz.com",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Headers": "Origin, Accept, Content-Type, X-Requested-With",
				"Access-Control-Allow-Methods": "GET, POST, HEAD",
			},
		},
		{
			"MaxAge",
			corsOptions{
				AllowedOrigins: "http://example.com/",
				AllowedMethods: "GET",
				MaxAge:         10,
			},
			"OPTIONS",
			map[string]string{
				"Origin":                        "http://example.com/",
				"Access-Control-Request-Method": "GET",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://example.com/",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Max-Age":       "10",
				"Access-Control-Allow-Headers": "Origin, Accept, Content-Type, X-Requested-With",
			},
		},
		{
			"AllowedMethod",
			corsOptions{
				AllowedOrigins: "http://foobar.com",
				AllowedMethods: "PUT, delete",
			},
			"OPTIONS",
			map[string]string{
				"Origin":                        "http://foobar.com",
				"Access-Control-Request-Method": "PUT",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://foobar.com",
				"Access-Control-Allow-Methods": "PUT, DELETE",
				"Access-Control-Allow-Headers": "Origin, Accept, Content-Type, X-Requested-With",
			},
		},
		{
			"AllowedHeaders",
			corsOptions{
				AllowedOrigins: "http://foobar.com",
				AllowedHeaders: "X-Header-1, x-header-2",
			},
			"OPTIONS",
			map[string]string{
				"Origin":                         "http://foobar.com",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "X-Header-2, X-HEADER-1",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://foobar.com",
				"Access-Control-Allow-Methods": "GET, POST, HEAD",
				"Access-Control-Allow-Headers": "X-Header-1, X-Header-2, Origin",
			},
		},
		{
			"DefaultAllowedHeaders",
			corsOptions{
				AllowedOrigins: "http://foobar.com",
				AllowedHeaders: "",
			},
			"OPTIONS",
			map[string]string{
				"Origin":                         "http://foobar.com",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "X-Requested-With",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://foobar.com",
				"Access-Control-Allow-Methods": "GET, POST, HEAD",
				"Access-Control-Allow-Headers": "Origin, Accept, Content-Type, X-Requested-With",
			},
		},
		{
			"AllowedAllHeaders",
			corsOptions{
				AllowedOrigins: "http://foobar.com",
				AllowedHeaders: "*",
			},
			"OPTIONS",
			map[string]string{
				"Origin":                         "http://foobar.com",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "X-Header-2, X-HEADER-1",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://foobar.com",
				"Access-Control-Allow-Methods": "GET, POST, HEAD",
				"Access-Control-Allow-Headers": "X-Header-2, X-HEADER-1",
			},
		},
		{
			"OriginHeader",
			corsOptions{
				AllowedOrigins: "http://foobar.com",
			},
			"OPTIONS",
			map[string]string{
				"Origin":                         "http://foobar.com",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "origin",
			},
			map[string]string{
				"Vary":                         "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":  "http://foobar.com",
				"Access-Control-Allow-Methods": "GET, POST, HEAD",
				"Access-Control-Allow-Headers": "Origin, Accept, Content-Type, X-Requested-With",
			},
		},
		{
			"ExposedHeader",
			corsOptions{
				AllowedOrigins: "http://foobar.com",
				ExposedHeaders: "X-Header-1, x-header-2",
			},
			"GET",
			map[string]string{
				"Origin": "http://foobar.com",
			},
			map[string]string{
				"Vary":                          "Origin",
				"Access-Control-Allow-Origin":   "http://foobar.com",
				"Access-Control-Expose-Headers": "X-Header-1, X-Header-2",
			},
		},
		{
			"AllowedCredentials",
			corsOptions{
				AllowedOrigins:   "http://foobar.com",
				AllowCredentials: true,
			},
			"OPTIONS",
			map[string]string{
				"Origin":                        "http://foobar.com",
				"Access-Control-Request-Method": "GET",
			},
			map[string]string{
				"Vary":                             "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Origin":      "http://foobar.com",
				"Access-Control-Allow-Methods":     "GET, POST, HEAD",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Allow-Headers":     "Origin, Accept, Content-Type, X-Requested-With",
			},
		},
	}
	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			s := newCors(test.NoOpLogger, tc.options)

			ctx := fasthttp.RequestCtx{}
			ctx.Request.Header.SetMethod(tc.method)
			ctx.Request.SetRequestURI("http://example.com/foo")
			for name, value := range tc.reqHeaders {
				ctx.Request.Header.Add(name, value)
			}

			s.Apply(testHandler)(&ctx)
			assertHeaders(t, &ctx.Response.Header, tc.resHeaders)
		})
	}
}

func TestHandlePreflightEmptyOriginAbortion(t *testing.T) {
	s := newCors(test.NoOpLogger, corsOptions{
		AllowedOrigins: "http://foo.com",
	})

	ctx := fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("OPTIONS")
	ctx.Request.SetRequestURI("http://example.com/foo")

	s.handlePreflight(&ctx)

	assertHeaders(t, &ctx.Response.Header, map[string]string{
		"Vary": "Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
	})
}

func TestHandleActualRequestEmptyOriginAbortion(t *testing.T) {
	s := newCors(test.NoOpLogger, corsOptions{
		AllowedOrigins: "http://foo.com",
	})

	ctx := fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("http://example.com/foo")

	s.handleActualRequest(&ctx)

	assertHeaders(t, &ctx.Response.Header, map[string]string{
		"Vary": "Origin",
	})
}

// Utils testing

func TestConvert(t *testing.T) {
	s := convert([]string{"A", "b", "C"}, strings.ToLower)
	e := []string{"a", "b", "c"}
	if s[0] != e[0] || s[1] != e[1] || s[2] != e[2] {
		t.Errorf("%v != %v", s, e)
	}
}
