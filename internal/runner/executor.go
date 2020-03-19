package runner

import (
	"bytes"
	"context"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

var errNoTimeoutProvided = errors.New("no timeout provided")

var disallowedHeaders = map[string]struct{}{
	"Host":            {},
	"Content-Type":    {},
	"Content-Length":  {},
	"Connection":      {},
	"Origin":          {},
	"Accept-Encoding": {},
}

const debugParamName = "_debug"

type DoneRequest Response
type DoneRequests []interface{}

type Executor struct {
	client domain.HttpClient
	log    domain.Logger
}

func (e Executor) DoStatement(ctx context.Context, statement domain.Statement, queryCtx QueryContext) (DoneRequest, error) {
	emptyChainedParams := getEmptyChainedParams(statement)
	if len(emptyChainedParams) > 0 {
		return newEmptyChainedResponse(emptyChainedParams), nil
	}

	e.log.Debug("executing request for statement", "resource", statement.Resource, "method", statement.Method)

	request := e.makeRequest(statement, queryCtx)

	timeout, err := parseTimeout(statement)
	if err == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	} else if err != errNoTimeoutProvided {
		e.log.Debug("failed to set timeout for statement", "error", err)
	}

	response, err := e.client.Do(ctx, request)
	switch {
	case err == domain.ErrRequestTimeout:
		return newTimeoutResponse(err, request, response, queryCtx), nil
	case err != nil:
		e.log.Debug("request failed", "error", err)
		return DoneRequest{}, err
	}

	dr := newDoneRequest(request, response, queryCtx)

	e.log.Debug("request execution done", "resource", statement.Resource, "method", statement.Method, "response", dr)

	return dr, nil
}

func (e Executor) DoMultiplexedStatement(ctx context.Context, statements []interface{}, queryCtx QueryContext) (DoneRequests, error) {
	responses := make(DoneRequests, len(statements))

	for i, stmt := range statements {
		switch stmt := stmt.(type) {
		case domain.Statement:
			r, err := e.DoStatement(ctx, stmt, queryCtx)
			if err != nil {
				return nil, err
			}
			responses[i] = r
		case []interface{}:
			r, err := e.DoMultiplexedStatement(ctx, stmt, queryCtx)
			if err != nil {
				return nil, err
			}
			responses[i] = r
		}
	}

	return responses, nil
}

func newDoneRequest(request domain.HttpRequest, response domain.HttpResponse, queryCtx QueryContext) DoneRequest {
	dr := DoneRequest{
		Details: Details{
			Status:  response.StatusCode,
			Success: response.StatusCode >= 200 && response.StatusCode < 400,
		},
		Result: response.Body,
	}

	debug := getDebug(queryCtx)

	if debug {
		dr.Details.Debug = newDebugging(request, response)
	}

	return dr
}

func newDebugging(request domain.HttpRequest, response domain.HttpResponse) *Debugging {
	return &Debugging{
		Url:             response.Url,
		Params:          request.Query,
		RequestHeaders:  request.Headers,
		ResponseHeaders: response.Headers,
		ResponseTime:    response.Duration.Milliseconds(),
	}
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

func (e Executor) makeRequest(statement domain.Statement, queryCtx QueryContext) domain.HttpRequest {
	mapping := queryCtx.Mappings[statement.Resource]
	url := makeUrl(mapping, statement)

	queryParams := e.makeQueryParams(statement, mapping, queryCtx)

	headers := e.makeHeaders(statement, queryCtx)

	return domain.HttpRequest{
		Schema:  mapping.Schema,
		Uri:     url,
		Query:   queryParams,
		Body:    nil,
		Headers: headers,
	}
}

func (e Executor) makeHeaders(statement domain.Statement, queryCtx QueryContext) map[string]string {
	headers := getForwardHeaders(queryCtx)
	for key, value := range statement.Headers {
		str, ok := value.(string)
		if !ok {
			e.log.Debug("skipping header on request build for failing string casting", "header-name", key, "header-value", value)
			continue
		}
		headers[key] = str
	}
	return headers
}

func getForwardHeaders(queryCtx QueryContext) map[string]string {
	r := make(map[string]string)
	for k, v := range queryCtx.Input.Headers {
		if _, found := disallowedHeaders[k]; !found {
			r[k] = v
		}
	}
	return r
}

func (e Executor) makeQueryParams(statement domain.Statement, mapping domain.Mapping, queryCtx QueryContext) map[string]interface{} {
	queryArgs := getForwardParams(queryCtx)
	for key, value := range statement.With {
		if mapping.HasParam(key) {
			continue
		}

		str, ok := value.(string)

		if !ok {
			e.log.Debug("skipping resources param on request build for failing string casting", "param-key", key, "param-value", value)
			continue
		}
		queryArgs[key] = str
	}
	return queryArgs
}

func getForwardParams(queryCtx QueryContext) map[string]interface{} {
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

func parseTimeout(statement domain.Statement) (time.Duration, error) {
	timeout := statement.Timeout
	if timeout == nil {
		return 0, errNoTimeoutProvided
	}

	duration, ok := timeout.(int)
	if !ok {
		return 0, errors.Errorf("statement timeout is not an int, got %T", timeout)
	}

	return time.Millisecond * time.Duration(duration), nil
}

func newTimeoutResponse(err error, request domain.HttpRequest, response domain.HttpResponse, queryCtx QueryContext) DoneRequest {
	dr := DoneRequest{
		Details: Details{
			Status:  408,
			Success: false,
		},
		Result: err.Error(),
	}

	debug := getDebug(queryCtx)

	if debug {
		dr.Details.Debug = newDebugging(request, response)
	}

	return dr
}

func newEmptyChainedResponse(params []string) DoneRequest {
	var buf bytes.Buffer

	buf.WriteString("The request was skipped due to missing { ")
	for _, p := range params {
		buf.WriteString(":")
		buf.WriteString(p)
		buf.WriteString(" ")
	}
	buf.WriteString(" } param value")

	return DoneRequest{Details: Details{Status: 400, Success: false}, Result: buf.String()}
}

func getEmptyChainedParams(statement domain.Statement) []string {
	var r []string
	for key, value := range statement.With {
		if value == EmptyChained {
			r = append(r, key)
		}
	}

	return r
}
