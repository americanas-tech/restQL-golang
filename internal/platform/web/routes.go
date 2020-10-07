package web

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/web/middleware"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"net/http"

	"github.com/b2wdigital/restQL-golang/v4/internal/eval"
	"github.com/b2wdigital/restQL-golang/v4/internal/parser"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/cache"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/httpclient"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/persistence"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/plugins"
	"github.com/b2wdigital/restQL-golang/v4/internal/runner"
	"github.com/valyala/fasthttp"
)

// API constructs a handler for the restQL query related endpoints
func API(log restql.Logger, cfg *conf.Config) (fasthttp.RequestHandler, error) {
	log.Debug("starting api")
	defaultParser, err := parser.New()
	if err != nil {
		log.Error("failed to compile parser", err)
		return nil, err
	}
	parserCacheLoader := cache.New(log, cfg.Cache.Parser.MaxSize, cache.ParserCacheLoader(defaultParser))
	parserCache := cache.NewParserCache(log, parserCacheLoader)

	databaseDisabled := cfg.Plugins.DisableDatabase
	db, err := persistence.NewDatabase(log, databaseDisabled)
	if err != nil {
		log.Error("failed to establish connection to database", err)
		return nil, err
	}

	lifecycle, err := plugins.NewLifecycle(log)
	if err != nil {
		log.Error("failed to initialize plugins", err)
	}

	md := middleware.NewDecorator(log, cfg, lifecycle)

	app := newApp(log, appOptions{MiddlewareDecorator: md})
	client := httpclient.New(log, lifecycle, cfg)
	executor := runner.NewExecutor(log, client, cfg.HTTP.QueryResourceTimeout, cfg.HTTP.ForwardPrefix)
	r := runner.NewRunner(log, executor, cfg.HTTP.GlobalQueryTimeout)

	mr := persistence.NewMappingReader(log, cfg.Env, cfg.Mappings, db)
	tenantCache := cache.New(log, cfg.Cache.Mappings.MaxSize,
		cache.TenantCacheLoader(mr),
		cache.WithExpiration(cfg.Cache.Mappings.Expiration),
		cache.WithRefreshInterval(cfg.Cache.Mappings.RefreshInterval),
		cache.WithRefreshQueueLength(cfg.Cache.Mappings.RefreshQueueLength),
	)
	cacheMr := cache.NewMappingsReaderCache(log, tenantCache)

	qr := persistence.NewQueryReader(log, cfg.Queries, db)
	queryCache := cache.New(log, cfg.Cache.Query.MaxSize, cache.QueryCacheLoader(qr))
	cacheQr := cache.NewQueryReaderCache(log, queryCache)

	e := eval.NewEvaluator(log, cacheMr, cacheQr, r, parserCache, lifecycle)

	restQl := newRestQl(log, cfg, e, defaultParser)

	app.Handle(http.MethodPost, "/validate-query", restQl.ValidateQuery)
	app.Handle(http.MethodPost, "/run-query", restQl.RunAdHocQuery)
	app.Handle(http.MethodGet, "/run-query/{namespace}/{queryId}/{revision}", restQl.RunSavedQuery)
	app.Handle(http.MethodPost, "/run-query/{namespace}/{queryId}/{revision}", restQl.RunSavedQuery)

	return app.RequestHandler(), nil
}

// Health constructs a handler for system checks endpoints
func Health(log restql.Logger, cfg *conf.Config) fasthttp.RequestHandler {
	app := newApp(log, appOptions{})
	check := newCheck(cfg.Build)

	app.Handle(http.MethodGet, "/health", check.Health)
	app.Handle(http.MethodGet, "/resource-status", check.ResourceStatus)

	return app.RequestHandler()
}

// Debug constructs a handler for profiling endpoints
func Debug(log restql.Logger) fasthttp.RequestHandler {
	app := newApp(log, appOptions{})
	d := newDebug()

	app.Handle(http.MethodGet, "/debug/pprof/goroutine", d.Index)
	app.Handle(http.MethodGet, "/debug/pprof/heap", d.Index)
	app.Handle(http.MethodGet, "/debug/pprof/threadcreate", d.Index)
	app.Handle(http.MethodGet, "/debug/pprof/block", d.Index)
	app.Handle(http.MethodGet, "/debug/pprof/mutex", d.Index)

	app.Handle(http.MethodGet, "/debug/pprof/profile", d.Profile)

	return app.RequestHandler()
}
