package persistence

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

type Database interface {
	FindMappingsForTenant(ctx context.Context, tenantId string) ([]restql.Mapping, error)
	FindQuery(ctx context.Context, namespace string, name string, revision int) (restql.SavedQuery, error)
}

func NewDatabase(log *logger.Logger) (Database, error) {
	pluginInfo, found := restql.GetDatabasePlugin()
	if !found {
		log.Info("no database plugin provided")
		return noOpDatabase{}, nil
	}

	dbPlugin, err := pluginInfo.New(log)
	if err != nil {
		return noOpDatabase{}, err
	}

	if dbPlugin == nil {
		log.Info("empty database instance returned by plugin", "plugin", pluginInfo.Name)
		return noOpDatabase{}, nil
	}

	database, ok := dbPlugin.(restql.DatabasePlugin)
	if !ok {
		return noOpDatabase{}, errors.Errorf("failed to cast database plugin, unknown type: %T", dbPlugin)
	}

	return database, nil
}

var ErrNoDatabase = errors.New("no op database")

type noOpDatabase struct{}

func (n noOpDatabase) FindMappingsForTenant(ctx context.Context, tenantId string) ([]restql.Mapping, error) {
	return []restql.Mapping{}, ErrNoDatabase
}

func (n noOpDatabase) FindQuery(ctx context.Context, namespace string, name string, revision int) (restql.SavedQuery, error) {
	return restql.SavedQuery{}, ErrNoDatabase
}
