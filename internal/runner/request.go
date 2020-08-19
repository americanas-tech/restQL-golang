package runner

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"net/http"
	"strings"
	"time"
)

var disallowedHeaders = map[string]struct{}{
	"host":            {},
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
	path := mapping.PathWithParams(statement.With.Values)
	queryParams := makeQueryParams(forwardPrefix, statement, mapping, queryCtx)
	timeout := parseTimeout(defaultResourceTimeout, statement)

	req := domain.HttpRequest{
		Method:  method,
		Schema:  mapping.Scheme(),
		Host:    mapping.Host(),
		Path:    path,
		Query:   queryParams,
		Headers: headers,
		Timeout: timeout,
	}

	if statement.Method == domain.ToMethod || statement.Method == domain.UpdateMethod || statement.Method == domain.IntoMethod {
		req.Body = makeBody(statement, mapping)
	}

	return req
}

func makeBody(statement domain.Statement, mapping domain.Mapping) domain.Body {
	if statement.With.Body != nil {
		return statement.With.Body
	}

	result := make(map[string]interface{})
	for key, value := range statement.With.Values {
		if mapping.IsPathParam(key) || mapping.IsQueryParam(key) {
			continue
		}

		if value, ok := value.(domain.AsBody); ok {
			return parseBodyValue(value.Target())
		}

		result[key] = parseBodyValue(value)
	}
	return result
}

func parseBodyValue(value interface{}) interface{} {
	switch value := value.(type) {
	case string:
		valid := json.Valid([]byte(value))
		if !valid {
			return value
		}

		var m interface{}
		_ = json.Unmarshal([]byte(value), &m)
		return m
	default:
		return value
	}
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

	_, found := headers["Content-Type"]
	if !found {
		headers["Content-Type"] = "application/json"
	}

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

	for key, value := range mapping.QueryWithParams(statement.With.Values) {
		queryArgs[key] = value
	}

	if statement.Method == domain.FromMethod || statement.Method == domain.DeleteMethod {
		for key, value := range statement.With.Values {
			if mapping.IsPathParam(key) {
				continue
			}
			queryArgs[key] = value
		}
	}

	return queryArgs
}

func getForwardParams(forwardPrefix string, queryCtx domain.QueryContext) map[string]interface{} {
	r := make(map[string]interface{})
	if forwardPrefix == "" {
		return r
	}

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
