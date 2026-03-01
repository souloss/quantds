package finnhub

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

func TestToResolution(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1m", Res1},
		{"5m", Res5},
		{"15m", Res15},
		{"30m", Res30},
		{"60m", Res60},
		{"1h", Res60},
		{"1d", ResD},
		{"1w", ResW},
		{"1M", ResM},
		{"", ResD},
		{"invalid", ResD},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToResolution(tt.input)
			if result != tt.expected {
				t.Errorf("ToResolution(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestClient_GetStockCandles(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	now := time.Now()
	from := now.AddDate(0, -1, 0).Unix()
	to := now.Unix()

	result, _, err := client.GetStockCandles(context.Background(), &CandleParams{
		Symbol:     "AAPL",
		Resolution: ResD,
		From:       from,
		To:         to,
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("AAPL stock candles: %d bars", result.Count)

	if len(result.Candles) == 0 {
		t.Fatal("Expected candle data, got 0")
	}

	candle := result.Candles[0]
	t.Logf("First: Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.0f",
		candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)

	if candle.Open <= 0 {
		t.Errorf("Expected Open > 0, got %f", candle.Open)
	}
}
