// Package config provides centralized configuration management for quantds.
//
// Configuration sources (in priority order):
// 1. Runtime overrides (code)
// 2. Environment variables
// 3. Configuration file (YAML, optional)
// 4. Code defaults
package config

import (
	"fmt"
	"os"
	"time"
)

// Config is the root configuration for the entire application.
type Config struct {
	HTTP      HTTPConfig
	Cache     CacheConfig
	Providers ProvidersConfig
}

// HTTPConfig holds shared HTTP client settings.
type HTTPConfig struct {
	Timeout      time.Duration
	MaxRetries   int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
	MaxIdleConns int
	UserAgent    string
}

// CacheConfig holds cache settings.
type CacheConfig struct {
	Enabled    bool
	DefaultTTL time.Duration
	MaxEntries int
}

// ProvidersConfig holds per-provider configurations keyed by provider name.
type ProvidersConfig struct {
	// CN data sources
	Tushare   ProviderConfig
	EastMoney ProviderConfig
	Sina      ProviderConfig
	Xueqiu    ProviderConfig
	CNInfo    ProviderConfig
	AKShare   ProviderConfig
	// US data sources
	Yahoo        ProviderConfig
	AlphaVantage ProviderConfig
	Finnhub      ProviderConfig
	Polygon      ProviderConfig
	TwelveData   ProviderConfig
	EODHD        ProviderConfig
	// HK data sources
	EastMoneyHK ProviderConfig
	// Crypto data sources
	Binance   ProviderConfig
	OKX       ProviderConfig
	CoinGecko ProviderConfig
}

// ProviderConfig holds configuration for a single data provider.
type ProviderConfig struct {
	Token          string
	BaseURL        string
	Priority       int
	Enabled        bool
	RateLimit      RateLimitConfig
	CircuitBreaker CircuitBreakerConfig
}

// RateLimitConfig holds rate limiting settings.
type RateLimitConfig struct {
	RequestsPerMinute float64
	Burst             int
	Jitter            float64 // jitter fraction (0.0-1.0), e.g. 0.1 = ±10%
	Enabled           bool
}

// CircuitBreakerConfig holds circuit breaker settings.
type CircuitBreakerConfig struct {
	FailureThreshold int
	SuccessThreshold int
	Timeout          time.Duration
	Enabled          bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		HTTP: HTTPConfig{
			Timeout:      30 * time.Second,
			MaxRetries:   3,
			RetryWaitMin: 1 * time.Second,
			RetryWaitMax: 5 * time.Second,
			MaxIdleConns: 100,
			UserAgent:    "quantds/1.0",
		},
		Cache: CacheConfig{
			Enabled:    true,
			DefaultTTL: 5 * time.Minute,
			MaxEntries: 10000,
		},
		Providers: ProvidersConfig{
			// CN
			Tushare: ProviderConfig{
				Token:          os.Getenv("TUSHARE_TOKEN"),
				Enabled:        true,
				Priority:       1,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 200, Burst: 10, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			EastMoney: ProviderConfig{
				Enabled:        true,
				Priority:       2,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 60, Burst: 5, Jitter: 0.1, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			Sina: ProviderConfig{
				Enabled:        true,
				Priority:       3,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 60, Burst: 5, Jitter: 0.15, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			Xueqiu: ProviderConfig{
				Enabled:        true,
				Priority:       4,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 30, Burst: 3, Jitter: 0.2, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			CNInfo: ProviderConfig{
				Enabled:        true,
				Priority:       5,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 30, Burst: 3, Jitter: 0.1, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			AKShare: ProviderConfig{
				Enabled:        true,
				Priority:       6,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 30, Burst: 3, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			// US
			Yahoo: ProviderConfig{
				Enabled:        true,
				Priority:       1,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 60, Burst: 5, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			AlphaVantage: ProviderConfig{
				Token:          os.Getenv("ALPHAVANTAGE_API_KEY"),
				Enabled:        true,
				Priority:       2,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 5, Burst: 2, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 3, SuccessThreshold: 2, Timeout: 60 * time.Second, Enabled: true},
			},
			Finnhub: ProviderConfig{
				Token:          os.Getenv("FINNHUB_API_KEY"),
				Enabled:        true,
				Priority:       3,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 60, Burst: 5, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			Polygon: ProviderConfig{
				Token:          os.Getenv("POLYGON_API_KEY"),
				Enabled:        true,
				Priority:       4,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 5, Burst: 2, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 3, SuccessThreshold: 2, Timeout: 60 * time.Second, Enabled: true},
			},
			TwelveData: ProviderConfig{
				Token:          os.Getenv("TWELVEDATA_API_KEY"),
				Enabled:        true,
				Priority:       5,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 8, Burst: 3, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 3, SuccessThreshold: 2, Timeout: 60 * time.Second, Enabled: true},
			},
			EODHD: ProviderConfig{
				Token:          os.Getenv("EODHD_API_KEY"),
				Enabled:        true,
				Priority:       6,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 60, Burst: 5, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			// HK
			EastMoneyHK: ProviderConfig{
				Enabled:        true,
				Priority:       1,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 60, Burst: 5, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			// Crypto
			Binance: ProviderConfig{
				Token:          os.Getenv("BINANCE_API_KEY"),
				Enabled:        true,
				Priority:       1,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 1200, Burst: 50, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 10, SuccessThreshold: 5, Timeout: 15 * time.Second, Enabled: true},
			},
			OKX: ProviderConfig{
				Token:          os.Getenv("OKX_API_KEY"),
				Enabled:        true,
				Priority:       2,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 60, Burst: 5, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 5, SuccessThreshold: 3, Timeout: 30 * time.Second, Enabled: true},
			},
			CoinGecko: ProviderConfig{
				Token:          os.Getenv("COINGECKO_API_KEY"),
				Enabled:        true,
				Priority:       3,
				RateLimit:      RateLimitConfig{RequestsPerMinute: 10, Burst: 3, Enabled: true},
				CircuitBreaker: CircuitBreakerConfig{FailureThreshold: 3, SuccessThreshold: 2, Timeout: 60 * time.Second, Enabled: true},
			},
		},
	}
}

// ProviderNames returns all provider names in order.
func (c *Config) ProviderNames() []string {
	return []string{
		"tushare", "eastmoney", "sina", "xueqiu", "cninfo", "akshare",
		"yahoo", "alphavantage", "finnhub", "polygon", "twelvedata", "eodhd",
		"eastmoneyhk",
		"binance", "okx", "coingecko",
	}
}

// GetProvider returns the configuration for a named provider.
func (c *Config) GetProvider(name string) (ProviderConfig, error) {
	switch name {
	case "tushare":
		return c.Providers.Tushare, nil
	case "eastmoney":
		return c.Providers.EastMoney, nil
	case "sina":
		return c.Providers.Sina, nil
	case "xueqiu":
		return c.Providers.Xueqiu, nil
	case "cninfo":
		return c.Providers.CNInfo, nil
	case "akshare":
		return c.Providers.AKShare, nil
	case "yahoo":
		return c.Providers.Yahoo, nil
	case "alphavantage":
		return c.Providers.AlphaVantage, nil
	case "finnhub":
		return c.Providers.Finnhub, nil
	case "polygon":
		return c.Providers.Polygon, nil
	case "twelvedata":
		return c.Providers.TwelveData, nil
	case "eodhd":
		return c.Providers.EODHD, nil
	case "eastmoneyhk":
		return c.Providers.EastMoneyHK, nil
	case "binance":
		return c.Providers.Binance, nil
	case "okx":
		return c.Providers.OKX, nil
	case "coingecko":
		return c.Providers.CoinGecko, nil
	default:
		return ProviderConfig{}, fmt.Errorf("unknown provider: %s", name)
	}
}

// EnabledProviders returns names of all enabled providers.
func (c *Config) EnabledProviders() []string {
	var enabled []string
	for _, name := range c.ProviderNames() {
		p, err := c.GetProvider(name)
		if err == nil && p.Enabled {
			enabled = append(enabled, name)
		}
	}
	return enabled
}
