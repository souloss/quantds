package yahoo

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient(nil)
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
		{"60m", Interval60m},
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

func TestParseTimestamp(t *testing.T) {
	// Test with known timestamp: 2024-01-01 00:00:00 UTC = 1704067200
	ts := int64(1704067200)
	result := ParseTimestamp(ts, "America/New_York")

	// Should be 2023-12-31 19:00:00 in New York (UTC-5)
	expectedYear := 2023
	expectedMonth := time.December
	expectedDay := 31

	if result.Year() != expectedYear {
		t.Errorf("ParseTimestamp year = %d, want %d", result.Year(), expectedYear)
	}
	if result.Month() != expectedMonth {
		t.Errorf("ParseTimestamp month = %v, want %v", result.Month(), expectedMonth)
	}
	if result.Day() != expectedDay {
		t.Errorf("ParseTimestamp day = %d, want %d", result.Day(), expectedDay)
	}
}

func TestParseUSSymbol(t *testing.T) {
	tests := []struct {
		input        string
		expectedCode string
		expectedExch string
		expectedOK   bool
	}{
		{"AAPL", "AAPL", "NASDAQ", true},
		{"MSFT.US", "MSFT", "NASDAQ", true},
		{"GOOGL.NASDAQ", "GOOGL", "NASDAQ", true},
		{"JPM.NYSE", "JPM", "NYSE", true},
		{"AAPL.US.NASDAQ", "AAPL", "NASDAQ", true},
		{"BRK.B", "BRK", "B", true}, // Special case for Berkshire
		{"", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			code, exchange, ok := ParseUSSymbol(tt.input)
			if ok != tt.expectedOK {
				t.Errorf("ParseUSSymbol(%s) ok = %v, want %v", tt.input, ok, tt.expectedOK)
			}
			if ok && code != tt.expectedCode {
				t.Errorf("ParseUSSymbol(%s) code = %s, want %s", tt.input, code, tt.expectedCode)
			}
			if ok && exchange != tt.expectedExch {
				t.Errorf("ParseUSSymbol(%s) exchange = %s, want %s", tt.input, exchange, tt.expectedExch)
			}
		})
	}
}

func TestToYahooSymbol(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"AAPL", "AAPL", false},
		{"AAPL.US", "AAPL", false},
		{"GOOGL.NASDAQ", "GOOGL", false},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ToYahooSymbol(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("ToYahooSymbol(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToYahooSymbol(%s) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToYahooSymbol(%s) = %s, want %s", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestIsUSSymbol(t *testing.T) {
	tests := []struct {
		symbol   string
		expected bool
	}{
		{"AAPL", true},
		{"MSFT.US", true},
		{"000001.SZ", false}, // A-share
		{"00700.HK", false},  // HK stock
		{"BTCUSDT", false},   // Crypto
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			result := IsUSSymbol(tt.symbol)
			if result != tt.expected {
				t.Errorf("IsUSSymbol(%s) = %v, want %v", tt.symbol, result, tt.expected)
			}
		})
	}
}
