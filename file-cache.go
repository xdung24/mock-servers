package main

import (
	"context"

	cache "github.com/eko/gocache/lib/v4/cache"
	gocachestore "github.com/eko/gocache/store/go_cache/v4"
	gocache "github.com/patrickmn/go-cache"
)

type CacheManager struct {
	manager *cache.Cache[[]byte]
	ctx     context.Context
}

func newCacheManager(ctx context.Context) *CacheManager {
	gocacheClient := gocache.New(0, 0) // No expiration, no cleanup
	gocacheStore := gocachestore.NewGoCache(gocacheClient)

	cacheManager := cache.New[[]byte](gocacheStore)
	return &CacheManager{
		manager: cacheManager,
		ctx:     ctx,
	}
}

func (c *CacheManager) read(id string) ([]byte, bool) {
	file_data, err := c.manager.Get(c.ctx, id)
	if err == nil {
		return file_data, true
	}
	return nil, false
}

func (c *CacheManager) update(id string, file_content []byte) {
	c.manager.Set(c.ctx, id, file_content)
}
