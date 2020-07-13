package database

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type tenant struct {
	Mappings map[string]string
}

type revision struct {
	Text       string
	Deprecated bool
}

type query struct {
	Name      string
	Namespace string
	Size      int
	Revisions []revision
}

type mongoDatabase struct {
	logger  *logger.Logger
	client  *mongo.Client
	options dbOptions
}

func (md mongoDatabase) FindMappingsForTenant(ctx context.Context, tenantId string) ([]domain.Mapping, error) {
	log := restql.GetLogger(ctx)
	mappingsTimeout := md.options.MappingsTimeout

	var cancel context.CancelFunc
	if mappingsTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, mappingsTimeout)
		defer cancel()
	}
	log.Debug("mappings timeout defined", "timeout", mappingsTimeout)

	var t tenant

	collection := md.client.Database(md.options.DatabaseName).Collection("tenant")
	err := collection.FindOne(ctx, bson.M{"_id": tenantId}).Decode(&t)
	if err != nil {
		return nil, err
	}

	i := 0
	result := make([]domain.Mapping, len(t.Mappings))
	for resourceName, url := range t.Mappings {
		mapping, err := domain.NewMapping(resourceName, url)
		if err != nil {
			continue
		}

		result[i] = mapping
		i++
	}

	return result, nil
}

func (md mongoDatabase) FindQuery(ctx context.Context, namespace string, name string, revision int) (domain.SavedQuery, error) {
	log := restql.GetLogger(ctx)

	queryTimeout := md.options.QueryTimeout
	if queryTimeout > 0 {
		ctx, _ = context.WithTimeout(ctx, queryTimeout)
	}
	log.Debug("query timeout defined", "timeout", queryTimeout)

	var q query

	collection := md.client.Database(md.options.DatabaseName).Collection("query")
	err := collection.FindOne(ctx, bson.M{"name": name, "namespace": namespace}).Decode(&q)
	if err != nil {
		return domain.SavedQuery{}, err
	}

	if q.Size < revision || revision < 0 {
		return domain.SavedQuery{}, errors.Errorf("invalid revision for query %s/%s : major revision %d : given revision %d", namespace, name, q.Size, revision)
	}

	r := q.Revisions[revision-1]

	return domain.SavedQuery{Text: r.Text, Deprecated: r.Deprecated}, nil
}
