package okx

import (
	"testing"

	okxclient "github.com/souloss/quantds/clients/okx"
	"github.com/souloss/quantds/domain"
)

func TestNewInstrumentAdapter(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewInstrumentAdapter(client)

	if adapter == nil {
		t.Error("NewInstrumentAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestInstrumentAdapter_Name(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewInstrumentAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestInstrumentAdapter_SupportedMarkets(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewInstrumentAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCrypto {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCrypto, markets)
	}
}

func TestInstrumentAdapter_CanHandle(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewInstrumentAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"BTCUSDT", true},
		{"ETHUSDT", true},
		{"BTC-USDT", true},
		{"000001.SZ", true}, // InstrumentAdapter always returns true
		{"AAPL.US", true},   // InstrumentAdapter always returns true
		{"", true},          // InstrumentAdapter always returns true
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
