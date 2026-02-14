package coingecko

import (
	"context"
	"testing"
)

// TestClient_Search tests coin search functionality
// API Rule: Searches across coins, exchanges, and categories
// Note: Useful for finding coin IDs (e.g. 'bitcoin') from symbols ('BTC')
func TestClient_Search(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	params := &SearchRequest{
		Query: "bitcoin",
	}

	result, _, err := client.Search(ctx, params)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	t.Logf("Search Coins count: %d", len(result.Coins))
	if len(result.Coins) > 0 {
		t.Logf("First Coin: %+v", result.Coins[0])
	}

	if len(result.Coins) == 0 {
		t.Error("Expected search results, got 0")
	}

	found := false
	for _, coin := range result.Coins {
		if coin.ID == "bitcoin" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Bitcoin not found in search results")
	}
}
