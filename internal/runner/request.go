package runner

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/b2wdigital/restQL-golang/v6/internal/domain"
	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
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
	timeout := parseTimeout(defaultResourceTimeout, statement)

	queryParams := makeQueryParams(forwardPrefix, statement, mapping, queryCtx)

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
	for key, value := range getValueForBody(statement, mapping) {
		if value, ok := value.(domain.AsBody); ok {
			return parseBodyValue(value.Target())
		}

		if !isPrimitiveValue(value) {
			continue
		}

		result[key] = parseBodyValue(value)
	}
	return result
}

func getValueForBody(statement domain.Statement, mapping restql.Mapping) map[string]interface{} {
	values := make(map[string]interface{})

	for key, value := range statement.With.Values {
		if mapping.IsPathParam(key) || mapping.IsQueryParam(key) {
			continue
		}

		if _, ok := value.(domain.AsQuery); ok {
			continue
		}

		values[key] = value
	}

	return values
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

func makeQueryParams(forwardPrefix string, statement domain.Statement, mapping restql.Mapping, queryCtx restql.QueryContext) map[string]interface{} {
	queryArgs := getForwardParams(forwardPrefix, queryCtx)

	for key, value := range mapping.QueryWithParams(statement.With.Values) {
		queryArgs[key] = value
	}

	for key, value := range getValueForQueryParams(statement, mapping) {
		if !isPrimitiveValue(value) {
			continue
		}

		queryArgs[key] = value
	}

	return queryArgs
}

func getValueForQueryParams(statement domain.Statement, mapping restql.Mapping) map[string]interface{} {
	values := make(map[string]interface{})

	for key, value := range statement.With.Values {
		if mapping.IsPathParam(key) {
			continue
		}

		if value, ok := value.(domain.AsQuery); ok {
			values[key] = value.Target()
			continue
		}

		if statement.Method == domain.FromMethod || statement.Method == domain.DeleteMethod {
			values[key] = value
			continue
		}
	}

	return values
}

func isPrimitiveValue(value interface{}) bool {
	if value == nil {
		return false
	}

	switch value.(type) {
	case string:
		return true
	case bool:
		return true
	case int:
		return true
	case float64:
		return true
	case map[string]interface{}:
		return true
	case []interface{}:
		return true
	default:
		return false
	}
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
