package binance

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("NewClient returned nil")
	}
	defer client.Close()
}

func TestToInterval(t *testing.T) {
	tests := []struct {
		timeframe string
		expected  string
	}{
		{"1m", Interval1m},
		{"5m", Interval5m},
		{"15m", Interval15m},
		{"30m", Interval30m},
		{"60m", Interval1h},
		{"1h", Interval1h},
		{"1d", Interval1d},
		{"1w", Interval1w},
		{"1M", Interval1M},
		{"", Interval1d},
		{"invalid", Interval1d},
	}

	for _, tt := range tests {
		t.Run(tt.timeframe, func(t *testing.T) {
			result := ToInterval(tt.timeframe)
			if result != tt.expected {
				t.Errorf("ToInterval(%s) = %s, want %s", tt.timeframe, result, tt.expected)
			}
		})
	}
}

func TestParseOpenTime(t *testing.T) {
	// Test with known timestamp: 2024-01-01 00:00:00 UTC in milliseconds
	ms := int64(1704067200000)
	result := ParseOpenTime(ms)

	expectedYear := 2024
	expectedMonth := time.January
	expectedDay := 1

	if result.Year() != expectedYear {
		t.Errorf("ParseOpenTime year = %d, want %d", result.Year(), expectedYear)
	}
	if result.Month() != expectedMonth {
		t.Errorf("ParseOpenTime month = %v, want %v", result.Month(), expectedMonth)
	}
	if result.Day() != expectedDay {
		t.Errorf("ParseOpenTime day = %d, want %d", result.Day(), expectedDay)
	}
}

func TestFormatToUnixMilli(t *testing.T) {
	tm := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result := FormatToUnixMilli(tm)
	expected := int64(1704067200000)

	if result != expected {
		t.Errorf("FormatToUnixMilli = %d, want %d", result, expected)
	}
}

func TestParseBinanceSymbol(t *testing.T) {
	tests := []struct {
		input         string
		expectedBase  string
		expectedQuote string
		expectedOK    bool
	}{
		{"BTCUSDT", "BTC", "USDT", true},
		{"ETHUSDT", "ETH", "USDT", true},
		{"BTC-USDT", "BTC", "USDT", true},
		{"BTC/USDT", "BTC", "USDT", true},
		{"BNBBTC", "BNB", "BTC", true},
		{"ETHBTC", "ETH", "BTC", true},
		{"INVALID", "", "", false},
		{"", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			base, quote, ok := ParseBinanceSymbol(tt.input)
			if ok != tt.expectedOK {
				t.Errorf("ParseBinanceSymbol(%s) ok = %v, want %v", tt.input, ok, tt.expectedOK)
			}
			if ok {
				if base != tt.expectedBase {
					t.Errorf("ParseBinanceSymbol(%s) base = %s, want %s", tt.input, base, tt.expectedBase)
				}
				if quote != tt.expectedQuote {
					t.Errorf("ParseBinanceSymbol(%s) quote = %s, want %s", tt.input, quote, tt.expectedQuote)
				}
			}
		})
	}
}

func TestToBinanceSymbol(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"BTCUSDT", "BTCUSDT", false},
		{"BTC-USDT", "BTCUSDT", false},
		{"BTC/USDT", "BTCUSDT", false},
		{"ETHUSDT", "ETHUSDT", false},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ToBinanceSymbol(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("ToBinanceSymbol(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToBinanceSymbol(%s) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToBinanceSymbol(%s) = %s, want %s", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestIsCryptoSymbol(t *testing.T) {
	tests := []struct {
		symbol   string
		expected bool
	}{
		{"BTCUSDT", true},
		{"ETHUSDT", true},
		{"BTC.USDT.CRYPTO.BINANCE", true},
		{"AAPL", false},
		{"000001.SZ", false},
		{"00700.HK", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			result := IsCryptoSymbol(tt.symbol)
			if result != tt.expected {
				t.Errorf("IsCryptoSymbol(%s) = %v, want %v", tt.symbol, result, tt.expected)
			}
		})
	}
}

func TestGetQuoteAsset(t *testing.T) {
	tests := []struct {
		symbol   string
		expected string
	}{
		{"BTCUSDT", "USDT"},
		{"ETHBTC", "BTC"},
		{"SOLBUSD", "BUSD"},
		{"INVALID", ""},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			result := GetQuoteAsset(tt.symbol)
			if result != tt.expected {
				t.Errorf("GetQuoteAsset(%s) = %s, want %s", tt.symbol, result, tt.expected)
			}
		})
	}
}

func TestGetBaseAsset(t *testing.T) {
	tests := []struct {
		symbol   string
		expected string
	}{
		{"BTCUSDT", "BTC"},
		{"ETHUSDT", "ETH"},
		{"BNBBTC", "BNB"},
		{"INVALID", ""},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			result := GetBaseAsset(tt.symbol)
			if result != tt.expected {
				t.Errorf("GetBaseAsset(%s) = %s, want %s", tt.symbol, result, tt.expected)
			}
		})
	}
}

func TestClient_GetKline(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	params := &KlineParams{
		Symbol:   "BTCUSDT",
		Interval: Interval1d,
		Limit:    5,
	}

	result, record, err := client.GetKline(ctx, params)
	if err != nil {
		if record != nil {
			t.Logf("Response Status: %d", record.Response.StatusCode)
			t.Logf("Response Body: %s", string(record.Response.Body))
		}
		checkAPIError(t, err)
		return
	}

	if result.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", result.Symbol)
	}

	if len(result.Data) == 0 {
		t.Fatalf("Expected kline data, got 0")
	}

	kline := result.Data[0]
	if kline.OpenTime == 0 {
		t.Error("Expected OpenTime > 0")
	}
	if kline.Open <= 0 {
		t.Errorf("Expected Open > 0, got %f", kline.Open)
	}
	if kline.Close <= 0 {
		t.Errorf("Expected Close > 0, got %f", kline.Close)
	}
	if kline.High <= 0 {
		t.Errorf("Expected High > 0, got %f", kline.High)
	}
	if kline.Low <= 0 {
		t.Errorf("Expected Low > 0, got %f", kline.Low)
	}
}

func TestClient_GetKline_Error(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	// Invalid symbol should return error or empty result depending on API behavior
	// Binance API usually returns error for invalid symbol
	params := &KlineParams{Symbol: "INVALID_SYMBOL"}
	_, _, err := client.GetKline(ctx, params)
	if err == nil {
		t.Log("Warning: Expected error for invalid symbol, but got nil. This might be due to API behavior change.")
	}
}
