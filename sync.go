package zzcache

import "sync"

type syncCache struct {
	sync.RWMutex
	lru lruCache
}

func (c *syncCache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	return c.lru.Get(key)
}


func (c *syncCache) Set(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()

	return c.lru.Set(key, value)
}

func (c *syncCache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	return c.lru.Delete(key)
}

func (c *syncCache) Len() uint64 {
	return c.lru.Len()
}

func NewSyncCache(cacheSize uint64) *syncCache {
	cache := new(syncCache)
	cache.lru = *NewLRU(cacheSize)

	return cache
}
