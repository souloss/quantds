package coingecko

import (
	"context"
	"testing"
)

// TestClient_GetCoinsList tests retrieving all supported coins
// API Rule: No pagination required for this endpoint
// Note: This endpoint returns a large list of coins (10k+ items)
func TestClient_GetCoinsList(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	params := &CoinsListRequest{
		IncludePlatform: true,
	}

	result, _, err := client.GetCoinsList(ctx, params)
	if err != nil {
		t.Fatalf("GetCoinsList failed: %v", err)
	}

	t.Logf("Coins count: %d", len(result))
	
	if len(result) == 0 {
		t.Fatal("Expected coins list, got 0")
	}

	// Verify Bitcoin exists in the list
	found := false
	for _, coin := range result {
		if coin.ID == "bitcoin" {
			found = true
			t.Logf("Found Bitcoin: %+v", coin)
			break
		}
	}

	if !found {
		t.Error("Bitcoin not found in coins list")
	}
}
