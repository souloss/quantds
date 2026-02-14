package bse

import (
	"testing"

	"github.com/souloss/quantds/clients/bse"
	"github.com/souloss/quantds/domain"
)

func TestNewInstrumentAdapter(t *testing.T) {
	client := bse.NewClient()
	adapter := NewInstrumentAdapter(client)

	if adapter == nil {
		t.Error("NewInstrumentAdapter returned nil")
	}

	if adapter.Name() != "bse" {
		t.Errorf("Expected name 'bse', got '%s'", adapter.Name())
	}
}

func TestInstrumentAdapter_Name(t *testing.T) {
	client := bse.NewClient()
	adapter := NewInstrumentAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestInstrumentAdapter_SupportedMarkets(t *testing.T) {
	client := bse.NewClient()
	adapter := NewInstrumentAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCN {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCN, markets)
	}
}

func TestInstrumentAdapter_CanHandle(t *testing.T) {
	client := bse.NewClient()
	adapter := NewInstrumentAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"000001.SZ", true},
		{"600001.SH", true},
		{"430001.BJ", true},
		{"00700.HK", false},
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
