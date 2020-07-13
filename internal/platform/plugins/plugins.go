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
	"runtime/debug"
)

type Manager interface {
	RunBeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context
	RunAfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context
	RunBeforeQuery(ctx context.Context, query string, queryCtx domain.QueryContext) context.Context
	RunAfterQuery(ctx context.Context, query string, result domain.Resources) context.Context
	RunBeforeRequest(ctx context.Context, request domain.HttpRequest) context.Context
	RunAfterRequest(ctx context.Context, request domain.HttpRequest, response domain.HttpResponse, err error) context.Context
}

type pluginExecutor func(ctx context.Context, p restql.Plugin) context.Context

type manager struct {
	log              *logger.Logger
	availablePlugins []restql.Plugin
}

func NewManager(log *logger.Logger, pluginsLocation string) (Manager, error) {
	ps := loadStaticPlugin(log)
	if len(ps) == 0 {
		log.Info("no plugins provided")
		return NoOpManager, nil
	}

	return manager{log: log, availablePlugins: ps}, nil
}

func (m manager) RunBeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return m.executeAllPluginsWithContext("BeforeTransaction", ctx, func(currentCtx context.Context, p restql.Plugin) context.Context {
		log := restql.GetLogger(ctx)
		tr := m.newTransactionRequest(log, requestCtx)
		return p.BeforeTransaction(currentCtx, tr)
	})
}

func (m manager) RunAfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return m.executeAllPluginsWithContext("AfterTransaction", ctx, func(currentCtx context.Context, p restql.Plugin) context.Context {
		tr := m.newTransactionResponse(requestCtx)
		return p.AfterTransaction(currentCtx, tr)
	})
}

func (m manager) RunBeforeQuery(ctx context.Context, query string, queryCtx domain.QueryContext) context.Context {
	return m.executeAllPluginsWithContext("BeforeQuery", ctx, func(currentCtx context.Context, p restql.Plugin) context.Context {
		return p.BeforeQuery(currentCtx, query, queryCtx)
	})
}

func (m manager) RunAfterQuery(ctx context.Context, query string, result domain.Resources) context.Context {
	return m.executeAllPluginsWithContext("AfterQuery", ctx, func(currentCtx context.Context, p restql.Plugin) context.Context {
		m := DecodeQueryResult(result)
		return p.AfterQuery(currentCtx, query, m)
	})
}

func (m manager) RunBeforeRequest(ctx context.Context, request domain.HttpRequest) context.Context {
	return m.executeAllPluginsWithContext("BeforeRequest", ctx, func(currentCtx context.Context, p restql.Plugin) context.Context {
		return p.BeforeRequest(currentCtx, request)
	})
}

func (m manager) RunAfterRequest(ctx context.Context, request domain.HttpRequest, response domain.HttpResponse, err error) context.Context {
	return m.executeAllPluginsWithContext("AfterRequest", ctx, func(currentCtx context.Context, p restql.Plugin) context.Context {
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

var NoOpManager Manager = noOpManager{}

type noOpManager struct{}

func (n noOpManager) RunBeforeTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return ctx
}
func (n noOpManager) RunAfterTransaction(ctx context.Context, requestCtx *fasthttp.RequestCtx) context.Context {
	return ctx
}
func (n noOpManager) RunBeforeQuery(ctx context.Context, query string, queryCtx domain.QueryContext) context.Context {
	return ctx
}
func (n noOpManager) RunAfterQuery(ctx context.Context, query string, result domain.Resources) context.Context {
	return ctx
}
func (n noOpManager) RunBeforeRequest(ctx context.Context, request domain.HttpRequest) context.Context {
	return ctx
}
func (n noOpManager) RunAfterRequest(ctx context.Context, request domain.HttpRequest, response domain.HttpResponse, err error) context.Context {
	return ctx
}
