package cache

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/internal/platform/persistence"
	"github.com/pkg/errors"
)

type MappingsReaderCache struct {
	log   *logger.Logger
	cache *Cache
	mr    persistence.MappingsReader
}

func NewMappingsReaderCache(log *logger.Logger, mr persistence.MappingsReader, c *Cache) *MappingsReaderCache {
	return &MappingsReaderCache{log: log, mr: mr, cache: c}
}

func (c *MappingsReaderCache) FromTenant(ctx context.Context, tenant string) (map[string]domain.Mapping, error) {
	result, err := c.cache.Get(ctx, tenant)
	if err != nil {
		return nil, err
	}

	mappings, ok := result.(map[string]domain.Mapping)
	if !ok {
		c.log.Info("failed to convert cache content", "content", result)
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
	qr    persistence.QueryReader
}

func NewQueryReaderCache(log *logger.Logger, qr persistence.QueryReader, c *Cache) *QueryReaderCache {
	return &QueryReaderCache{log: log, qr: qr, cache: c}
}

func (c *QueryReaderCache) Get(ctx context.Context, namespace, id string, revision int) (domain.SavedQuery, error) {
	cacheKey := cacheQueryKey{namespace: namespace, id: id, revision: revision}
	result, err := c.cache.Get(ctx, cacheKey)
	if err != nil {
		return domain.SavedQuery{}, err
	}

	query, ok := result.(domain.SavedQuery)
	if !ok {
		c.log.Info("failed to convert cache content", "content", result)
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
