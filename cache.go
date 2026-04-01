package erlc

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data      []byte
	ExpiresAt time.Time
}

func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

type Cache struct {
	data        map[string]CacheEntry
	maxSize     int
	defaultTTL  time.Duration
	currentSize int
	mu          sync.RWMutex
	hits        int
	misses      int
	evictions   int
	lastCleanup time.Time
	cleanupTTL  time.Duration
}

func NewCache(maxSize int, defaultTTL time.Duration) *Cache {
	if maxSize < 1 {
		maxSize = DefaultMaxCacheSize
	}
	if defaultTTL < 1 {
		defaultTTL = DefaultCacheTTL
	}

	c := &Cache{
		data:        make(map[string]CacheEntry),
		maxSize:     maxSize,
		defaultTTL:  defaultTTL,
		cleanupTTL:  time.Minute,
		lastCleanup: time.Now(),
	}

	go c.backgroundCleanup()
	return c
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.data[key]
	if !ok {
		c.misses++
		return nil, false
	}

	if entry.IsExpired() {
		c.misses++
		return nil, false
	}

	c.hits++
	return entry.Data, true
}

func (c *Cache) Set(key string, data []byte) {
	c.SetWithTTL(key, data, c.defaultTTL)
}

func (c *Cache) SetWithTTL(key string, data []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttl < 1 {
		ttl = c.defaultTTL
	}

	if existing, ok := c.data[key]; ok {
		c.currentSize -= len(existing.Data)
	}

	if c.currentSize+len(data) > c.maxSize {
		c.evictOldest()
	}

	c.data[key] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.currentSize += len(data)
}

func (c *Cache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.data[key]; ok {
		c.currentSize -= len(entry.Data)
		delete(c.data, key)
		return true
	}
	return false
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]CacheEntry)
	c.currentSize = 0
	c.hits = 0
	c.misses = 0
	c.evictions = 0
}

func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.data[key]
	if !ok {
		return false
	}

	return !entry.IsExpired()
}

func (c *Cache) Stats() *CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(c.hits) / float64(total)
	}

	return &CacheStats{
		Hits:      c.hits,
		Misses:    c.misses,
		Evictions: c.evictions,
		HitRate:   hitRate,
		Size:      c.currentSize,
		MaxSize:   c.maxSize,
		Entries:   len(c.data),
	}
}

type CacheStats struct {
	Hits      int
	Misses    int
	Evictions int
	HitRate   float64
	Size      int
	MaxSize   int
	Entries   int
}

func (c *Cache) evictOldest() {
	for key, entry := range c.data {
		if entry.IsExpired() {
			c.currentSize -= len(entry.Data)
			delete(c.data, key)
			c.evictions++
			return
		}
	}

	if len(c.data) > 0 {
		var oldestKey string
		var oldestTime time.Time
		first := true

		for key, entry := range c.data {
			if first {
				oldestKey = key
				oldestTime = entry.ExpiresAt
				first = false
			} else if entry.ExpiresAt.Before(oldestTime) {
				oldestKey = key
				oldestTime = entry.ExpiresAt
			}
		}

		if entry, ok := c.data[oldestKey]; ok {
			c.currentSize -= len(entry.Data)
			delete(c.data, oldestKey)
			c.evictions++
		}
	}
}

func (c *Cache) backgroundCleanup() {
	ticker := time.NewTicker(c.cleanupTTL)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.After(entry.ExpiresAt) {
			c.currentSize -= len(entry.Data)
			delete(c.data, key)
		}
	}
	c.lastCleanup = now
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentSize
}

func (c *Cache) Entries() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.data)
}
