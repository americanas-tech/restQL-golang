package runner

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
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

func MakeRequest(statement domain.Statement, queryCtx domain.QueryContext) domain.HttpRequest {
	mapping := queryCtx.Mappings[statement.Resource]
	url := makeUrl(mapping, statement)

	queryParams := makeQueryParams(statement, mapping, queryCtx)

	headers := makeHeaders(statement, queryCtx)

	return domain.HttpRequest{
		Schema:  mapping.Schema,
		Uri:     url,
		Query:   queryParams,
		Body:    nil,
		Headers: headers,
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

		str, ok := value.(string)

		if !ok {
			continue
		}
		queryArgs[key] = str
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
		resource = strings.Replace(resource, fmt.Sprintf(":%v", pathParam), fmt.Sprintf("%v", statement.With[pathParam]), 1)
	}

	return resource
}
