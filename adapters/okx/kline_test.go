package okx

import (
	"testing"

	okxclient "github.com/souloss/quantds/clients/okx"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
)

func TestNewKlineAdapter(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_Name(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_SupportedMarkets(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewKlineAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCrypto {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCrypto, markets)
	}
}

func TestKlineAdapter_CanHandle(t *testing.T) {
	client := okxclient.NewClient()
	adapter := NewKlineAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"BTCUSDT", true},
		{"ETHUSDT", true},
		{"BTC-USDT", true}, // OKX format
		{"ETH-USDC", true}, // OKX format
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

func TestToOKXInstID(t *testing.T) {
	tests := []struct {
		name     string
		symbol   string
		expected string
	}{
		{"Already OKX format", "BTC-USDT", "BTC-USDT"},
		{"Convert BTCUSDT", "BTCUSDT", "BTC-USDT"},
		{"Convert ETHUSDT", "ETHUSDT", "ETH-USDT"},
		{"Convert with USDC", "ETHUSDC", "ETH-USDC"},
		{"Convert with USD", "BTCUSD", "BTC-USD"},
		{"Convert with BTC quote", "ETHBTC", "ETH-BTC"},
		{"Convert with ETH quote", "BNBETH", "BNB-ETH"},
		{"Lowercase input", "btcusdt", "BTC-USDT"},
		{"Unknown format", "UNKNOWN", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toOKXInstID(tt.symbol)
			if result != tt.expected {
				t.Errorf("toOKXInstID(%s) = %s, want %s", tt.symbol, result, tt.expected)
			}
		})
	}
}

func TestFromOKXInstID(t *testing.T) {
	tests := []struct {
		name     string
		instID   string
		expected string
	}{
		{"BTC-USDT", "BTC-USDT", "BTCUSDT"},
		{"ETH-USDC", "ETH-USDC", "ETHUSDC"},
		{"Already domain format", "BTCUSDT", "BTCUSDT"},
		{"Multiple dashes", "BTC-USDT-SWAP", "BTCUSDTSWAP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fromOKXInstID(tt.instID)
			if result != tt.expected {
				t.Errorf("fromOKXInstID(%s) = %s, want %s", tt.instID, result, tt.expected)
			}
		})
	}
}

func TestToOKXBar(t *testing.T) {
	tests := []struct {
		name      string
		timeframe kline.Timeframe
		expected  string
	}{
		{"1 minute", kline.Timeframe1m, "1m"},
		{"5 minutes", kline.Timeframe5m, "5m"},
		{"15 minutes", kline.Timeframe15m, "15m"},
		{"30 minutes", kline.Timeframe30m, "30m"},
		{"60 minutes", kline.Timeframe60m, "1H"},
		{"1 day", kline.Timeframe1d, "1D"},
		{"1 week", kline.Timeframe1w, "1W"},
		{"1 month", kline.Timeframe1M, "1M"},
		{"Unknown timeframe", kline.Timeframe("unknown"), "1D"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toOKXBar(tt.timeframe)
			if result != tt.expected {
				t.Errorf("toOKXBar(%s) = %s, want %s", tt.timeframe, result, tt.expected)
			}
		})
	}
}
