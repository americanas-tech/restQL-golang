package httpclient

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/internal/platform/plugins"
)

func New(log *logger.Logger, pm plugins.Manager, cfg *conf.Config) domain.HttpClient {
	return newNativeHttpClient(log, pm, cfg)
}
