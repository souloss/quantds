package binance

import (
	"testing"

	"github.com/souloss/quantds/clients/binance"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
)

func TestNewKlineAdapter(t *testing.T) {
	client := binance.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_Name(t *testing.T) {
	client := binance.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_SupportedMarkets(t *testing.T) {
	client := binance.NewClient()
	adapter := NewKlineAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCrypto {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCrypto, markets)
	}
}

func TestKlineAdapter_CanHandle(t *testing.T) {
	client := binance.NewClient()
	adapter := NewKlineAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"BTCUSDT", true},
		{"ETHUSDT", true},
		{"BTC.USDT.CRYPTO.BINANCE", true},
		{"BTC-USDT", true},
		{"000001.SZ", false}, // A-share
		{"00700.HK", false},  // HK stock
		{"AAPL.US", false},   // US stock
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
	client := binance.NewClient()
	adapter := NewSpotAdapter(client)

	if adapter == nil {
		t.Error("NewSpotAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestSpotAdapter_SupportedMarkets(t *testing.T) {
	client := binance.NewClient()
	adapter := NewSpotAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCrypto {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCrypto, markets)
	}
}

func TestSpotAdapter_CanHandle(t *testing.T) {
	client := binance.NewClient()
	adapter := NewSpotAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"BTCUSDT", true},
		{"ETHUSDT", true},
		{"BTC.USDT.CRYPTO.BINANCE", true},
		{"000001.SZ", false},
		{"00700.HK", false},
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

func TestInferIntervalFromTimeframe(t *testing.T) {
	tests := []struct {
		timeframe kline.Timeframe
	}{
		{kline.Timeframe1m},
		{kline.Timeframe5m},
		{kline.Timeframe15m},
		{kline.Timeframe30m},
		{kline.Timeframe60m},
		{kline.Timeframe1d},
		{kline.Timeframe1w},
		{kline.Timeframe1M},
	}

	for _, tt := range tests {
		t.Run(string(tt.timeframe), func(t *testing.T) {
			result := inferIntervalFromTimeframe(tt.timeframe)
			if result <= 0 {
				t.Errorf("inferIntervalFromTimeframe(%s) returned non-positive duration", tt.timeframe)
			}
		})
	}
}
