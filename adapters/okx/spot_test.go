package okx

import (
	"testing"

	okxclient "github.com/souloss/quantds/clients/okx"
	"github.com/souloss/quantds/domain"
)

func TestNewSpotAdapter(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewSpotAdapter(client)

	if adapter == nil {
		t.Error("NewSpotAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestSpotAdapter_Name(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewSpotAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestSpotAdapter_SupportedMarkets(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewSpotAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCrypto {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCrypto, markets)
	}
}

func TestSpotAdapter_CanHandle(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewSpotAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"BTCUSDT", true},
		{"ETHUSDT", true},
		{"BTC-USDT", true}, // Domain.Symbol.Parse normalizes dashes
		{"ETH-USDC", true}, // Domain.Symbol.Parse normalizes dashes
		{"000001.SZ", false},
		{"AAPL.US", false},
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
