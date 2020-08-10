package cache

import (
	"context"
	"time"

	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/bluele/gcache"
	"github.com/pkg/errors"
)

type cacheItem struct {
	key        interface{}
	value      interface{}
	expiration time.Time
}

func (i cacheItem) Expired() bool {
	return !i.expiration.IsZero() && time.Now().After(i.expiration)
}

type Loader func(ctx context.Context, key interface{}) (interface{}, error)

type Option func(c *Cache)

func WithRefreshInterval(interval time.Duration) Option {
	return func(c *Cache) {
		c.refreshInterval = interval
	}
}

func WithRefreshQueueLength(length int) Option {
	return func(c *Cache) {
		c.refreshQueueLength = length
	}
}

func WithExpiration(expiration time.Duration) Option {
	return func(c *Cache) {
		c.expiration = expiration
	}
}

type Cache struct {
	log                *logger.Logger
	gcache             gcache.Cache
	loader             Loader
	refreshWorkCh      chan interface{}
	expiration         time.Duration
	refreshInterval    time.Duration
	refreshQueueLength int
}

func New(log *logger.Logger, size int, loader Loader, options ...Option) *Cache {
	c := gcache.New(size).LRU().Build()

	cache := Cache{
		log:    log,
		gcache: c,
		loader: loader,
	}

	for _, option := range options {
		option(&cache)
	}

	if cache.refreshInterval > 0 && cache.refreshQueueLength > 0 {
		rw := cache.setupRefreshWorker()
		go rw.Run()
	}

	return &cache
}

func (c *Cache) Get(ctx context.Context, key interface{}) (interface{}, error) {
	obj, err := c.gcache.Get(key)

	switch {
	case err == gcache.KeyNotFoundError:
		item, err := c.loadItem(ctx, key)
		if err != nil {
			return nil, err
		}

		return item.value, nil
	case err != nil:
		c.log.Error("failed to retrieve value from cache", err, "key", key)
		return nil, err
	}

	item, ok := obj.(cacheItem)
	if !ok {
		err := errors.Errorf("invalid cache item : %v", item)
		c.log.Error("failed to cast cache item", err, "key", key)
		return nil, err
	}

	if item.Expired() {
		go func() {
			c.refreshWorkCh <- item.key
		}()
	}

	return item.value, nil
}

func (c *Cache) loadItem(ctx context.Context, key interface{}) (cacheItem, error) {
	value, err := c.loader(ctx, key)
	if err != nil {
		c.log.Debug("failed to load value to populate cache", "error", err)
		return cacheItem{}, err
	}

	item := cacheItem{
		key:   key,
		value: value,
	}
	if c.expiration > 0 {
		item.expiration = time.Now().Add(c.expiration)
	}

	err = c.gcache.Set(key, item)
	if err != nil {
		c.log.Error("failed to set value on cache", err)
		return cacheItem{}, err
	}
	return item, nil
}

func (c *Cache) setupRefreshWorker() *refreshWorker {
	ticker := time.NewTicker(c.refreshInterval)
	refreshWorkCh := make(chan interface{}, c.refreshQueueLength)

	c.refreshWorkCh = refreshWorkCh

	rw := refreshWorker{
		log:           c.log,
		cache:         c,
		refreshWorkCh: refreshWorkCh,
		ticker:        ticker,
	}

	return &rw
}

type refreshWorker struct {
	log           *logger.Logger
	cache         *Cache
	refreshFn     Loader
	refreshWorkCh chan interface{}
	ticker        *time.Ticker
}

func (rw *refreshWorker) Run() {
	for {
		select {
		case <-rw.ticker.C:
			for key := range rw.refreshWorkCh {
				key := key
				go func() {
					_, err := rw.cache.loadItem(context.Background(), key)
					if err != nil {
						rw.log.Error("failed to refresh cache item in background", err)
					}
				}()
			}
		}
	}
}
