package runner

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"net/http"
	"strings"
	"time"
)

var disallowedHeaders = map[string]struct{}{
	"Host":            {},
	"Content-Type":    {},
	"Content-Length":  {},
	"Connection":      {},
	"Origin":          {},
	"Accept-Encoding": {},
}

var queryMethodToHttpMethod = map[string]string{
	domain.FromMethod:   http.MethodGet,
	domain.ToMethod:     http.MethodPost,
	domain.IntoMethod:   http.MethodPut,
	domain.UpdateMethod: http.MethodPatch,
	domain.DeleteMethod: http.MethodDelete,
}

func MakeRequest(defaultResourceTimeout time.Duration, forwardPrefix string, statement domain.Statement, queryCtx domain.QueryContext) domain.HttpRequest {
	mapping := queryCtx.Mappings[statement.Resource]
	method := queryMethodToHttpMethod[statement.Method]
	headers := makeHeaders(statement, queryCtx)
	path := mapping.PathWithParams(statement.With)
	timeout := parseTimeout(defaultResourceTimeout, statement)

	req := domain.HttpRequest{
		Method:  method,
		Schema:  mapping.Schema,
		Host:    mapping.Host,
		Path:    path,
		Headers: headers,
		Timeout: timeout,
	}

	if statement.Method == domain.FromMethod || statement.Method == domain.DeleteMethod {
		req.Query = makeQueryParams(forwardPrefix, statement, mapping, queryCtx)
	} else {
		req.Query = getForwardParams(forwardPrefix, queryCtx)
	}

	if statement.Method == domain.ToMethod || statement.Method == domain.UpdateMethod || statement.Method == domain.IntoMethod {
		req.Body = makeBody(statement, mapping)
	}

	return req
}

func makeBody(statement domain.Statement, mapping domain.Mapping) domain.Body {
	result := make(map[string]interface{})
	for key, value := range statement.With {
		if mapping.HasParam(key) {
			continue
		}

		result[key] = value
	}
	return result
}

func makeHeaders(statement domain.Statement, queryCtx domain.QueryContext) map[string]string {
	headers := getForwardHeaders(queryCtx)
	for key, value := range statement.Headers {
		str, ok := value.(string)
		if !ok {
			continue
		}
		headers[key] = str
	}
	headers["Content-Type"] = "application/json"

	return headers
}

func getForwardHeaders(queryCtx domain.QueryContext) map[string]string {
	r := make(map[string]string)
	for k, v := range queryCtx.Input.Headers {
		if _, found := disallowedHeaders[k]; !found {
			r[k] = v
		}
	}
	return r
}

func makeQueryParams(forwardPrefix string, statement domain.Statement, mapping domain.Mapping, queryCtx domain.QueryContext) map[string]interface{} {
	queryArgs := getForwardParams(forwardPrefix, queryCtx)
	for key, value := range statement.With {
		if mapping.HasParam(key) {
			continue
		}
		queryArgs[key] = value
	}
	return queryArgs
}

func getForwardParams(forwardPrefix string, queryCtx domain.QueryContext) map[string]interface{} {
	r := make(map[string]interface{})
	for k, v := range queryCtx.Input.Params {
		if strings.HasPrefix(k, forwardPrefix) {
			r[k] = v
		}
	}

	return r
}

func parseTimeout(defaultResourceTimeout time.Duration, statement domain.Statement) time.Duration {
	timeout := statement.Timeout
	if timeout == nil {
		return defaultResourceTimeout
	}

	duration, ok := timeout.(int)
	if !ok {
		return defaultResourceTimeout
	}

	return time.Millisecond * time.Duration(duration)
}
