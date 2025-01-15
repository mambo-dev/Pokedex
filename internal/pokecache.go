package internal

import (
	"sync"
	"time"
)

type Cache struct {
	mu         sync.Mutex
	cacheEntry map[string]CacheEntry
}

type CacheEntry struct {
	value     []byte
	createdAt time.Time
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheEntry[key] = CacheEntry{
		value:     val,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cacheEntry, ok := c.cacheEntry[key]

	if ok {
		return cacheEntry.value, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for _, cache := range c.cacheEntry {

		for range ticker.C {
			if time.Since(cache.createdAt) > interval {
				delete(c.cacheEntry, "")
			}
		}
	}

}

func NewCache(val []byte, key string, interval time.Duration) *Cache {
	cache := &Cache{}

	cache.reapLoop(interval)
	return cache

}
