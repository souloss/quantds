package okx

import (
	"context"
	"testing"
)

// TestClient_GetTicker tests retrieving latest ticker info
// API Rule: Rate limit 20 req/2s
// Geo-Restriction: OKX API may be blocked in US, China Mainland, etc.
func TestClient_GetTicker(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	params := &TickerRequest{
		InstID: "BTC-USDT",
	}

	ticker, _, err := client.GetTicker(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Ticker Response: %+v", ticker)

	if ticker.InstID != "BTC-USDT" {
		t.Errorf("Expected InstID BTC-USDT, got %s", ticker.InstID)
	}
	if ticker.Last == "" {
		t.Error("Expected Last price, got empty")
	}
}
