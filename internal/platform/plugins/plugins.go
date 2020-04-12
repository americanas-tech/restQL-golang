package plugins

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
)

type Manager interface {
	RunBeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context
	RunAfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx)
	RunBeforeQuery(ctx context.Context, query string, queryCtx domain.QueryContext)
	RunAfterQuery(ctx context.Context, query string, result domain.Resources)
	RunBeforeRequest(ctx context.Context, request domain.HttpRequest)
	RunAfterRequest(ctx context.Context, request domain.HttpRequest, response domain.HttpResponse, err error)
}

type manager struct {
	log              *logger.Logger
	availablePlugins []restql.Plugin
}

func NewManager(log *logger.Logger, pluginsLocation string) (Manager, error) {
	ps, err := loadPlugins(log, pluginsLocation)
	if err != nil {
		return NoOpManager, err
	}

	return manager{log: log, availablePlugins: ps}, nil
}

func (m manager) RunBeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	var pluginCtx context.Context

	pluginCtx = ctx
	for _, p := range m.availablePlugins {
		m.safeExecute(p.Name(), "BeforeTransaction", func() {
			tr := m.newTransactionRequest(requestCtx)
			pluginCtx = p.BeforeTransaction(pluginCtx, tr)
		})
	}

	return pluginCtx
}

func (m manager) RunAfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) {
	for _, p := range m.availablePlugins {
		m.safeExecute(p.Name(), "AfterTransaction", func() {
			tr := m.newTransactionResponse(requestCtx)
			p.AfterTransaction(ctx, tr)
		})
	}
}

func (m manager) RunBeforeQuery(ctx context.Context, query string, queryCtx domain.QueryContext) {
	for _, p := range m.availablePlugins {
		m.safeExecute(p.Name(), "BeforeQuery", func() {
			p.BeforeQuery(ctx, query, queryCtx)
		})
	}
}

func (m manager) RunAfterQuery(ctx context.Context, query string, result domain.Resources) {
	for _, p := range m.availablePlugins {
		m.safeExecute(p.Name(), "AfterQuery", func() {
			m := convertQueryResult(result)
			p.AfterQuery(ctx, query, m)
		})
	}
}

func (m manager) RunBeforeRequest(ctx context.Context, request domain.HttpRequest) {
	for _, p := range m.availablePlugins {
		m.safeExecute(p.Name(), "BeforeRequest", func() {
			p.BeforeRequest(ctx, request)
		})
	}
}

func (m manager) RunAfterRequest(ctx context.Context, request domain.HttpRequest, response domain.HttpResponse, err error) {
	for _, p := range m.availablePlugins {
		m.safeExecute(p.Name(), "AfterRequest", func() {
			p.AfterRequest(ctx, request, response, err)
		})
	}
}

func (m manager) safeExecute(pluginName string, hook string, fn func()) {
	defer func() {
		if reason := recover(); reason != nil {
			err := errors.Errorf("reason : %v", reason)
			m.log.Error("plugin produced a panic", err, "name", pluginName, "hook", hook)
		}
	}()

	fn()
}

func (m manager) newTransactionRequest(ctx *fasthttp.RequestCtx) restql.TransactionRequest {
	uri, err := url.ParseRequestURI(string(ctx.RequestURI()))
	if err != nil {
		m.log.Error("failed to parse request uri for plugin", err)
	}

	header := make(http.Header)
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		header.Set(string(k), string(v))
	})

	//todo: add header to ctx

	return restql.TransactionRequest{
		Url:    uri,
		Method: string(ctx.Method()),
		Header: header,
	}
}

func (m manager) newTransactionResponse(ctx *fasthttp.RequestCtx) restql.TransactionResponse {
	return restql.TransactionResponse{
		Status: ctx.Response.StatusCode(),
		Header: ctx.Response.Header.Header(),
		Body:   ctx.Response.Body(),
	}
}

func convertQueryResult(resource interface{}) map[string]interface{} {
	switch resource := resource.(type) {
	case domain.Resources:
		m := make(map[string]interface{})
		for k, v := range resource {
			m[string(k)] = convertDoneResource(v)
		}
		return m
	case domain.Details:
		return map[string]interface{}{
			"status":       resource.Status,
			"success":      resource.Success,
			"ignoreErrors": resource.IgnoreErrors,
			"debugging":    convertQueryResult(resource.Debug),
		}
	case *domain.Debugging:
		return map[string]interface{}{
			"method":          resource.Method,
			"url":             resource.Url,
			"requestHeaders":  resource.RequestHeaders,
			"responseHeaders": resource.ResponseHeaders,
			"params":          resource.Params,
			"requestBody":     resource.RequestBody,
			"responseTime":    resource.ResponseTime,
		}
	default:
		return nil
	}
}

func convertDoneResource(doneResource interface{}) interface{} {
	switch resource := doneResource.(type) {
	case domain.DoneResource:
		return map[string]interface{}{
			"details": convertQueryResult(resource.Details),
			"result":  resource.Result,
		}
	case domain.DoneResources:
		l := make([]interface{}, len(resource))
		for i, r := range resource {
			l[i] = convertQueryResult(r)
		}
		return l
	default:
		return resource
	}
}

var NoOpManager Manager = noOpManager{}

type noOpManager struct{}

func (n noOpManager) RunBeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return nil
}
func (n noOpManager) RunAfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx)       {}
func (n noOpManager) RunBeforeQuery(ctx context.Context, query string, queryCtx domain.QueryContext) {}
func (n noOpManager) RunAfterQuery(ctx context.Context, query string, result domain.Resources)       {}
func (n noOpManager) RunBeforeRequest(ctx context.Context, request domain.HttpRequest)               {}
func (n noOpManager) RunAfterRequest(ctx context.Context, request domain.HttpRequest, response domain.HttpResponse, err error) {
}
