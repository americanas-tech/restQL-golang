package httpclient

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/plugins"
)

func New(log *logger.Logger, pm plugins.Lifecycle, cfg *conf.Config) domain.HttpClient {
	return newNativeHttpClient(log, pm, cfg)
}
