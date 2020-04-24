package plugins

import (
	"github.com/b2wdigital/restQL-golang/pkg/restql"
)

func loadStaticPlugin(logger restql.Logger) []restql.Plugin {
	var ps []restql.Plugin
	for _, loader := range restql.GetPluginLoaders() {
		p, err := loader(logger)
		if err != nil {
			logger.Error("failed to load plugin", err)
			continue
		}

		logger.Debug("plugin loaded", "name", p.Name())
		ps = append(ps, p)
	}
	return ps
}
