package persistence

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"github.com/b2wdigital/restQL-golang/internal/platform/cache"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/pkg/errors"
)

type CacheMappingsReader struct {
	log   *logger.Logger
	mr    eval.MappingsReader
	cache *cache.Cache
}

func NewCacheMappingsReader(log *logger.Logger, mr eval.MappingsReader, c *cache.Cache) eval.MappingsReader {
	return &CacheMappingsReader{log: log, mr: mr, cache: c}
}

func (c *CacheMappingsReader) FromTenant(ctx context.Context, tenant string) (map[string]domain.Mapping, error) {
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

func TenantCacheLoader(mr eval.MappingsReader) cache.Loader {
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

type CacheQueryReader struct {
	log   *logger.Logger
	qr    eval.QueryReader
	cache *cache.Cache
}

func NewCacheQueryReader(log *logger.Logger, qr eval.QueryReader, c *cache.Cache) eval.QueryReader {
	return &CacheQueryReader{log: log, qr: qr, cache: c}
}

func (c *CacheQueryReader) Get(ctx context.Context, namespace, id string, revision int) (string, error) {
	cacheKey := cacheQueryKey{namespace: namespace, id: id, revision: revision}
	result, err := c.cache.Get(ctx, cacheKey)
	if err != nil {
		return "", err
	}

	query, ok := result.(string)
	if !ok {
		c.log.Info("failed to convert cache content", "content", result)
	}

	return query, nil
}

func QueryCacheLoader(qr eval.QueryReader) cache.Loader {
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
