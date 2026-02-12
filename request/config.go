package request

import (
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/ratelimiter"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/failsafe-go/failsafe-go/timeout"
)

type Config struct {
	RetryPolicy    failsafe.Policy[Response]
	CircuitBreaker failsafe.Policy[Response]
	RateLimiter    failsafe.Policy[Response]
	Timeout        failsafe.Policy[Response]
}

type ConfigOption func(*Config)

func WithRetryPolicy(policy failsafe.Policy[Response]) ConfigOption {
	return func(c *Config) {
		c.RetryPolicy = policy
	}
}

func WithCircuitBreaker(policy failsafe.Policy[Response]) ConfigOption {
	return func(c *Config) {
		c.CircuitBreaker = policy
	}
}

func WithRateLimiter(policy failsafe.Policy[Response]) ConfigOption {
	return func(c *Config) {
		c.RateLimiter = policy
	}
}

func WithTimeout(policy failsafe.Policy[Response]) ConfigOption {
	return func(c *Config) {
		c.Timeout = policy
	}
}

func DefaultConfig(opts ...ConfigOption) *Config {
	cfg := &Config{
		RetryPolicy:    DefaultRetryPolicy(),
		CircuitBreaker: DefaultCircuitBreaker(),
		Timeout:        DefaultTimeout(),
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func DefaultRetryPolicy() failsafe.Policy[Response] {
	return retrypolicy.NewBuilder[Response]().
		HandleIf(func(resp Response, err error) bool {
			return IsRetryableError(err)
		}).
		WithMaxRetries(3).
		WithBackoff(time.Second, 30*time.Second).
		Build()
}

func DefaultCircuitBreaker() failsafe.Policy[Response] {
	return circuitbreaker.NewBuilder[Response]().
		WithFailureRateThreshold(0.5, 5, time.Minute).
		WithDelay(30 * time.Second).
		WithSuccessThreshold(3).
		Build()
}

func DefaultTimeout() failsafe.Policy[Response] {
	return timeout.New[Response](30 * time.Second)
}

func DefaultRateLimiter() failsafe.Policy[Response] {
	return ratelimiter.NewBursty[Response](100, time.Second)
}
