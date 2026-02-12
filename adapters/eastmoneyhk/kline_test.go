package eastmoneyhk

import (
	"testing"

	"github.com/souloss/quantds/clients/eastmoneyhk"
	"github.com/souloss/quantds/domain"
)

func TestNewKlineAdapter(t *testing.T) {
	client := eastmoneyhk.NewClient(nil)
	adapter := NewKlineAdapter(client)

	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_Name(t *testing.T) {
	client := eastmoneyhk.NewClient(nil)
	adapter := NewKlineAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_SupportedMarkets(t *testing.T) {
	client := eastmoneyhk.NewClient(nil)
	adapter := NewKlineAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketHK {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketHK, markets)
	}
}

func TestKlineAdapter_CanHandle(t *testing.T) {
	client := eastmoneyhk.NewClient(nil)
	adapter := NewKlineAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"00700.HK", true},
		{"00941.HK.HKEX", true},
		{"00700.HK.HKEX", true},
		{"000001.SZ", false}, // A-share
		{"AAPL.US", false},   // US stock
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

func TestNewSpotAdapter(t *testing.T) {
	client := eastmoneyhk.NewClient(nil)
	adapter := NewSpotAdapter(client)

	if adapter == nil {
		t.Error("NewSpotAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestSpotAdapter_SupportedMarkets(t *testing.T) {
	client := eastmoneyhk.NewClient(nil)
	adapter := NewSpotAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketHK {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketHK, markets)
	}
}

func TestSpotAdapter_CanHandle(t *testing.T) {
	client := eastmoneyhk.NewClient(nil)
	adapter := NewSpotAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"00700.HK", true},
		{"00941.HK.HKEX", true},
		{"000001.SZ", false},
		{"AAPL.US", false},
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

func TestFormatHKSymbol(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"00700", "00700.HK.HKEX"},
		{"00941", "00941.HK.HKEX"},
		{"09988", "09988.HK.HKEX"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := formatHKSymbol(tt.code)
			if result != tt.expected {
				t.Errorf("formatHKSymbol(%s) = %s, want %s", tt.code, result, tt.expected)
			}
		})
	}
}
