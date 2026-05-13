package middleware

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/manager"
)

func TestRateLimiter_AllowsBurst(t *testing.T) {
	// 60 req/min = 1 req/sec, burst of 5
	rl := RateLimiter[string, string](60, 5, 0)

	var callCount atomic.Int32
	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			callCount.Add(1)
			return "ok", nil, nil
		},
	}

	wrapped := rl(next)

	// Should allow burst of 5 immediately
	for i := 0; i < 5; i++ {
		resp, _, err := wrapped.Fetch(context.Background(), "req")
		if err != nil {
			t.Errorf("burst request %d: unexpected error: %v", i, err)
		}
		if resp != "ok" {
			t.Errorf("burst request %d: got %q, want %q", i, resp, "ok")
		}
	}
	if callCount.Load() != 5 {
		t.Errorf("expected 5 calls, got %d", callCount.Load())
	}
}

func TestRateLimiter_BlocksAfterBurst(t *testing.T) {
	// 60 req/min = 1 req/sec, burst of 2
	rl := RateLimiter[string, string](60, 2, 0)

	next := &mockProvider[string, string]{
		name:  "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "ok", nil, nil },
	}

	wrapped := rl(next)

	// Consume burst
	wrapped.Fetch(context.Background(), "req")
	wrapped.Fetch(context.Background(), "req")

	// Next request should take ~1 second (rate limited)
	start := time.Now()
	resp, _, err := wrapped.Fetch(context.Background(), "req")
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("rate-limited request: unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Errorf("rate-limited request: got %q, want %q", resp, "ok")
	}
	// Should have waited approximately 1 second (allow some tolerance)
	if elapsed < 500*time.Millisecond {
		t.Errorf("expected rate limiting delay, got %v", elapsed)
	}
}

func TestRateLimiter_RespectsContext(t *testing.T) {
	// 1 req/min, burst of 1 — very slow rate
	rl := RateLimiter[string, string](1, 1, 0)

	next := &mockProvider[string, string]{
		name:  "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "ok", nil, nil },
	}

	wrapped := rl(next)

	// Consume the burst
	wrapped.Fetch(context.Background(), "req")

	// Cancel context quickly
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, _, err := wrapped.Fetch(ctx, "req")
	if err == nil {
		t.Error("expected context cancellation error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestRateLimiter_PreservesProviderInfo(t *testing.T) {
	rl := RateLimiter[string, string](60, 5, 0)

	next := &mockProvider[string, string]{
		name:    "myprovider",
		fetch:   func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "ok", nil, nil },
		markets: []domain.Market{domain.MarketCN},
	}

	wrapped := rl(next)

	if wrapped.Name() != "myprovider" {
		t.Errorf("Name() = %q, want %q", wrapped.Name(), "myprovider")
	}
	if !wrapped.CanHandle("anything") {
		t.Error("CanHandle() = false, want true")
	}
}

func TestRateLimiterFromConfig_Disabled(t *testing.T) {
	cfg := RateLimiterConfig{Enabled: false}
	rl := RateLimiterFromConfig[string, string](cfg)

	next := &mockProvider[string, string]{
		name:  "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "ok", nil, nil },
	}

	wrapped := rl(next)

	resp, _, err := wrapped.Fetch(context.Background(), "req")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Errorf("got %q, want %q", resp, "ok")
	}
}

func TestRateLimiterFromConfig_Enabled(t *testing.T) {
	cfg := RateLimiterConfig{
		Enabled:           true,
		RequestsPerMinute: 60,
		Burst:             5,
	}
	rl := RateLimiterFromConfig[string, string](cfg)

	var callCount atomic.Int32
	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			callCount.Add(1)
			return "ok", nil, nil
		},
	}

	wrapped := rl(next)

	// Should allow burst
	for i := 0; i < 5; i++ {
		resp, _, err := wrapped.Fetch(context.Background(), "req")
		if err != nil {
			t.Errorf("request %d: unexpected error: %v", i, err)
		}
		if resp != "ok" {
			t.Errorf("request %d: got %q, want %q", i, resp, "ok")
		}
	}
	if callCount.Load() != 5 {
		t.Errorf("expected 5 calls, got %d", callCount.Load())
	}
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	// High rate to avoid blocking, burst of 10
	rl := RateLimiter[string, string](6000, 10, 0)

	var callCount atomic.Int32
	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			callCount.Add(1)
			return "ok", nil, nil
		},
	}

	wrapped := rl(next)

	// Launch 10 concurrent requests
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			wrapped.Fetch(context.Background(), "req")
			done <- struct{}{}
		}()
	}

	// Wait for all to complete with timeout
	timeout := time.After(5 * time.Second)
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-timeout:
			t.Fatal("timeout waiting for concurrent requests")
		}
	}

	if callCount.Load() != 10 {
		t.Errorf("expected 10 calls, got %d", callCount.Load())
	}
}
