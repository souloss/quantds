package manager

import (
	"testing"
	"time"
)

func TestMemoryCacheStats(t *testing.T) {
	c := NewMemoryCache()
	defer c.Close()

	// Initially, no hits or misses
	stats := c.Stats()
	if stats.Hits != 0 || stats.Misses != 0 || stats.EntryCount != 0 {
		t.Errorf("initial stats: got hits=%d, misses=%d, entries=%d; want all 0", stats.Hits, stats.Misses, stats.EntryCount)
	}

	// Set and get → hit
	c.Set("key1", []byte("val1"), time.Minute)
	data, ok := c.Get("key1")
	if !ok || string(data) != "val1" {
		t.Fatalf("Get(key1) = %q, %v; want val1, true", data, ok)
	}
	stats = c.Stats()
	if stats.Hits != 1 {
		t.Errorf("after 1 hit: hits = %d, want 1", stats.Hits)
	}
	if stats.EntryCount != 1 {
		t.Errorf("after 1 set: entries = %d, want 1", stats.EntryCount)
	}

	// Get non-existent key → miss
	_, ok = c.Get("nonexistent")
	if ok {
		t.Fatal("Get(nonexistent) should return false")
	}
	stats = c.Stats()
	if stats.Misses != 1 {
		t.Errorf("after 1 miss: misses = %d, want 1", stats.Misses)
	}

	// Get expired key → miss
	c.Set("expired", []byte("val"), 1*time.Nanosecond)
	time.Sleep(time.Millisecond) // wait for expiry
	_, ok = c.Get("expired")
	if ok {
		t.Fatal("Get(expired) should return false")
	}
	stats = c.Stats()
	if stats.Misses != 2 {
		t.Errorf("after expired miss: misses = %d, want 2", stats.Misses)
	}

	// Hit rate calculation
	// hits=1, misses=2, total=3 → hitRate ≈ 0.333
	stats = c.Stats()
	if stats.HitRate < 0.33 || stats.HitRate > 0.34 {
		t.Errorf("hit rate = %f, want ~0.333", stats.HitRate)
	}
}

func TestTwoLevelCacheStats(t *testing.T) {
	c := NewTwoLevelCache(time.Minute, time.Minute)
	defer c.Close()

	// Populate request cache
	c.SetRequest("req1", []byte("reqval1"))
	// Populate fetch cache
	c.SetFetch("fetch1", []byte("fetchval1"))

	// Hit on request cache
	data, ok := c.GetRequest("req1")
	if !ok || string(data) != "reqval1" {
		t.Fatalf("GetRequest(req1) = %q, %v; want reqval1, true", data, ok)
	}

	// Miss on request cache
	_, ok = c.GetRequest("nonexistent")
	if ok {
		t.Fatal("GetRequest(nonexistent) should return false")
	}

	// Hit on fetch cache
	data, ok = c.GetFetch("fetch1")
	if !ok || string(data) != "fetchval1" {
		t.Fatalf("GetFetch(fetch1) = %q, %v; want fetchval1, true", data, ok)
	}

	stats := c.Stats()
	if stats.RequestCache.Hits != 1 {
		t.Errorf("request cache hits = %d, want 1", stats.RequestCache.Hits)
	}
	if stats.RequestCache.Misses != 1 {
		t.Errorf("request cache misses = %d, want 1", stats.RequestCache.Misses)
	}
	if stats.FetchCache.Hits != 1 {
		t.Errorf("fetch cache hits = %d, want 1", stats.FetchCache.Hits)
	}
	if stats.FetchCache.Misses != 0 {
		t.Errorf("fetch cache misses = %d, want 0", stats.FetchCache.Misses)
	}
}

func TestManagerCacheStats(t *testing.T) {
	// Manager without cache → nil
	m1 := NewManager[klineReq, klineResp]()
	if cs := m1.CacheStats(); cs != nil {
		t.Errorf("CacheStats() without cache = %v, want nil", cs)
	}
	m1.Close()

	// Manager with cache → stats available
	m2 := NewManager[klineReq, klineResp](
		WithTwoLevelCache[klineReq, klineResp](time.Minute, time.Minute),
	)
	cs := m2.CacheStats()
	if cs == nil {
		t.Fatal("CacheStats() with cache = nil, want non-nil")
	}
	if cs.RequestCache.Hits != 0 || cs.FetchCache.Hits != 0 {
		t.Errorf("initial CacheStats: request hits=%d, fetch hits=%d; want both 0", cs.RequestCache.Hits, cs.FetchCache.Hits)
	}
	m2.Close()
}

// Minimal types for Manager generic instantiation in tests
type klineReq struct{}
type klineResp struct{}
