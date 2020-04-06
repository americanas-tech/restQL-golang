package database

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const databaseName = "restql"

type tenant struct {
	Mappings map[string]string
}

type mongoDatabase struct {
	logger *logger.Logger
	client *mongo.Client
}

func (md mongoDatabase) FindMappingsForTenant(ctx context.Context, tenantId string) ([]domain.Mapping, error) {
	var tenant tenant

	collection := md.client.Database(databaseName).Collection("tenant")
	err := collection.FindOne(ctx, bson.M{"_id": tenantId}).Decode(&tenant)
	if err != nil {
		return nil, err
	}

	i := 0
	result := make([]domain.Mapping, len(tenant.Mappings))
	for resourceName, url := range tenant.Mappings {
		mapping, err := domain.NewMapping(resourceName, url)
		if err != nil {
			continue
		}

		result[i] = mapping
		i++
	}

	return result, nil
}
