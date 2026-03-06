package twelvedata

import (
	"testing"

	"github.com/souloss/quantds/clients/twelvedata"
	"github.com/souloss/quantds/domain"
)

func TestNewKlineAdapter(t *testing.T) {
	adapter := NewKlineAdapter(twelvedata.NewClient())
	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}
	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_SupportedMarkets(t *testing.T) {
	adapter := NewKlineAdapter(twelvedata.NewClient())
	found := map[domain.Market]bool{}
	for _, m := range adapter.SupportedMarkets() {
		found[m] = true
	}
	for _, want := range []domain.Market{domain.MarketUS, domain.MarketForex, domain.MarketCrypto} {
		if !found[want] {
			t.Errorf("Expected market %s", want)
		}
	}
}

func TestKlineAdapter_CanHandle(t *testing.T) {
	adapter := NewKlineAdapter(twelvedata.NewClient())
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
