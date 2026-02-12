package manager

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/souloss/quantds/request"
)

type FetchResult[Resp any] struct {
	Data     Resp
	Trace    *RequestTrace
	Provider string
	Cached   bool
}

func (r *FetchResult[Resp]) IsCached() bool {
	return r.Cached
}

func (r *FetchResult[Resp]) HasError() bool {
	if r.Trace == nil {
		return false
	}
	return r.Trace.FailedRequests() > 0
}

type Manager[Req, Resp any] struct {
	mu           sync.RWMutex
	providers    map[string]Provider[Req, Resp]
	providerInfo map[string]ProviderInfo
	client       request.Client
	cache        *TwoLevelCache
	metrics      Collector
	selector     Selector
}

type ManagerOption[Req, Resp any] func(*Manager[Req, Resp])

func WithClient[Req, Resp any](client request.Client) ManagerOption[Req, Resp] {
	return func(m *Manager[Req, Resp]) {
		m.client = client
	}
}

func WithTwoLevelCache[Req, Resp any](requestTTL, fetchTTL time.Duration) ManagerOption[Req, Resp] {
	return func(m *Manager[Req, Resp]) {
		m.cache = NewTwoLevelCache(requestTTL, fetchTTL)
	}
}

func WithMetrics[Req, Resp any](collector Collector) ManagerOption[Req, Resp] {
	return func(m *Manager[Req, Resp]) {
		m.metrics = collector
	}
}

func WithSelector[Req, Resp any](selector Selector) ManagerOption[Req, Resp] {
	return func(m *Manager[Req, Resp]) {
		m.selector = selector
	}
}

func WithProvider[Req, Resp any](p Provider[Req, Resp], opts ...ProviderOption) ManagerOption[Req, Resp] {
	return func(m *Manager[Req, Resp]) {
		m.Register(p, opts...)
	}
}

func NewManager[Req, Resp any](opts ...ManagerOption[Req, Resp]) *Manager[Req, Resp] {
	m := &Manager[Req, Resp]{
		providers:    make(map[string]Provider[Req, Resp]),
		providerInfo: make(map[string]ProviderInfo),
		selector:     NewPrioritySelector(),
		metrics:      NewNoopCollector(),
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.client == nil {
		m.client = request.NewClient(request.DefaultConfig())
	}

	return m
}

func (m *Manager[Req, Resp]) Register(p Provider[Req, Resp], opts ...ProviderOption) {
	m.mu.Lock()
	defer m.mu.Unlock()

	info := ProviderInfo{
		Name:     p.Name(),
		Priority: 0,
		Weight:   1,
		Tags:     make(map[string]string),
	}

	for _, opt := range opts {
		opt(&info)
	}

	m.providers[p.Name()] = p
	m.providerInfo[p.Name()] = info
}

func (m *Manager[Req, Resp]) Fetch(ctx context.Context, req Req) (*FetchResult[Resp], error) {
	startTime := time.Now()

	if m.cache != nil {
		cacheKey := BuildCacheKey(req)
		if data, ok := m.cache.GetFetch(cacheKey); ok {
			var result FetchResult[Resp]
			if err := json.Unmarshal(data, &result); err == nil {
				result.Cached = true
				m.metrics.RecordFetch(Metric{
					Provider: "cache",
					Duration: time.Since(startTime),
					Success:  true,
					CacheHit: true,
				})
				return &result, nil
			}
		}
	}

	providerNames := m.getOrderedProviders()
	if len(providerNames) == 0 {
		return nil, ErrNoProvider
	}

	var lastErr error
	for _, name := range providerNames {
		m.mu.RLock()
		provider, ok := m.providers[name]
		m.mu.RUnlock()

		if !ok {
			continue
		}

		resp, trace, err := provider.Fetch(ctx, m.client, req)
		if err != nil {
			lastErr = err
			m.metrics.RecordFetch(Metric{
				Provider:  name,
				Duration:  time.Since(startTime),
				Success:   false,
				ErrorType: "fetch_error",
			})
			continue
		}

		result := &FetchResult[Resp]{
			Data:     resp,
			Trace:    trace,
			Provider: name,
			Cached:   false,
		}

		if m.cache != nil {
			cacheKey := BuildCacheKey(req)
			if data, merr := json.Marshal(result); merr == nil {
				m.cache.SetFetch(cacheKey, data)
			}
		}

		m.metrics.RecordFetch(Metric{
			Provider: name,
			Duration: time.Since(startTime),
			Success:  true,
		})

		return result, nil
	}

	return nil, errors.Join(ErrAllProviderFailed, lastErr)
}

func (m *Manager[Req, Resp]) FetchFrom(ctx context.Context, providerName string, req Req) (*FetchResult[Resp], error) {
	m.mu.RLock()
	provider, ok := m.providers[providerName]
	m.mu.RUnlock()

	if !ok {
		return nil, ErrNoProvider
	}

	startTime := time.Now()
	resp, trace, err := provider.Fetch(ctx, m.client, req)
	if err != nil {
		m.metrics.RecordFetch(Metric{
			Provider:  providerName,
			Duration:  time.Since(startTime),
			Success:   false,
			ErrorType: "fetch_error",
		})
		return nil, err
	}

	m.metrics.RecordFetch(Metric{
		Provider: providerName,
		Duration: time.Since(startTime),
		Success:  true,
	})

	return &FetchResult[Resp]{
		Data:     resp,
		Trace:    trace,
		Provider: providerName,
		Cached:   false,
	}, nil
}

func (m *Manager[Req, Resp]) getOrderedProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	providers := make([]ProviderInfo, 0, len(m.providerInfo))
	for _, info := range m.providerInfo {
		providers = append(providers, info)
	}

	return m.selector.Select(providers)
}

func (m *Manager[Req, Resp]) Providers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.providers))
	for name := range m.providers {
		names = append(names, name)
	}
	return names
}

func (m *Manager[Req, Resp]) Stats() Stats {
	return m.metrics.GetStats()
}

func (m *Manager[Req, Resp]) Close() {
	if m.client != nil {
		m.client.Close()
	}
}
