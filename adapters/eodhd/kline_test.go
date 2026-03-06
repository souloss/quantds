package eodhd

import (
	"testing"

	"github.com/souloss/quantds/clients/eodhd"
	"github.com/souloss/quantds/domain"
)

func TestNewKlineAdapter(t *testing.T) {
	adapter := NewKlineAdapter(eodhd.NewClient())
	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}
	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_SupportedMarkets(t *testing.T) {
	adapter := NewKlineAdapter(eodhd.NewClient())
	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketUS {
		t.Errorf("Expected [%s], got %v", domain.MarketUS, markets)
	}
}

func TestKlineAdapter_CanHandle(t *testing.T) {
	adapter := NewKlineAdapter(eodhd.NewClient())
	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"AAPL.US", true},
		{"000001.SZ", false},
		{"00700.HK.HKEX", false},
	}
	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			if got := adapter.CanHandle(tt.symbol); got != tt.canHandle {
				t.Errorf("CanHandle(%s) = %v, want %v", tt.symbol, got, tt.canHandle)
			}
		})
	}
}
