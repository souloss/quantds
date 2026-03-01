package eastmoneyhk

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

func TestParseHKSymbol(t *testing.T) {
	tests := []struct {
		input        string
		expectedCode string
		expectedOK   bool
	}{
		{"00700", "00700", true},
		{"00700.HK", "00700", true},
		{"0700.HK", "00700", true},
		{"700.HK", "00700", true},
		{"00941.HKEX", "00941", true},
		{"0941", "00941", true},
		{"9988", "09988", true},
		{"", "", false},
		{"ABC", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			code, ok := ParseHKSymbol(tt.input)
			if ok != tt.expectedOK {
				t.Errorf("ParseHKSymbol(%s) ok = %v, want %v", tt.input, ok, tt.expectedOK)
			}
			if ok && code != tt.expectedCode {
				t.Errorf("ParseHKSymbol(%s) code = %s, want %s", tt.input, code, tt.expectedCode)
			}
		})
	}
}

func TestToHKSecid(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		{"00700.HK", "116.00700", false},
		{"0700", "116.00700", false},
		{"700", "116.00700", false},
		{"00941", "116.00941", false},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := toHKSecid(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("toHKSecid(%s) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("toHKSecid(%s) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("toHKSecid(%s) = %s, want %s", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestToPeriod(t *testing.T) {
	tests := []struct {
		timeframe string
		expected  string
	}{
		{"1m", Period1m},
		{"5m", Period5m},
		{"15m", Period15m},
		{"30m", Period30m},
		{"60m", Period60m},
		{"1d", Period1d},
		{"1w", Period1w},
		{"1M", Period1M},
		{"", Period1d},
		{"invalid", Period1d},
	}

	for _, tt := range tests {
		t.Run(tt.timeframe, func(t *testing.T) {
			result := ToPeriod(tt.timeframe)
			if result != tt.expected {
				t.Errorf("ToPeriod(%s) = %s, want %s", tt.timeframe, result, tt.expected)
			}
		})
	}
}

func TestToAdjust(t *testing.T) {
	tests := []struct {
		adjust   string
		expected string
	}{
		{"qfq", "1"},
		{"hfq", "2"},
		{"", "0"},
		{"other", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.adjust, func(t *testing.T) {
			result := ToAdjust(tt.adjust)
			if result != tt.expected {
				t.Errorf("ToAdjust(%s) = %s, want %s", tt.adjust, result, tt.expected)
			}
		})
	}
}

func TestClient_GetKline(t *testing.T) {
	client := NewClient()
	defer client.Close()

	tests := []struct {
		name   string
		params *KlineParams
	}{
		{
			"00700 daily",
			&KlineParams{Symbol: "00700.HK", Period: Period1d, StartDate: "20240101", EndDate: "20240131"},
		},
		{
			"00941 weekly",
			&KlineParams{Symbol: "00941", Period: Period1w, StartDate: "20240101", EndDate: "20240630"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, record, err := client.GetKline(ctx, tt.params)
			if err != nil {
				checkAPIError(t, err)
				return
			}

			if record == nil {
				t.Fatal("record is nil")
			}

			t.Logf("Symbol: %s, Count: %d", result.Symbol, result.Count)

			if len(result.Data) == 0 {
				t.Log("Warning: no kline data returned")
				return
			}

			kline := result.Data[0]
			t.Logf("First bar: Date=%s, Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.0f",
				kline.Date, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume)

			if kline.Open <= 0 {
				t.Errorf("Expected Open > 0, got %f", kline.Open)
			}
			if kline.Close <= 0 {
				t.Errorf("Expected Close > 0, got %f", kline.Close)
			}
		})
	}
}

func TestPadHKCode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"700", "00700"},
		{"00700", "00700"},
		{"0700", "00700"},
		{"941", "00941"},
		{"9988", "09988"},
		{"", "00000"},
		{"0", "00000"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := padHKCode(tt.input)
			if result != tt.expected {
				t.Errorf("padHKCode(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
