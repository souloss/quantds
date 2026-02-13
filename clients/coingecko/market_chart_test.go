package coingecko

import (
	"context"
	"testing"
)

// TestClient_GetMarketChart tests retrieving historical K-line data
// API Rule: Data granularity is automatic based on 'days' parameter:
// - 1 day: 5 minute interval
// - 1-90 days: hourly interval
// - >90 days: daily interval
func TestClient_GetMarketChart(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	params := &MarketChartRequest{
		ID:         "bitcoin",
		VsCurrency: "usd",
		Days:       "1",
	}

	result, _, err := client.GetMarketChart(ctx, params)
	if err != nil {
		t.Fatalf("GetMarketChart failed: %v", err)
	}

	t.Logf("Market Chart Prices count: %d", len(result.Prices))
	if len(result.Prices) > 0 {
		t.Logf("First Price Point: %v", result.Prices[0])
	}

	if len(result.Prices) == 0 {
		t.Error("Expected price data, got 0")
	}
}
