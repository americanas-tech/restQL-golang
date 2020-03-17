package runner

import (
	"context"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"strconv"
	"strings"
)

const debugParamName = "_debug"

type DoneRequest Response
type DoneRequests []interface{}

type Executor struct {
	client domain.HttpClient
	log    domain.Logger
}

func (e Executor) DoStatement(ctx context.Context, statement domain.Statement, queryCtx QueryContext) DoneRequest {
	request := e.makeRequest(queryCtx.Mappings, statement)
	response, err := e.client.Do(ctx, request)
	if err != nil {
		e.log.Debug("request failed", "error", err)
		return DoneRequest{}
	}

	dr := newDoneRequest(queryCtx, request, response)

	return dr
}

func newDoneRequest(queryCtx QueryContext, request domain.HttpRequest, response domain.HttpResponse) DoneRequest {
	dr := DoneRequest{
		Details: Details{
			Status:  response.StatusCode,
			Success: response.StatusCode >= 200 && response.StatusCode < 400,
		},
		Result: response.Body,
	}

	debug := getDebug(queryCtx)

	if debug {
		dr.Details.Debug = &Debugging{
			Url:             request.Schema + "://" + request.Uri,
			Params:          request.Query,
			RequestHeaders:  request.Headers,
			ResponseHeaders: response.Headers,
		}
	}

	return dr
}

func getDebug(queryCtx QueryContext) bool {
	param, found := queryCtx.Input.Params[debugParamName]
	if !found {
		return false
	}

	debug, ok := param.(string)
	if !ok {
		return false
	}

	d, err := strconv.ParseBool(debug)
	if err != nil {
		return false
	}

	return d
}

func (e Executor) DoMultiplexedStatement(ctx context.Context, statements []interface{}, queryCtx QueryContext) DoneRequests {
	responses := make(DoneRequests, len(statements))

	for i, stmt := range statements {
		switch stmt := stmt.(type) {
		case domain.Statement:
			responses[i] = e.DoStatement(ctx, stmt, queryCtx)
		case []interface{}:
			responses[i] = e.DoMultiplexedStatement(ctx, stmt, queryCtx)
		}
	}

	return responses
}

func (e Executor) makeRequest(mappings map[string]domain.Mapping, statement domain.Statement) domain.HttpRequest {
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

	return domain.HttpRequest{
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
