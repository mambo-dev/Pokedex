package internal

import (
	"fmt"
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
	fmt.Printf("=== cache %v added \n", key)
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheEntry[key] = CacheEntry{
		value:     val,
		createdAt: time.Now(),
	}

}

func (c *Cache) Get(key string) ([]byte, bool) {
	fmt.Printf("=== getting %v from cache \n", key)
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

	ticker := time.NewTicker(interval)

	for range ticker.C {
		c.mu.Lock()
		for key, value := range c.cacheEntry {
			if time.Since(value.createdAt) > interval {

				delete(c.cacheEntry, key)
			}
		}
		c.mu.Unlock()
	}

}

func NewCache(interval time.Duration) *Cache {
	fmt.Println("=== new cache created")
	cache := &Cache{
		cacheEntry: make(map[string]CacheEntry),
	}

	go cache.reapLoop(interval)
	return cache

}
