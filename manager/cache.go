package manager

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"
)

type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, ttl time.Duration)
	Delete(key string)
	Clear()
}

type MemoryCache struct {
	mu   sync.RWMutex
	data map[string]*cacheEntry
}

type cacheEntry struct {
	data      []byte
	expiresAt time.Time
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]*cacheEntry),
	}
}

func (c *MemoryCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.data[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.expiresAt) {
		return nil, false
	}
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

func BuildCacheKey(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	h := sha256.New()
	h.Write(jsonData)
	return hex.EncodeToString(h.Sum(nil))
}
