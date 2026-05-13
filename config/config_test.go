package config

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.HTTP.Timeout != 30*time.Second {
		t.Errorf("Expected HTTP timeout 30s, got %v", cfg.HTTP.Timeout)
	}
	if cfg.Cache.DefaultTTL != 5*time.Minute {
		t.Errorf("Expected cache TTL 5m, got %v", cfg.Cache.DefaultTTL)
	}
	if !cfg.Providers.Tushare.Enabled {
		t.Error("Tushare should be enabled by default")
	}
	if cfg.Providers.AlphaVantage.RateLimit.RequestsPerMinute != 5 {
		t.Errorf("AlphaVantage rate limit should be 5/min, got %v", cfg.Providers.AlphaVantage.RateLimit.RequestsPerMinute)
	}
	if cfg.Providers.Binance.RateLimit.RequestsPerMinute != 1200 {
		t.Errorf("Binance rate limit should be 1200/min, got %v", cfg.Providers.Binance.RateLimit.RequestsPerMinute)
	}
}

func TestGetProvider(t *testing.T) {
	cfg := DefaultConfig()
	p, err := cfg.GetProvider("tushare")
	if err != nil {
		t.Fatalf("GetProvider(tushare) error: %v", err)
	}
	if !p.Enabled {
		t.Error("Tushare should be enabled")
	}
	_, err = cfg.GetProvider("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent provider")
	}
}

func TestEnabledProviders(t *testing.T) {
	cfg := DefaultConfig()
	enabled := cfg.EnabledProviders()
	if len(enabled) == 0 {
		t.Error("Expected some enabled providers")
	}
}

func TestProviderNames(t *testing.T) {
	cfg := DefaultConfig()
	names := cfg.ProviderNames()
	if len(names) != 16 {
		t.Errorf("Expected 16 providers, got %d", len(names))
	}
}
