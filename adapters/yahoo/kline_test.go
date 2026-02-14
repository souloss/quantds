package yahoo

import (
	"testing"

	"github.com/souloss/quantds/clients/yahoo"
	"github.com/souloss/quantds/domain"
)

func TestNewKlineAdapter(t *testing.T) {
	client := yahoo.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_Name(t *testing.T) {
	client := yahoo.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_SupportedMarkets(t *testing.T) {
	client := yahoo.NewClient()
	adapter := NewKlineAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketUS {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketUS, markets)
	}
}

func TestKlineAdapter_CanHandle(t *testing.T) {
	client := yahoo.NewClient()
	adapter := NewKlineAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"AAPL.US", true},
		{"MSFT.US.NASDAQ", true},
		{"JPM.US.NYSE", true},
		{"000001.SZ", false}, // A-share
		{"00700.HK.HKEX", false}, // HK stock
		{"BTCUSDT", false},   // Crypto
		{"invalid", false},
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
