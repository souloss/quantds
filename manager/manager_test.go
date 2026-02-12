package manager

import (
	"context"
	"testing"
	"time"

	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/request"
)

type testReq struct {
	Symbol string
}

type testResp struct {
	Data string
}

type testProvider struct {
	name string
	data string
	err  error
}

func (p *testProvider) Name() string {
	return p.name
}

func (p *testProvider) SupportedMarkets() []domain.Market {
	return []domain.Market{domain.MarketCN}
}

func (p *testProvider) CanHandle(symbol string) bool {
	return true
}

func (p *testProvider) Fetch(ctx context.Context, client request.Client, req testReq) (testResp, *RequestTrace, error) {
	return testResp{Data: p.data}, NewRequestTrace(p.name), p.err
}

func TestNewManager(t *testing.T) {
	m := NewManager[testReq, testResp]()
	if m == nil {
		t.Fatal("NewManager() returned nil")
	}
	defer m.Close()
}

func TestManager_Register(t *testing.T) {
	m := NewManager[testReq, testResp]()
	defer m.Close()

	p := &testProvider{name: "test", data: "hello"}
	m.Register(p, WithPriority(10))

	providers := m.Providers()
	if len(providers) != 1 {
		t.Errorf("Providers() = %v, want 1", len(providers))
	}
	if providers[0] != "test" {
		t.Errorf("Providers()[0] = %v, want test", providers[0])
	}
}

func TestManager_Fetch(t *testing.T) {
	m := NewManager[testReq, testResp](
		WithProvider[testReq, testResp](&testProvider{name: "p1", data: "data1"}, WithPriority(10)),
		WithProvider[testReq, testResp](&testProvider{name: "p2", data: "data2"}, WithPriority(5)),
	)
	defer m.Close()

	result, err := m.Fetch(context.Background(), testReq{Symbol: "test"})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}

	if result.Provider != "p1" {
		t.Errorf("Provider = %v, want p1", result.Provider)
	}
	if result.Data.Data != "data1" {
		t.Errorf("Data = %v, want data1", result.Data.Data)
	}
}

func TestManager_Fetch_Fallback(t *testing.T) {
	m := NewManager[testReq, testResp](
		WithProvider[testReq, testResp](&testProvider{name: "p1", data: "", err: ErrAllProviderFailed}, WithPriority(10)),
		WithProvider[testReq, testResp](&testProvider{name: "p2", data: "data2"}, WithPriority(5)),
	)
	defer m.Close()

	result, err := m.Fetch(context.Background(), testReq{Symbol: "test"})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}

	if result.Provider != "p2" {
		t.Errorf("Provider = %v, want p2 (fallback)", result.Provider)
	}
}

func TestManager_Fetch_NoProvider(t *testing.T) {
	m := NewManager[testReq, testResp]()
	defer m.Close()

	_, err := m.Fetch(context.Background(), testReq{Symbol: "test"})
	if err != ErrNoProvider {
		t.Errorf("Fetch() error = %v, want %v", err, ErrNoProvider)
	}
}

func TestManager_FetchFrom(t *testing.T) {
	m := NewManager[testReq, testResp](
		WithProvider[testReq, testResp](&testProvider{name: "p1", data: "data1"}, WithPriority(10)),
	)
	defer m.Close()

	result, err := m.FetchFrom(context.Background(), "p1", testReq{Symbol: "test"})
	if err != nil {
		t.Fatalf("FetchFrom() error = %v", err)
	}

	if result.Provider != "p1" {
		t.Errorf("Provider = %v, want p1", result.Provider)
	}
}

func TestManager_Cache(t *testing.T) {
	m := NewManager[testReq, testResp](
		WithTwoLevelCache[testReq, testResp](time.Minute, time.Minute),
		WithProvider[testReq, testResp](&testProvider{name: "p1", data: "data1"}),
	)
	defer m.Close()

	result1, err := m.Fetch(context.Background(), testReq{Symbol: "test"})
	if err != nil {
		t.Fatalf("First Fetch() error = %v", err)
	}
	if result1.Cached {
		t.Error("First result should not be cached")
	}

	result2, err := m.Fetch(context.Background(), testReq{Symbol: "test"})
	if err != nil {
		t.Fatalf("Second Fetch() error = %v", err)
	}
	if !result2.Cached {
		t.Error("Second result should be cached")
	}
}

func TestPrioritySelector(t *testing.T) {
	selector := NewPrioritySelector()

	providers := []ProviderInfo{
		{Name: "low", Priority: 1},
		{Name: "high", Priority: 10},
		{Name: "mid", Priority: 5},
	}

	names := selector.Select(providers)

	if len(names) != 3 {
		t.Fatalf("Select() returned %d names, want 3", len(names))
	}

	if names[0] != "high" {
		t.Errorf("names[0] = %v, want high", names[0])
	}
	if names[1] != "mid" {
		t.Errorf("names[1] = %v, want mid", names[1])
	}
	if names[2] != "low" {
		t.Errorf("names[2] = %v, want low", names[2])
	}
}

func TestMemoryCache(t *testing.T) {
	cache := NewMemoryCache()

	cache.Set("key1", []byte("value1"), time.Minute)

	data, ok := cache.Get("key1")
	if !ok {
		t.Fatal("Get() returned not ok")
	}
	if string(data) != "value1" {
		t.Errorf("Get() = %v, want value1", string(data))
	}

	_, ok = cache.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for nonexistent key")
	}

	cache.Delete("key1")
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Get() should return false after Delete()")
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	cache := NewMemoryCache()

	cache.Set("key1", []byte("value1"), 10*time.Millisecond)

	data, ok := cache.Get("key1")
	if !ok {
		t.Fatal("Get() returned not ok immediately after Set()")
	}
	if string(data) != "value1" {
		t.Errorf("Get() = %v, want value1", string(data))
	}

	time.Sleep(20 * time.Millisecond)

	_, ok = cache.Get("key1")
	if ok {
		t.Error("Get() should return false after expiration")
	}
}

func TestMemoryCollector(t *testing.T) {
	collector := NewMemoryCollector()

	collector.RecordFetch(Metric{Provider: "p1", Success: true, Duration: time.Second})
	collector.RecordFetch(Metric{Provider: "p1", Success: false, Duration: time.Second})
	collector.RecordFetch(Metric{Provider: "p2", Success: true, CacheHit: true, Duration: time.Second})

	stats := collector.GetStats()
	if stats.TotalFetches != 3 {
		t.Errorf("TotalFetches = %v, want 3", stats.TotalFetches)
	}
	if stats.SuccessFetches != 2 {
		t.Errorf("SuccessFetches = %v, want 2", stats.SuccessFetches)
	}
	if stats.FailedFetches != 1 {
		t.Errorf("FailedFetches = %v, want 1", stats.FailedFetches)
	}
	if stats.CacheHits != 1 {
		t.Errorf("CacheHits = %v, want 1", stats.CacheHits)
	}

	collector.Reset()
	stats = collector.GetStats()
	if stats.TotalFetches != 0 {
		t.Errorf("After Reset(), TotalFetches = %v, want 0", stats.TotalFetches)
	}
}
