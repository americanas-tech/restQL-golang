package plugins

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
	"github.com/pkg/errors"
	"os"
	"path"
	"plugin"
)

func loadPlugins(log *logger.Logger, location string) ([]restql.Plugin, error) {
	if location == "" {
		log.Info("no plugin location provided")
		return nil, nil
	}

	if _, err := os.Stat(location); os.IsNotExist(err) {
		log.Info("provided plugin location does not exist", "path", location)
		return nil, nil
	}

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
			p, err := loadPlugin(log, pluginPath)
			if err != nil {
				log.Error("failed to load plugin", err, "path", pluginPath)
				continue
			}

			log.Info("plugin loaded", "plugin", p.Name())
			availablePlugins = append(availablePlugins, p)
		}
	}

	return availablePlugins, nil
}

func loadPlugin(log *logger.Logger, pluginPath string) (restql.Plugin, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	addPluginSym, err := p.Lookup("AddPlugin")
	if err != nil {
		return nil, err
	}

	fmt.Printf("addPluginSym type : %T\n", addPluginSym)
	addPlugin, ok := addPluginSym.(func(log restql.Logger) restql.Plugin)
	if !ok {
		return nil, errors.New("failed to load plugin : AddPlugin function has wrong signature")
	}

	return addPlugin(log), nil
}
