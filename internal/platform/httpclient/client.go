package httpclient

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/plugins"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
)

// New constructs an HTTPClient instances.
func New(log restql.Logger, pm plugins.Lifecycle, cfg *conf.Config) domain.HTTPClient {
	return newFastHttpClient(log, pm, cfg)
}
