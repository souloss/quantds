package coingecko

import (
	"context"
	"testing"
)

// TestClient_GetSimplePrice tests retrieving current prices
// API Rule: Supports multiple coins and currencies in one request
// Note: Price data is updated frequently but not real-time stream
func TestClient_GetSimplePrice(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	params := &SimplePriceRequest{
		IDs:          []string{"bitcoin", "ethereum"},
		VsCurrencies: []string{"usd", "cny"},
	}

	result, _, err := client.GetSimplePrice(ctx, params)
	if err != nil {
		checkAPIError(t, err)
	}

	t.Logf("Simple Price Response: %+v", result)

	if len(result) == 0 {
		t.Fatal("Expected prices, got 0")
	}

	if _, ok := result["bitcoin"]; !ok {
		t.Error("bitcoin price not found")
	}
	if _, ok := result["ethereum"]; !ok {
		t.Error("ethereum price not found")
	}

	btc := result["bitcoin"]
	if btc["usd"] <= 0 {
		t.Errorf("Expected BTC/USD > 0, got %f", btc["usd"])
	}
}
