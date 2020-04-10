package web

import (
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"github.com/b2wdigital/restQL-golang/internal/parser"
	"github.com/b2wdigital/restQL-golang/internal/platform/cache"
	"github.com/b2wdigital/restQL-golang/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/internal/platform/httpclient"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/internal/platform/persistence"
	"github.com/b2wdigital/restQL-golang/internal/platform/persistence/database"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"github.com/valyala/fasthttp"
	"net/http"
)

func API(log *logger.Logger, cfg *conf.Config) (fasthttp.RequestHandler, error) {
	defaultParser, err := parser.New()
	if err != nil {
		log.Error("failed to compile parser", err)
		return nil, err
	}
	parserCacheLoader := cache.New(log, cfg.Cache.Parser.MaxSize, cache.ParserCacheLoader(defaultParser))
	parserCache := cache.NewParserCache(log, parserCacheLoader)

	db, err := database.New(log, cfg.Database.ConnectionString,
		database.WithConnectionTimeout(cfg.Database.Timeouts.Connection),
		database.WithMappingsTimeout(cfg.Database.Timeouts.Mappings),
		database.WithQueryTimeout(cfg.Database.Timeouts.Query),
	)
	if err != nil {
		log.Error("failed to establish connection to database", err)
	}

	app := NewApp(log, cfg)
	client := httpclient.New(log, cfg)

	executor := runner.NewExecutor(log, client, cfg.QueryResourceTimeout)
	r := runner.NewRunner(log, executor, cfg.GlobalQueryTimeout)

	mr := persistence.NewMappingReader(log, cfg.Env, cfg.Mappings, db)
	tenantCache := cache.New(log, cfg.Cache.Mappings.MaxSize,
		cache.TenantCacheLoader(mr),
		cache.WithExpiration(cfg.Cache.Mappings.Expiration),
		cache.WithRefreshInterval(cfg.Cache.Mappings.RefreshInterval),
		cache.WithRefreshQueueLength(cfg.Cache.Mappings.RefreshQueueLength),
	)
	cacheMr := cache.NewMappingsReaderCache(log, mr, tenantCache)

	qr := persistence.NewQueryReader(log, cfg.Queries, db)
	queryCache := cache.New(log, cfg.Cache.Query.MaxSize, cache.QueryCacheLoader(qr))
	cacheQr := cache.NewQueryReaderCache(log, qr, queryCache)

	e := eval.NewEvaluator(log, cacheMr, cacheQr, r, parserCache)

	restQl := NewRestQl(log, cfg, e, defaultParser)

	app.Handle(http.MethodPost, "/validate-query", restQl.ValidateQuery)
	app.Handle(http.MethodGet, "/run-query/:namespace/:queryId/:revision", restQl.RunSavedQuery)

	return app.RequestHandler(), nil
}

func Health(log *logger.Logger, cfg *conf.Config) fasthttp.RequestHandler {
	app := NewApp(log, cfg)
	check := NewCheck(cfg.Build)

	app.Handle(http.MethodGet, "/health", check.Health)
	app.Handle(http.MethodGet, "/resource-status", check.ResourceStatus)

	return app.RequestHandlerWithoutMiddlewares()
}

func Debug(log *logger.Logger, cfg *conf.Config) fasthttp.RequestHandler {
	app := NewApp(log, cfg)
	pprof := NewPprof()

	app.Handle(http.MethodGet, "/debug/pprof/goroutine", pprof.Index)
	app.Handle(http.MethodGet, "/debug/pprof/heap", pprof.Index)
	app.Handle(http.MethodGet, "/debug/pprof/threadcreate", pprof.Index)
	app.Handle(http.MethodGet, "/debug/pprof/block", pprof.Index)
	app.Handle(http.MethodGet, "/debug/pprof/mutex", pprof.Index)

	app.Handle(http.MethodGet, "/debug/pprof/profile", pprof.Profile)

	return app.RequestHandlerWithoutMiddlewares()
}
