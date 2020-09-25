package plugins

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

// Lifecycle represent the hooks on the query execution.
type Lifecycle interface {
	BeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context
	AfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context
	BeforeQuery(ctx context.Context, query string, queryCtx restql.QueryContext) context.Context
	AfterQuery(ctx context.Context, query string, result domain.Resources) context.Context
	BeforeRequest(ctx context.Context, request restql.HTTPRequest) context.Context
	AfterRequest(ctx context.Context, request restql.HTTPRequest, response restql.HTTPResponse, err error) context.Context
}

type pluginExecutor func(ctx context.Context, p restql.LifecyclePlugin) context.Context

type manager struct {
	log              restql.Logger
	availablePlugins []restql.LifecyclePlugin
}

// NewLifecycle constructs a Lifecycle instance.
func NewLifecycle(log restql.Logger) (Lifecycle, error) {
	ps := loadLifecyclePlugins(log)
	if len(ps) == 0 {
		log.Info("no lifecycle hook provided")
		return NoOpLifecycle, nil
	}

	return manager{log: log, availablePlugins: ps}, nil
}

func (m manager) BeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return m.executeAllPluginsWithContext(ctx, "BeforeTransaction", func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		log := restql.GetLogger(ctx)
		tr := m.newTransactionRequest(log, requestCtx)
		return p.BeforeTransaction(currentCtx, tr)
	})
}

func (m manager) AfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return m.executeAllPluginsWithContext(ctx, "AfterTransaction", func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		tr := m.newTransactionResponse(requestCtx)
		return p.AfterTransaction(currentCtx, tr)
	})
}

func (m manager) BeforeQuery(ctx context.Context, query string, queryCtx restql.QueryContext) context.Context {
	return m.executeAllPluginsWithContext(ctx, "BeforeQuery", func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		return p.BeforeQuery(currentCtx, query, queryCtx)
	})
}

func (m manager) AfterQuery(ctx context.Context, query string, result domain.Resources) context.Context {
	return m.executeAllPluginsWithContext(ctx, "AfterQuery", func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		r := make(map[string]interface{})
		for id, resource := range result {
			r[string(id)] = resource
		}

		return p.AfterQuery(currentCtx, query, r)
	})
}

func (m manager) BeforeRequest(ctx context.Context, request restql.HTTPRequest) context.Context {
	return m.executeAllPluginsWithContext(ctx, "BeforeRequest", func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		return p.BeforeRequest(currentCtx, request)
	})
}

func (m manager) AfterRequest(ctx context.Context, request restql.HTTPRequest, response restql.HTTPResponse, err error) context.Context {
	return m.executeAllPluginsWithContext(ctx, "AfterRequest", func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		return p.AfterRequest(currentCtx, request, response, err)
	})
}
func (m manager) executeAllPluginsWithContext(ctx context.Context, hook string, fn pluginExecutor) context.Context {
	log := restql.GetLogger(ctx)

	var pluginCtx context.Context

	pluginCtx = ctx
	for _, p := range m.availablePlugins {
		m.safeExecute(log, p.Name(), hook, func() {
			pluginCtx = fn(pluginCtx, p)
		})
	}

	return pluginCtx
}

func (m manager) safeExecute(log restql.Logger, pluginName string, hook string, fn func()) {
	defer func() {
		if reason := recover(); reason != nil {
			err := errors.Errorf("reason : %v\n\t stack : %v", reason, string(debug.Stack()))
			log.Error("plugin produced a panic", err, "name", pluginName, "hook", hook)
		}
	}()

	fn()
}

func (m manager) newTransactionRequest(log restql.Logger, ctx *fasthttp.RequestCtx) restql.TransactionRequest {
	uri, err := url.ParseRequestURI(string(ctx.RequestURI()))
	if err != nil {
		log.Error("failed to parse request uri for plugin", err)
	}

	header := make(http.Header)
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		header.Add(string(k), string(v))
	})

	//todo: add header to ctx

	return restql.TransactionRequest{
		Url:    uri,
		Method: string(ctx.Method()),
		Header: header,
	}
}

func (m manager) newTransactionResponse(ctx *fasthttp.RequestCtx) restql.TransactionResponse {
	header := make(http.Header)
	ctx.Response.Header.VisitAll(func(k, v []byte) {
		header.Add(string(k), string(v))
	})

	return restql.TransactionResponse{
		Status: ctx.Response.StatusCode(),
		Header: header,
		Body:   ctx.Response.Body(),
	}
}

// NoOpLifecycle is Lifecycle implementation with no handlers.
var NoOpLifecycle Lifecycle = noOpLifecycle{}

type noOpLifecycle struct{}

func (n noOpLifecycle) BeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return ctx
}
func (n noOpLifecycle) AfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return ctx
}
func (n noOpLifecycle) BeforeQuery(ctx context.Context, query string, queryCtx restql.QueryContext) context.Context {
	return ctx
}
func (n noOpLifecycle) AfterQuery(ctx context.Context, query string, result domain.Resources) context.Context {
	return ctx
}
func (n noOpLifecycle) BeforeRequest(ctx context.Context, request restql.HTTPRequest) context.Context {
	return ctx
}
func (n noOpLifecycle) AfterRequest(ctx context.Context, request restql.HTTPRequest, response restql.HTTPResponse, err error) context.Context {
	return ctx
}
