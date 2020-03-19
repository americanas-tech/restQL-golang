package runner

import (
	"bytes"
	"context"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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

type Executor struct {
	client domain.HttpClient
	log    domain.Logger
}

func (e Executor) DoStatement(ctx context.Context, statement domain.Statement, queryCtx domain.QueryContext) (domain.DoneResource, error) {
	ignoreErrors := statement.IgnoreErrors
	debug := isDebugEnabled(queryCtx)
	drOptions := doneResourceOptions{Debugging: debug, IgnoreErrors: ignoreErrors}

	emptyChainedParams := getEmptyChainedParams(statement)
	if len(emptyChainedParams) > 0 {
		emptyChainedResponse := newEmptyChainedResponse(emptyChainedParams, drOptions)
		e.log.Debug("request execution skipped due to empty chained parameters", "resource", statement.Resource, "method", statement.Method)
		return emptyChainedResponse, nil
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
		return newTimeoutResponse(err, request, response, drOptions), nil
	case err != nil:
		e.log.Debug("request failed", "error", err)
		return domain.DoneResource{}, err
	}

	dr := newDoneResource(request, response, drOptions)

	e.log.Debug("request execution done", "resource", statement.Resource, "method", statement.Method, "response", dr)

	return dr, nil
}

func (e Executor) DoMultiplexedStatement(ctx context.Context, statements []interface{}, queryCtx domain.QueryContext) (domain.DoneResources, error) {
	responseChans := make([]chan interface{}, len(statements))
	for i := range responseChans {
		responseChans[i] = make(chan interface{}, 1)
	}
	defer func() {
		for _, ch := range responseChans {
			close(ch)
		}
	}()

	var g errgroup.Group

	for i, stmt := range statements {
		i, stmt := i, stmt
		ch := responseChans[i]

		g.Go(func() error {
			response, err := e.doCurrentStatement(stmt, ctx, queryCtx)
			if err != nil {
				return err
			}
			ch <- response
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	responses := make(domain.DoneResources, len(statements))
	for i, ch := range responseChans {
		responses[i] = <-ch
	}

	return responses, nil
}

func (e Executor) doCurrentStatement(stmt interface{}, ctx context.Context, queryCtx domain.QueryContext) (interface{}, error) {
	switch stmt := stmt.(type) {
	case domain.Statement:
		r, err := e.DoStatement(ctx, stmt, queryCtx)
		if err != nil {
			return nil, err
		}
		return r, nil
	case []interface{}:
		r, err := e.DoMultiplexedStatement(ctx, stmt, queryCtx)
		if err != nil {
			return nil, err
		}
		return r, nil
	default:
		return nil, errors.Errorf("unknown statement type: %T", stmt)
	}
}

type doneResourceOptions struct {
	Debugging    bool
	IgnoreErrors bool
}

func newDoneResource(request domain.HttpRequest, response domain.HttpResponse, options doneResourceOptions) domain.DoneResource {
	dr := domain.DoneResource{
		Details: domain.Details{
			Status:       response.StatusCode,
			Success:      response.StatusCode >= 200 && response.StatusCode < 400,
			IgnoreErrors: options.IgnoreErrors,
		},
		Result: response.Body,
	}

	if options.Debugging {
		dr.Details.Debug = newDebugging(request, response)
	}

	return dr
}

func newDebugging(request domain.HttpRequest, response domain.HttpResponse) *domain.Debugging {
	return &domain.Debugging{
		Url:             response.Url,
		Params:          request.Query,
		RequestHeaders:  request.Headers,
		ResponseHeaders: response.Headers,
		ResponseTime:    response.Duration.Milliseconds(),
	}
}

func isDebugEnabled(queryCtx domain.QueryContext) bool {
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

func (e Executor) makeRequest(statement domain.Statement, queryCtx domain.QueryContext) domain.HttpRequest {
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

func (e Executor) makeHeaders(statement domain.Statement, queryCtx domain.QueryContext) map[string]string {
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

func getForwardHeaders(queryCtx domain.QueryContext) map[string]string {
	r := make(map[string]string)
	for k, v := range queryCtx.Input.Headers {
		if _, found := disallowedHeaders[k]; !found {
			r[k] = v
		}
	}
	return r
}

func (e Executor) makeQueryParams(statement domain.Statement, mapping domain.Mapping, queryCtx domain.QueryContext) map[string]interface{} {
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

func newTimeoutResponse(err error, request domain.HttpRequest, response domain.HttpResponse, options doneResourceOptions) domain.DoneResource {
	dr := domain.DoneResource{
		Details: domain.Details{
			Status:       408,
			Success:      false,
			IgnoreErrors: options.IgnoreErrors,
		},
		Result: err.Error(),
	}

	if options.Debugging {
		dr.Details.Debug = newDebugging(request, response)
	}

	return dr
}

func newEmptyChainedResponse(params []string, options doneResourceOptions) domain.DoneResource {
	var buf bytes.Buffer

	buf.WriteString("The request was skipped due to missing { ")
	for _, p := range params {
		buf.WriteString(":")
		buf.WriteString(p)
		buf.WriteString(" ")
	}
	buf.WriteString("} param value")

	return domain.DoneResource{
		Details: domain.Details{Status: 400, Success: false, IgnoreErrors: options.IgnoreErrors},
		Result:  buf.String(),
	}
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
