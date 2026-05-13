package manager

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"
)

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, ttl time.Duration)
	Delete(key string)
	Clear()
}

type MemoryCache struct {
	mu       sync.RWMutex
	data     map[string]*cacheEntry
	stopCh   chan struct{}
	cleanupT time.Duration
	hits     atomic.Uint64
	misses   atomic.Uint64
}

type cacheEntry struct {
	data      []byte
	expiresAt time.Time
}

func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		data:     make(map[string]*cacheEntry),
		stopCh:   make(chan struct{}),
		cleanupT: time.Minute,
	}
	go c.cleanupLoop()
	return c
}

func (c *MemoryCache) cleanupLoop() {
	ticker := time.NewTicker(c.cleanupT)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCh:
			return
		}
	}
}

func (c *MemoryCache) Close() {
	close(c.stopCh)
}

func (c *MemoryCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.data[key]
	if !ok {
		c.misses.Add(1)
		return nil, false
	}
	if time.Now().After(entry.expiresAt) {
		c.misses.Add(1)
		return nil, false
	}
	c.hits.Add(1)
	return entry.data, true
}

func (c *MemoryCache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = &cacheEntry{
		data:      value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*cacheEntry)
}

func (c *MemoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, v := range c.data {
		if now.After(v.expiresAt) {
			delete(c.data, k)
		}
	}
}

var _ Cache = (*MemoryCache)(nil)

// CacheStats holds cache hit/miss statistics.
type CacheStats struct {
	Hits       uint64  // 缓存命中次数
	Misses     uint64  // 缓存未命中次数
	HitRate    float64 // 命中率 (0.0 - 1.0)
	EntryCount int     // 当前缓存条目数
}

// Stats returns the current cache statistics.
func (c *MemoryCache) Stats() CacheStats {
	c.mu.RLock()
	count := len(c.data)
	c.mu.RUnlock()

	hits := c.hits.Load()
	misses := c.misses.Load()
	total := hits + misses
	var hitRate float64
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}

	return CacheStats{
		Hits:       hits,
		Misses:     misses,
		HitRate:    hitRate,
		EntryCount: count,
	}
}

type TwoLevelCache struct {
	requestCache Cache
	fetchCache   Cache
	requestTTL   time.Duration
	fetchTTL     time.Duration
}

func NewTwoLevelCache(requestTTL, fetchTTL time.Duration) *TwoLevelCache {
	return &TwoLevelCache{
		requestCache: NewMemoryCache(),
		fetchCache:   NewMemoryCache(),
		requestTTL:   requestTTL,
		fetchTTL:     fetchTTL,
	}
}

// TwoLevelCacheStats holds aggregated stats for both cache levels.
type TwoLevelCacheStats struct {
	RequestCache CacheStats // 请求级缓存统计
	FetchCache   CacheStats // 数据级缓存统计
}

// Stats returns the aggregated cache statistics for both levels.
func (c *TwoLevelCache) Stats() TwoLevelCacheStats {
	var stats TwoLevelCacheStats
	if mc, ok := c.requestCache.(*MemoryCache); ok {
		stats.RequestCache = mc.Stats()
	}
	if mc, ok := c.fetchCache.(*MemoryCache); ok {
		stats.FetchCache = mc.Stats()
	}
	return stats
}

func (c *TwoLevelCache) GetRequest(key string) ([]byte, bool) {
	return c.requestCache.Get(key)
}

func (c *TwoLevelCache) SetRequest(key string, value []byte) {
	c.requestCache.Set(key, value, c.requestTTL)
}

func (c *TwoLevelCache) GetFetch(key string) ([]byte, bool) {
	return c.fetchCache.Get(key)
}

func (c *TwoLevelCache) SetFetch(key string, value []byte) {
	c.fetchCache.Set(key, value, c.fetchTTL)
}

func (c *TwoLevelCache) Clear() {
	c.requestCache.Clear()
	c.fetchCache.Clear()
}

func (c *TwoLevelCache) Close() {
	if mc, ok := c.requestCache.(*MemoryCache); ok {
		mc.Close()
	}
	if mc, ok := c.fetchCache.(*MemoryCache); ok {
		mc.Close()
	}
}

func BuildCacheKey(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	h := sha256.New()
	h.Write(jsonData)
	return hex.EncodeToString(h.Sum(nil))
}
