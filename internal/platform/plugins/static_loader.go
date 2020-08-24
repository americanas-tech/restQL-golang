package plugins

import (
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

func loadLifecyclePlugins(logger restql.Logger) []restql.LifecyclePlugin {
	var ps []restql.LifecyclePlugin
	for _, pluginInfo := range restql.GetLifecyclePlugins() {
		p, err := pluginInfo.New(logger)
		if err != nil {
			logger.Error("failed to load plugin", err)
			continue
		}

		pluginInstance, ok := p.(restql.LifecyclePlugin)
		if !ok {
			logger.Error("failed to load plugin", errors.Errorf("plugin of incorrect type: %T", p))
			continue
		}

		logger.Debug("plugin loaded", "name", pluginInstance.Name())
		ps = append(ps, pluginInstance)
	}
	return ps
}
