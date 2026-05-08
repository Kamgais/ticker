package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data      []byte
	ExpiresAt time.Time
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
	ttl     time.Duration
}

func New(ttl time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]CacheEntry),
		ttl:     ttl,
	}
	go c.cleanup()
	return c
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, found := c.entries[key]
	if !found || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Data, true
}

func (c *Cache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.entries {
			if time.Now().After(entry.ExpiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}