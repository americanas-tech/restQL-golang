package runner

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"net/http"
	"strings"
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

func MakeRequest(statement domain.Statement, queryCtx domain.QueryContext) domain.HttpRequest {
	mapping := queryCtx.Mappings[statement.Resource]
	method := queryMethodToHttpMethod[statement.Method]
	url := makeUrl(mapping, statement)
	headers := makeHeaders(statement, queryCtx)

	req := domain.HttpRequest{
		Method:  method,
		Schema:  mapping.Schema,
		Uri:     url,
		Headers: headers,
	}

	if statement.Method == domain.FromMethod {
		req.Query = makeQueryParams(statement, mapping, queryCtx)
	} else {
		req.Query = getForwardParams(queryCtx)
	}

	if statement.Method == domain.ToMethod || statement.Method == domain.UpdateMethod {
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

func makeQueryParams(statement domain.Statement, mapping domain.Mapping, queryCtx domain.QueryContext) map[string]interface{} {
	queryArgs := getForwardParams(queryCtx)
	for key, value := range statement.With {
		if mapping.HasParam(key) {
			continue
		}
		queryArgs[key] = value
	}
	return queryArgs
}

func getForwardParams(queryCtx domain.QueryContext) map[string]interface{} {
	r := make(map[string]interface{})
	for k, v := range queryCtx.Input.Params {
		if strings.HasPrefix(k, "c_") {
			r[k] = v
		}
	}

	return r
}

func makeUrl(mapping domain.Mapping, statement domain.Statement) string {
	resource := mapping.Uri
	for _, pathParam := range mapping.PathParams {
		pathParamValue, found := statement.With[pathParam]
		if !found {
			pathParamValue = ""
		}

		resource = strings.Replace(resource, fmt.Sprintf(":%v", pathParam), fmt.Sprintf("%v", pathParamValue), 1)
	}

	return resource
}
