package utils

import (
	"sync"
	"sync/atomic"
	"time"

	"audit-query-mcp-server/types"
)

// CacheEntry represents a cached audit result
type CacheEntry struct {
	Result    *types.AuditResult
	Timestamp time.Time
	TTL       time.Duration
}

// Cache provides a simple in-memory cache for audit results
type Cache struct {
	entries map[string]*CacheEntry
	mutex   sync.RWMutex
	ttl     time.Duration
	hits    int64
	misses  int64
}

// NewCache creates a new cache instance with default TTL
func NewCache(defaultTTL time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]*CacheEntry),
		ttl:     defaultTTL,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a cached result by query ID
func (c *Cache) Get(queryID string) (*types.AuditResult, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.entries[queryID]
	if !exists {
		atomic.AddInt64(&c.misses, 1)
		return nil, false
	}

	// Check if entry has expired
	if time.Since(entry.Timestamp) > entry.TTL {
		// Entry has expired, remove it
		c.mutex.RUnlock()
		c.mutex.Lock()
		delete(c.entries, queryID)
		c.mutex.Unlock()
		c.mutex.RLock()
		atomic.AddInt64(&c.misses, 1)
		return nil, false
	}

	atomic.AddInt64(&c.hits, 1)
	return entry.Result, true
}

// Set stores a result in the cache
func (c *Cache) Set(queryID string, result *types.AuditResult) {
	c.SetWithTTL(queryID, result, c.ttl)
}

// SetWithTTL stores a result in the cache with custom TTL
func (c *Cache) SetWithTTL(queryID string, result *types.AuditResult, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[queryID] = &CacheEntry{
		Result:    result,
		Timestamp: time.Now(),
		TTL:       ttl,
	}
}

// Delete removes a result from the cache
func (c *Cache) Delete(queryID string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.entries, queryID)
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries = make(map[string]*CacheEntry)
}

// ResetStats resets the cache hit/miss statistics
func (c *Cache) ResetStats() {
	atomic.StoreInt64(&c.hits, 0)
	atomic.StoreInt64(&c.misses, 0)
}

// Size returns the number of entries in the cache
func (c *Cache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.entries)
}

// cleanup periodically removes expired entries
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for queryID, entry := range c.entries {
			if now.Sub(entry.Timestamp) > entry.TTL {
				delete(c.entries, queryID)
			}
		}
		c.mutex.Unlock()
	}
}

// GetStats returns cache statistics
func (c *Cache) GetStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	stats := make(map[string]interface{})
	stats["size"] = len(c.entries)
	stats["default_ttl"] = c.ttl.String()
	stats["hits"] = atomic.LoadInt64(&c.hits)
	stats["misses"] = atomic.LoadInt64(&c.misses)

	// Calculate hit rate
	total := atomic.LoadInt64(&c.hits) + atomic.LoadInt64(&c.misses)
	if total > 0 {
		hitRate := float64(atomic.LoadInt64(&c.hits)) / float64(total) * 100
		stats["hit_rate"] = hitRate
	} else {
		stats["hit_rate"] = 0.0
	}

	// Count entries by age
	ageStats := make(map[string]int)
	now := time.Now()
	for _, entry := range c.entries {
		age := now.Sub(entry.Timestamp)
		switch {
		case age < time.Minute:
			ageStats["<1m"]++
		case age < time.Hour:
			ageStats["<1h"]++
		case age < 24*time.Hour:
			ageStats["<24h"]++
		default:
			ageStats[">24h"]++
		}
	}
	stats["age_distribution"] = ageStats

	return stats
}
