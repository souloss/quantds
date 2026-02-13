package okx

import (
	"context"
	"testing"
)

// TestClient_GetCandlesticks tests retrieving historical K-line data
// API Rule: Max 1440 data points per request
// Supported bars: 1m, 3m, 5m, 15m, 30m, 1H, 2H, 4H...
func TestClient_GetCandlesticks(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	params := &CandlestickRequest{
		InstID: "BTC-USDT",
		Bar:    "1D",
		Limit:  5,
	}

	candles, _, err := client.GetCandlesticks(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Candlesticks count: %d", len(candles))
	if len(candles) > 0 {
		t.Logf("First Candle: %v", candles[0])
	}

	if len(candles) == 0 {
		t.Error("Expected candles, got 0")
	}

	for _, c := range candles {
		if len(c) < 5 {
			t.Errorf("Expected candle data length >= 5, got %d", len(c))
		}
	}
}
