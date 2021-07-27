package zzcache

import (
	"sync"
)

type GroupCache struct {
	name string
	syncCache
}

var (
	locker sync.RWMutex
	groups = make(map[string]*GroupCache)
)

func NewGroup(name string, cacheSize uint64) *GroupCache {
	cache := new(GroupCache)
	cache.name = name
	cache.syncCache = *NewSyncCache(cacheSize)

	locker.Lock()
	defer locker.Unlock()

	groups[name] = cache

	return cache
}

func GetGroup(name string) *GroupCache {
	locker.RLock()
	defer locker.RUnlock()

	if group, ok := groups[name]; ok {
		return group
	} else {
		return nil
	}
}
