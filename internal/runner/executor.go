package runner

import (
	"context"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"strings"
)

type DoneRequest domain.Response
type DoneRequests []interface{}

type Executor struct {
	mappings map[string]domain.Mapping
	client   domain.HttpClient
	log      domain.Logger
}

func (e Executor) DoStatement(ctx context.Context, statement domain.Statement) DoneRequest {
	req := e.makeRequest(e.mappings, statement)
	response, err := e.client.Do(ctx, req)
	if err != nil {
		e.log.Debug("request failed", "error", err)
		return DoneRequest{}
	}

	return DoneRequest(response)
}

func (e Executor) DoMultiplexedStatement(ctx context.Context, statements []interface{}) DoneRequests {
	responses := make(DoneRequests, len(statements))

	for i, stmt := range statements {
		switch stmt := stmt.(type) {
		case domain.Statement:
			responses[i] = e.DoStatement(ctx, stmt)
		case []interface{}:
			responses[i] = e.DoMultiplexedStatement(ctx, stmt)
		}
	}

	return responses
}

func (e Executor) makeRequest(mappings map[string]domain.Mapping, statement domain.Statement) domain.Request {
	mapping := mappings[statement.Resource]
	url := makeUrl(mapping, statement)

	queryArgs := make(map[string]string)
	for key, value := range statement.With {
		if contains(mapping.PathParams, key) {
			continue
		}

		str, ok := value.(string)

		if !ok {
			e.log.Debug("skipping resources param on request build for failing string casting", "param-key", key, "param-value", value)
			continue
		}
		queryArgs[key] = str
	}

	headers := make(map[string]string)
	for key, value := range statement.Headers {
		str, ok := value.(string)
		if !ok {
			e.log.Debug("skipping header on request build for failing string casting", "header-name", key, "header-value", value)
			continue
		}
		headers[key] = str
	}

	return domain.Request{
		Schema:  mapping.Schema,
		Uri:     url,
		Query:   queryArgs,
		Body:    nil,
		Headers: headers,
	}
}

func makeUrl(mapping domain.Mapping, statement domain.Statement) string {
	resource := mapping.Uri
	for _, pathParam := range mapping.PathParams {
		resource = strings.Replace(resource, fmt.Sprintf(":%v", pathParam), fmt.Sprintf("%v", statement.With[pathParam]), 1)
	}

	return resource
}

func contains(list []string, item string) bool {
	for _, el := range list {
		if el == item {
			return true
		}
	}

	return false
}
