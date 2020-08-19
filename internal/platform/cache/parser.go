package cache

import (
	"context"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/parser"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/pkg/errors"
)

type ParserCache struct {
	log   *logger.Logger
	cache *Cache
}

func NewParserCache(log *logger.Logger, c *Cache) ParserCache {
	return ParserCache{log: log, cache: c}
}

func (p ParserCache) Parse(queryStr string) (domain.Query, error) {
	result, err := p.cache.Get(context.Background(), queryStr)
	if err != nil {
		return domain.Query{}, err
	}

	query, ok := result.(domain.Query)
	if !ok {
		p.log.Info("failed to convert cache content", "content", result)
	}

	return query, nil
}

func ParserCacheLoader(p parser.Parser) Loader {
	return func(ctx context.Context, key interface{}) (interface{}, error) {
		queryStr, ok := key.(string)
		if !ok {
			return nil, errors.Errorf("invalid key type : got %T", key)
		}

		query, err := p.Parse(queryStr)
		if err != nil {
			return nil, err
		}

		return query, nil
	}
}
