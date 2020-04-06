package database

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Database interface {
	FindMappingsForTenant(ctx context.Context, tenantId string) ([]domain.Mapping, error)
	FindQuery(ctx context.Context, namespace string, name string, revision int) (string, error)
}

func New(log *logger.Logger, connectionString string) (Database, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	return mongoDatabase{logger: log, client: client}, nil
}
