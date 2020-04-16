package plugins

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
	"github.com/pkg/errors"
	"os"
	"path"
	"plugin"
	"time"
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
			timeout, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)

			p, err := loadPlugin(timeout, log, pluginPath)
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

func loadPlugin(ctx context.Context, log *logger.Logger, pluginPath string) (restql.Plugin, error) {
	out := make(chan restql.Plugin)
	var pluginerr error

	go func() {
		p, err := plugin.Open(pluginPath)
		if err != nil {
			pluginerr = err
			return
		}

		addPluginSym, err := p.Lookup("AddPlugin")
		if err != nil {
			pluginerr = err
			return
		}

		addPlugin, ok := addPluginSym.(func(log restql.Logger) (restql.Plugin, error))
		if !ok {
			pluginerr = errors.New("failed to load plugin : AddPlugin function has wrong signature")
			return
		}

		result, err := addPlugin(log)
		if err != nil {
			pluginerr = err
			return
		}

		out <- result
	}()

	select {
	case p := <-out:
		if pluginerr != nil {
			return nil, pluginerr
		}

		return p, nil
	case <-ctx.Done():
		return nil, errors.New("failed to load plugin : timed out")
	}
}
