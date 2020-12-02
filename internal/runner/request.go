package runner

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
)

var disallowedHeaders = []string{
	"host",
	"content-type",
	"content-length",
	"connection",
	"origin",
	"accept-encoding",
}

var queryMethodToHTTPMethod = map[string]string{
	domain.FromMethod:   http.MethodGet,
	domain.ToMethod:     http.MethodPost,
	domain.IntoMethod:   http.MethodPut,
	domain.UpdateMethod: http.MethodPatch,
	domain.DeleteMethod: http.MethodDelete,
}

// MakeRequest builds a HTTPRequest from a statement.
func MakeRequest(defaultResourceTimeout time.Duration, forwardPrefix string, statement domain.Statement, queryCtx restql.QueryContext) restql.HTTPRequest {
	mapping := queryCtx.Mappings[statement.Resource]
	method := queryMethodToHTTPMethod[statement.Method]
	headers := makeHeaders(statement, queryCtx)
	path := mapping.PathWithParams(statement.With.Values)
	queryParams := makeQueryParams(forwardPrefix, statement, mapping, queryCtx)
	timeout := parseTimeout(defaultResourceTimeout, statement)

	req := restql.HTTPRequest{
		Method:  method,
		Schema:  mapping.Schema(),
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

func makeBody(statement domain.Statement, mapping restql.Mapping) restql.Body {
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

func makeHeaders(statement domain.Statement, queryCtx restql.QueryContext) map[string]string {
	headers := getForwardHeaders(queryCtx)
	for key, value := range statement.Headers {
		str, ok := value.(string)
		if !ok {
			continue
		}
		key = http.CanonicalHeaderKey(key)
		headers[key] = str
	}

	_, found := headers["Content-Type"]
	if !found {
		headers["Content-Type"] = "application/json"
	}

	return headers
}

func isDisallowedHeader(header string) bool {
	for _, disallowedHeader := range disallowedHeaders {
		if strings.EqualFold(header, disallowedHeader) {
			return true
		}
	}
	return false
}

func getForwardHeaders(queryCtx restql.QueryContext) map[string]string {
	r := make(map[string]string)
	for k, v := range queryCtx.Input.Headers {
		if !isDisallowedHeader(k) {
			k = http.CanonicalHeaderKey(k)
			r[k] = v
		}
	}
	return r
}

func makeQueryParams(forwardPrefix string, statement domain.Statement, mapping restql.Mapping, queryCtx restql.QueryContext) map[string]interface{} {
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

func getForwardParams(forwardPrefix string, queryCtx restql.QueryContext) map[string]interface{} {
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
