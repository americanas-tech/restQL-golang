package cache

import (
	"context"

	"github.com/b2wdigital/restQL-golang/v4/internal/platform/persistence"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

// MappingsReaderCache is a caching wrapper that
// implements the MappingsReader interface.
type MappingsReaderCache struct {
	log   restql.Logger
	cache *Cache
}

// NewMappingsReaderCache constructs a MappingsReaderCache instance.
func NewMappingsReaderCache(log restql.Logger, c *Cache) *MappingsReaderCache {
	return &MappingsReaderCache{log: log, cache: c}
}

// FromTenant returns a cached mapping index if present, fetching it otherwise.
func (c *MappingsReaderCache) FromTenant(ctx context.Context, tenant string) (map[string]restql.Mapping, error) {
	result, err := c.cache.Get(ctx, tenant)
	if err != nil {
		return nil, err
	}

	mappings, ok := result.(map[string]restql.Mapping)
	if !ok {
		log := restql.GetLogger(ctx)
		err := errors.Errorf("invalid mapping cache content type: %T", result)

		log.Error("failed to convert cache content", err)
		return nil, err
	}

	return mappings, nil
}

// TenantCacheLoader is the strategy to load
// values for the cached mappings reader.
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

// QueryReaderCache is a caching wrapper that
// implements the QueryReader interface.
type QueryReaderCache struct {
	log   restql.Logger
	cache *Cache
}

// NewQueryReaderCache constructs a QueryReaderCache instance.
func NewQueryReaderCache(log restql.Logger, c *Cache) *QueryReaderCache {
	return &QueryReaderCache{log: log, cache: c}
}

// Get returns a cached saved query if present, fetching it otherwise.
func (c *QueryReaderCache) Get(ctx context.Context, namespace, id string, revision int) (restql.SavedQuery, error) {
	cacheKey := cacheQueryKey{namespace: namespace, id: id, revision: revision}
	result, err := c.cache.Get(ctx, cacheKey)
	if err != nil {
		return restql.SavedQuery{}, err
	}

	query, ok := result.(restql.SavedQuery)
	if !ok {
		log := restql.GetLogger(ctx)
		err := errors.Errorf("invalid query cache content type: %T", result)

		log.Error("failed to convert cache content", err)
		return restql.SavedQuery{}, err
	}

	return query, nil
}

// QueryCacheLoader is the strategy to load
// values for the cached query reader.
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
