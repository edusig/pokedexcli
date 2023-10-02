package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.RWMutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c *Cache) readLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		tick, ok := <-ticker.C
		if !ok {
			break
		}
		for key, val := range c.entries {
			if tick.Sub(val.createdAt) > interval {
				c.mu.Lock()
				delete(c.entries, key)
				c.mu.Unlock()
			}
		}
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		entries: make(map[string]cacheEntry),
		mu:      sync.RWMutex{},
	}
	go cache.readLoop(interval)
	return &cache
}
