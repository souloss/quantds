package middleware

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/souloss/quantds/manager"
)

// rateLimiter implements a token bucket rate limiter.
type rateLimiter struct {
	mu         sync.Mutex
	rate       float64   // tokens per second
	burst      int       // max tokens
	tokens     float64   // current tokens
	lastRefill time.Time // last refill time
	jitter     float64   // jitter fraction (0.0 - 1.0), e.g. 0.1 = ±10%
}

func newRateLimiter(requestsPerMinute float64, burst int, jitter float64) *rateLimiter {
	return &rateLimiter{
		rate:       requestsPerMinute / 60.0, // convert to per-second
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: time.Now(),
		jitter:     jitter,
	}
}

// allow tries to consume one token. Returns true if allowed.
func (rl *rateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()
	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

// wait blocks until a token is available or context is cancelled.
func (rl *rateLimiter) wait(ctx context.Context) error {
	for {
		rl.mu.Lock()
		rl.refill()
		if rl.tokens >= 1 {
			rl.tokens--
			rl.mu.Unlock()
			return nil
		}
		rl.mu.Unlock()

		// Calculate wait time
		waitTime := rl.waitTime()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			// try again
		}
	}
}

func (rl *rateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens += elapsed * rl.rate
	if rl.tokens > float64(rl.burst) {
		rl.tokens = float64(rl.burst)
	}
	rl.lastRefill = now
}

func (rl *rateLimiter) waitTime() time.Duration {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	deficit := 1.0 - rl.tokens
	if deficit <= 0 {
		return 0
	}
	base := time.Duration(deficit/rl.rate) * time.Second
	if rl.jitter > 0 {
		// Add random jitter: ±jitter fraction of the base wait time
		jitterNs := float64(base) * rl.jitter * (2*rand.Float64() - 1)
		base = time.Duration(float64(base) + jitterNs)
		if base < 0 {
			base = 0
		}
	}
	return base
}

// RateLimiter creates a rate limiting middleware with the given configuration.
// requestsPerMinute is the sustained rate; burst allows short bursts above that rate.
// jitter is a fraction (0.0-1.0) of the wait time to add as random delay, e.g. 0.1 = ±10%.
// When the rate limit is exceeded, the middleware waits (blocks) until a token
// is available, respecting context cancellation.
func RateLimiter[Req, Resp any](requestsPerMinute float64, burst int, jitter float64) Middleware[Req, Resp] {
	rl := newRateLimiter(requestsPerMinute, burst, jitter)
	return func(next manager.Provider[Req, Resp]) manager.Provider[Req, Resp] {
		return &providerFunc[Req, Resp]{
			name:             next.Name(),
			supportedMarkets: next.SupportedMarkets(),
			canHandle:        next.CanHandle,
			fetch: func(ctx context.Context, req Req) (Resp, *manager.RequestTrace, error) {
				if err := rl.wait(ctx); err != nil {
					var zero Resp
					return zero, nil, err
				}
				return next.Fetch(ctx, req)
			},
		}
	}
}

// RateLimiterConfig holds configuration for rate limiter middleware.
type RateLimiterConfig struct {
	Enabled           bool
	RequestsPerMinute float64
	Burst             int
	Jitter            float64 // jitter fraction (0.0-1.0), e.g. 0.1 = ±10%
}

// RateLimiterFromConfig creates a rate limiter middleware from config.
// If cfg.Enabled is false, returns an identity middleware (no-op).
func RateLimiterFromConfig[Req, Resp any](cfg RateLimiterConfig) Middleware[Req, Resp] {
	if !cfg.Enabled {
		return Identity[Req, Resp]()
	}
	return RateLimiter[Req, Resp](cfg.RequestsPerMinute, cfg.Burst, cfg.Jitter)
}
