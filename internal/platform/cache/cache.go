package cache

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/bluele/gcache"
	"github.com/pkg/errors"
	"time"
)

type cacheItem struct {
	key        interface{}
	value      interface{}
	expiration time.Time
}

func (i cacheItem) Expired() bool {
	return time.Now().After(i.expiration)
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

type Cache struct {
	log                *logger.Logger
	gcache             gcache.Cache
	loader             Loader
	refreshWorkCh      chan interface{}
	expiration         time.Duration
	refreshInterval    time.Duration
	refreshQueueLength int
}

func New(log *logger.Logger, loader Loader, size int, options ...Option) *Cache {
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
			c.refreshWorkCh <- item
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
		key:        key,
		value:      value,
		expiration: time.Now().Add(c.expiration),
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
		cache:         c,
		refreshWorkCh: refreshWorkCh,
		ticker:        ticker,
	}

	return &rw
}

type refreshWorker struct {
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
				go func() {
					//todo: tratar context para atualização em background
					rw.cache.loadItem(context.Background(), key)
				}()
			}
		}
	}
}
