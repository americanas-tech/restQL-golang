package plugins

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
)

type Manager interface {
	RunBeforeQuery(query string, queryCtx domain.QueryContext)
	RunAfterQuery(query string, result domain.Resources)
	RunBeforeRequest(request domain.HttpRequest)
	RunAfterRequest(request domain.HttpRequest, response domain.HttpResponse, err error)
}

type manager struct {
	log              *logger.Logger
	availablePlugins []restql.Plugin
	runner           *Runner
}

func NewManager(log *logger.Logger, pluginsLocation string) (Manager, error) {
	ps, err := loadPlugins(log, pluginsLocation)
	if err != nil {
		return noOpManager{}, err
	}

	runner := NewRunner(log)

	return manager{log: log, availablePlugins: ps, runner: runner}, nil
}

func (m manager) RunBeforeQuery(query string, queryCtx domain.QueryContext) {
	for _, p := range m.availablePlugins {
		m.runner.BeforeQuery(p, query, queryCtx)
	}
}

func (m manager) RunAfterQuery(query string, result domain.Resources) {
	for _, p := range m.availablePlugins {
		m.runner.AfterQuery(p, query, result)
	}
}

func (m manager) RunBeforeRequest(request domain.HttpRequest) {
	for _, p := range m.availablePlugins {
		m.runner.BeforeRequest(p, request)
	}
}

func (m manager) RunAfterRequest(request domain.HttpRequest, response domain.HttpResponse, err error) {
	for _, p := range m.availablePlugins {
		m.runner.AfterRequest(p, request, response, err)
	}
}

type noOpManager struct{}

func (n noOpManager) RunBeforeQuery(query string, queryCtx domain.QueryContext) {}
func (n noOpManager) RunAfterQuery(query string, result domain.Resources)       {}
func (n noOpManager) RunBeforeRequest(request domain.HttpRequest)               {}
func (n noOpManager) RunAfterRequest(request domain.HttpRequest, response domain.HttpResponse, err error) {
}
