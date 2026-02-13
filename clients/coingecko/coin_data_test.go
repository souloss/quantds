package coingecko

import (
	"context"
	"testing"
)

// TestClient_GetCoinData tests retrieving detailed coin information
// API Rule: Rate limit 10-30 calls/min
// Note: Returns comprehensive data including description, links, market data
func TestClient_GetCoinData(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	params := &CoinDataRequest{
		ID:            "bitcoin",
		Localization:  false,
		Tickers:       false,
		MarketData:    true,
		CommunityData: true,
		DeveloperData: true,
		Sparkline:     false,
	}

	result, _, err := client.GetCoinData(ctx, params)
	if err != nil {
		t.Fatalf("GetCoinData failed: %v", err)
	}

	t.Logf("Coin Name: %s", result.Name)
	t.Logf("Symbol: %s", result.Symbol)
	
	if result.ID != "bitcoin" {
		t.Errorf("Expected ID bitcoin, got %s", result.ID)
	}

	if result.MarketData == nil {
		t.Error("Expected MarketData, got nil")
	} else {
		price := result.MarketData.CurrentPrice["usd"]
		t.Logf("Current Price (USD): %f", price)
		if price <= 0 {
			t.Errorf("Expected price > 0, got %f", price)
		}
	}

	if result.CommunityData != nil {
		t.Logf("Twitter Followers: %d", result.CommunityData.TwitterFollowers)
	}
}
