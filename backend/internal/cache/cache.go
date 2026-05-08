package cache

import (
	"sync"
	"time"
)

// CacheEntry speichert Daten mit Ablaufzeit
type CacheEntry struct {
	Data      []byte
	ExpiresAt time.Time
}

// Cache ist ein einfacher In-Memory Cache
type Cache struct {
	mu      sync.RWMutex
	entries map[string]CacheEntry
	ttl     time.Duration
}

// New erstellt einen neuen Cache mit TTL
func New(ttl time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]CacheEntry),
		ttl:     ttl,
	}
	// Abgelaufene Einträge alle 30 Sekunden bereinigen
	go c.cleanup()
	return c
}

// Get gibt einen Cache-Eintrag zurück
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, found := c.entries[key]
	if !found || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Data, true
}

// Set speichert einen Eintrag im Cache
func (c *Cache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Invalidate löscht einen Eintrag aus dem Cache
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// cleanup bereinigt abgelaufene Einträge regelmäßig
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