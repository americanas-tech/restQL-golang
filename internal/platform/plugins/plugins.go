package plugins

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
	"github.com/pkg/errors"
	"os"
	"path"
	"plugin"
)

type Manager struct {
	log              *logger.Logger
	availablePlugins []restql.Plugin
	runner           Runner
}

func NewManager(log *logger.Logger, pluginsLocation string) (Manager, error) {
	ps, err := loadPlugins(log, pluginsLocation)
	if err != nil {
		return Manager{}, err
	}

	return Manager{log: log, availablePlugins: ps, runner: Runner{}}, nil
}

func (m Manager) RunBeforeQuery(query string, queryCtx domain.QueryContext) {
	for _, p := range m.availablePlugins {
		m.runner.BeforeQuery(p, query, queryCtx)
	}
}

func (m Manager) RunAfterQuery(query string, result domain.Resources) {
	for _, p := range m.availablePlugins {
		m.runner.AfterQuery(p, query, result)
	}
}

func (m Manager) RunBeforeRequest(request domain.HttpRequest) {
	for _, p := range m.availablePlugins {
		m.runner.BeforeRequest(p, request)
	}
}

func (m Manager) RunAfterRequest(response domain.HttpResponse, err error) {
	for _, p := range m.availablePlugins {
		m.runner.AfterRequest(p, response, err)
	}
}

func loadPlugins(log *logger.Logger, location string) ([]restql.Plugin, error) {
	dir, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := dir.Close()
		if closeErr != nil {
			log.Error("failed to access plugin directory", closeErr)
		}
	}()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var availablePlugins []restql.Plugin
	for _, info := range fileInfos {
		if !info.IsDir() {
			pluginPath := path.Join(location, info.Name())
			p, err := loadPlugin(pluginPath)
			if err != nil {
				log.Error("failed to load plugin", err, "path", pluginPath)
				continue
			}

			log.Info("plugin loaded", "plugin", info.Name())
			availablePlugins = append(availablePlugins, p)
		}
	}

	return availablePlugins, nil
}

func loadPlugin(pluginPath string) (restql.Plugin, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	addPluginSym, err := p.Lookup("AddPlugin")
	if err != nil {
		return nil, err
	}

	fmt.Printf("addPluginSym type : %T\n", addPluginSym)
	addPlugin, ok := addPluginSym.(func() restql.Plugin)
	if !ok {
		return nil, errors.New("failed to load plugin : AddPlugin function has wrong signature")
	}

	return addPlugin(), nil
}
