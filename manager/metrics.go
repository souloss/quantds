package manager

import (
	"sync"
	"sync/atomic"
	"time"
)

type Metric struct {
	Provider  string
	Duration  time.Duration
	Success   bool
	CacheHit  bool
	ErrorType string
}

type Collector interface {
	RecordFetch(metric Metric)
	RecordRequest(metric Metric)
	GetStats() Stats
	Reset()
}

type Stats struct {
	TotalFetches   int64
	SuccessFetches int64
	FailedFetches  int64
	CacheHits      int64
	AvgLatency     time.Duration

	ByProvider map[string]ProviderMetric
}

type ProviderMetric struct {
	Fetches  int64
	Success  int64
	Failed   int64
	Duration int64
}

type MemoryCollector struct {
	totalFetches   int64
	successFetches int64
	failedFetches  int64
	cacheHits      int64
	totalLatency   int64

	mu         sync.RWMutex
	byProvider map[string]*ProviderMetric
}

func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{
		byProvider: make(map[string]*ProviderMetric),
	}
}

func (c *MemoryCollector) RecordFetch(metric Metric) {
	atomic.AddInt64(&c.totalFetches, 1)
	atomic.AddInt64(&c.totalLatency, int64(metric.Duration))

	if metric.Success {
		atomic.AddInt64(&c.successFetches, 1)
	} else {
		atomic.AddInt64(&c.failedFetches, 1)
	}

	if metric.CacheHit {
		atomic.AddInt64(&c.cacheHits, 1)
	}

	if metric.Provider != "" {
		c.mu.Lock()
		stats, ok := c.byProvider[metric.Provider]
		if !ok {
			stats = &ProviderMetric{}
			c.byProvider[metric.Provider] = stats
		}
		c.mu.Unlock()

		atomic.AddInt64(&stats.Fetches, 1)
		atomic.AddInt64(&stats.Duration, int64(metric.Duration))
		if metric.Success {
			atomic.AddInt64(&stats.Success, 1)
		} else {
			atomic.AddInt64(&stats.Failed, 1)
		}
	}
}

func (c *MemoryCollector) RecordRequest(metric Metric) {
}

func (c *MemoryCollector) GetStats() Stats {
	stats := Stats{
		TotalFetches:   atomic.LoadInt64(&c.totalFetches),
		SuccessFetches: atomic.LoadInt64(&c.successFetches),
		FailedFetches:  atomic.LoadInt64(&c.failedFetches),
		CacheHits:      atomic.LoadInt64(&c.cacheHits),
		ByProvider:     make(map[string]ProviderMetric),
	}

	totalLatency := atomic.LoadInt64(&c.totalLatency)
	if stats.TotalFetches > 0 {
		stats.AvgLatency = time.Duration(totalLatency / stats.TotalFetches)
	}

	c.mu.RLock()
	for name, ps := range c.byProvider {
		stats.ByProvider[name] = ProviderMetric{
			Fetches:  atomic.LoadInt64(&ps.Fetches),
			Success:  atomic.LoadInt64(&ps.Success),
			Failed:   atomic.LoadInt64(&ps.Failed),
			Duration: atomic.LoadInt64(&ps.Duration),
		}
	}
	c.mu.RUnlock()

	return stats
}

func (c *MemoryCollector) Reset() {
	atomic.StoreInt64(&c.totalFetches, 0)
	atomic.StoreInt64(&c.successFetches, 0)
	atomic.StoreInt64(&c.failedFetches, 0)
	atomic.StoreInt64(&c.cacheHits, 0)
	atomic.StoreInt64(&c.totalLatency, 0)

	c.mu.Lock()
	c.byProvider = make(map[string]*ProviderMetric)
	c.mu.Unlock()
}

type NoopCollector struct{}

func NewNoopCollector() *NoopCollector {
	return &NoopCollector{}
}

func (c *NoopCollector) RecordFetch(metric Metric)   {}
func (c *NoopCollector) RecordRequest(metric Metric) {}
func (c *NoopCollector) GetStats() Stats {
	return Stats{ByProvider: make(map[string]ProviderMetric)}
}
func (c *NoopCollector) Reset() {}

var _ Collector = (*MemoryCollector)(nil)
var _ Collector = (*NoopCollector)(nil)
