package finnhub

import (
	"testing"

	"github.com/souloss/quantds/clients/finnhub"
	"github.com/souloss/quantds/domain"
)

func TestNewInstrumentAdapter(t *testing.T) {
	adapter := NewInstrumentAdapter(finnhub.NewClient())
	if adapter == nil {
		t.Error("NewInstrumentAdapter returned nil")
	}
	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestInstrumentAdapter_SupportedMarkets(t *testing.T) {
	adapter := NewInstrumentAdapter(finnhub.NewClient())
	markets := adapter.SupportedMarkets()
	found := map[domain.Market]bool{}
	for _, m := range markets {
		found[m] = true
	}
	if !found[domain.MarketUS] {
		t.Error("Expected MarketUS in supported markets")
	}
}

func TestInstrumentAdapter_CanHandle(t *testing.T) {
	adapter := NewInstrumentAdapter(finnhub.NewClient())
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
