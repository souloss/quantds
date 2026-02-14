package yahoo

import (
	"testing"

	"github.com/souloss/quantds/clients/yahoo"
	"github.com/souloss/quantds/domain"
)

func TestNewSpotAdapter(t *testing.T) {
	client := yahoo.NewClient()
	adapter := NewSpotAdapter(client)

	if adapter == nil {
		t.Error("NewSpotAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestSpotAdapter_SupportedMarkets(t *testing.T) {
	client := yahoo.NewClient()
	adapter := NewSpotAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketUS {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketUS, markets)
	}
}

func TestSpotAdapter_CanHandle(t *testing.T) {
	client := yahoo.NewClient()
	adapter := NewSpotAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		// US symbols should match
		{"AAPL.US", true},
		{"MSFT.US", true},
		// Non-US symbols should NOT match
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
