package coingecko

import (
	"context"
	"testing"
)

// TestClient_Ping verifies connectivity to CoinGecko API
// API Rule: No authentication required
// Rate Limit: 10-30 calls/min for public API
func TestClient_Ping(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	result, err := client.Ping(ctx)
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}

	t.Logf("Ping response: %s", result.GeckoSays)

	if result.GeckoSays != "(V3) To the Moon!" {
		t.Errorf("Unexpected ping response: %v", result)
	}
}
