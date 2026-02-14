package eastmoneyhk

import (
	"testing"

	"github.com/souloss/quantds/clients/eastmoneyhk"
	"github.com/souloss/quantds/domain"
)

func TestNewInstrumentAdapter(t *testing.T) {
	client := eastmoneyhk.NewClient()
	adapter := NewInstrumentAdapter(client)

	if adapter == nil {
		t.Error("NewInstrumentAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestInstrumentAdapter_SupportedMarkets(t *testing.T) {
	client := eastmoneyhk.NewClient()
	adapter := NewInstrumentAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketHK {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketHK, markets)
	}
}

func TestInstrumentAdapter_CanHandle(t *testing.T) {
	client := eastmoneyhk.NewClient()
	adapter := NewInstrumentAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		// HK symbols should match
		{"00700.HK.HKEX", true},
		{"09988.HK.HKEX", true},
		// Non-HK symbols should NOT match
		{"000001.SZ", false},
		{"BTCUSDT", false},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			result := adapter.CanHandle(tt.symbol)
			if result != tt.canHandle {
				t.Errorf("CanHandle(%s) = %v, want %v", tt.symbol, result, tt.canHandle)
			}
		})
	}
}
