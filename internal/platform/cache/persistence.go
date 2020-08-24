package cache

import (
	"context"

	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/persistence"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

type MappingsReaderCache struct {
	log   *logger.Logger
	cache *Cache
}

func NewMappingsReaderCache(log *logger.Logger, c *Cache) *MappingsReaderCache {
	return &MappingsReaderCache{log: log, cache: c}
}

func (c *MappingsReaderCache) FromTenant(ctx context.Context, tenant string) (map[string]restql.Mapping, error) {
	result, err := c.cache.Get(ctx, tenant)
	if err != nil {
		return nil, err
	}

	mappings, ok := result.(map[string]restql.Mapping)
	if !ok {
		log := restql.GetLogger(ctx)
		log.Info("failed to convert cache content", "content", result)
	}

	return mappings, nil
}

func TenantCacheLoader(mr persistence.MappingsReader) Loader {
	return func(ctx context.Context, key interface{}) (interface{}, error) {
		tenant, ok := key.(string)
		if !ok {
			return nil, errors.Errorf("invalid key type : got %T", key)
		}

		mappings, err := mr.FromTenant(ctx, tenant)
		if err != nil {
			return nil, err
		}

		return mappings, nil
	}
}

type cacheQueryKey struct {
	namespace string
	id        string
	revision  int
}

type QueryReaderCache struct {
	log   *logger.Logger
	cache *Cache
}

func NewQueryReaderCache(log *logger.Logger, c *Cache) *QueryReaderCache {
	return &QueryReaderCache{log: log, cache: c}
}

func (c *QueryReaderCache) Get(ctx context.Context, namespace, id string, revision int) (restql.SavedQuery, error) {
	cacheKey := cacheQueryKey{namespace: namespace, id: id, revision: revision}
	result, err := c.cache.Get(ctx, cacheKey)
	if err != nil {
		return restql.SavedQuery{}, err
	}

	query, ok := result.(restql.SavedQuery)
	if !ok {
		log := restql.GetLogger(ctx)
		log.Info("failed to convert cache content", "content", result)
	}

	return query, nil
}

func QueryCacheLoader(qr persistence.QueryReader) Loader {
	return func(ctx context.Context, key interface{}) (interface{}, error) {
		cacheKey, ok := key.(cacheQueryKey)
		if !ok {
			return nil, errors.Errorf("invalid key type : got %T", key)
		}

		query, err := qr.Get(ctx, cacheKey.namespace, cacheKey.id, cacheKey.revision)
		if err != nil {
			return nil, err
		}

		return query, nil
	}
}
