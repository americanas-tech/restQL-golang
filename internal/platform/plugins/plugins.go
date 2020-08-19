package plugins

import (
	"context"
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Lifecycle interface {
	BeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context
	AfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context
	BeforeQuery(ctx context.Context, query string, queryCtx domain.QueryContext) context.Context
	AfterQuery(ctx context.Context, query string, result domain.Resources) context.Context
	BeforeRequest(ctx context.Context, request domain.HttpRequest) context.Context
	AfterRequest(ctx context.Context, request domain.HttpRequest, response domain.HttpResponse, err error) context.Context
}

type pluginExecutor func(ctx context.Context, p restql.LifecyclePlugin) context.Context

type manager struct {
	log              *logger.Logger
	availablePlugins []restql.LifecyclePlugin
}

func NewLifecycle(log *logger.Logger) (Lifecycle, error) {
	ps := loadLifecyclePlugins(log)
	if len(ps) == 0 {
		log.Info("no plugins provided")
		return NoOpLifecycle, nil
	}

	return manager{log: log, availablePlugins: ps}, nil
}

func (m manager) BeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return m.executeAllPluginsWithContext("BeforeTransaction", ctx, func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		log := restql.GetLogger(ctx)
		tr := m.newTransactionRequest(log, requestCtx)
		return p.BeforeTransaction(currentCtx, tr)
	})
}

func (m manager) AfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return m.executeAllPluginsWithContext("AfterTransaction", ctx, func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		tr := m.newTransactionResponse(requestCtx)
		return p.AfterTransaction(currentCtx, tr)
	})
}

func (m manager) BeforeQuery(ctx context.Context, query string, queryCtx domain.QueryContext) context.Context {
	return m.executeAllPluginsWithContext("BeforeQuery", ctx, func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		return p.BeforeQuery(currentCtx, query, queryCtx)
	})
}

func (m manager) AfterQuery(ctx context.Context, query string, result domain.Resources) context.Context {
	return m.executeAllPluginsWithContext("AfterQuery", ctx, func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		m := DecodeQueryResult(result)
		return p.AfterQuery(currentCtx, query, m)
	})
}

func (m manager) BeforeRequest(ctx context.Context, request domain.HttpRequest) context.Context {
	return m.executeAllPluginsWithContext("BeforeRequest", ctx, func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		return p.BeforeRequest(currentCtx, request)
	})
}

func (m manager) AfterRequest(ctx context.Context, request domain.HttpRequest, response domain.HttpResponse, err error) context.Context {
	return m.executeAllPluginsWithContext("AfterRequest", ctx, func(currentCtx context.Context, p restql.LifecyclePlugin) context.Context {
		return p.AfterRequest(currentCtx, request, response, err)
	})
}
func (m manager) executeAllPluginsWithContext(hook string, ctx context.Context, fn pluginExecutor) context.Context {
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

var NoOpLifecycle Lifecycle = noOpLifecycle{}

type noOpLifecycle struct{}

func (n noOpLifecycle) BeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return ctx
}
func (n noOpLifecycle) AfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return ctx
}
func (n noOpLifecycle) BeforeQuery(ctx context.Context, query string, queryCtx domain.QueryContext) context.Context {
	return ctx
}
func (n noOpLifecycle) AfterQuery(ctx context.Context, query string, result domain.Resources) context.Context {
	return ctx
}
func (n noOpLifecycle) BeforeRequest(ctx context.Context, request domain.HttpRequest) context.Context {
	return ctx
}
func (n noOpLifecycle) AfterRequest(ctx context.Context, request domain.HttpRequest, response domain.HttpResponse, err error) context.Context {
	return ctx
}
