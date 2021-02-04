package cache

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/parser"
	"github.com/pkg/errors"
)

// ParserCache is a caching wrapper that implements the Parser interface.
type ParserCache struct {
	log   restql.Logger
	cache *Cache
}

// NewParserCache constructs a ParserCache instance.
func NewParserCache(log restql.Logger, c *Cache) ParserCache {
	return ParserCache{log: log, cache: c}
}

// Parse returns a cached QueryRevisions internal representation if
// present, transforming the query text into one otherwise.
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

// ParserCacheLoader is the strategy to load
// values for the cached parser.
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
