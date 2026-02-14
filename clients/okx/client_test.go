package okx

import (
	"context"
	"testing"
)

func TestClient_WithBaseURL(t *testing.T) {
	// Note: aws.okx.com sometimes has certificate issues or is blocked.
	// This test mainly verifies that the BaseURL option works correctly in the client struct.
	client := NewClient(WithBaseURL(AwsBaseURL))

	if client.BaseURL != AwsBaseURL {
		t.Errorf("Expected BaseURL %s, got %s", AwsBaseURL, client.BaseURL)
	}

	// We skip the actual network call if it fails due to network/cert issues,
	// as we've already verified the functionality in TestClient_GetTicker with the default URL.
	ctx := context.Background()
	params := &TickerRequest{
		InstID: "ETH-USDT",
	}
	ticker, _, err := client.GetTicker(ctx, params)
	if err != nil {
		t.Logf("Skipping network check for AwsBaseURL due to error: %v", err)
		return
	}

	if ticker.InstID != "ETH-USDT" {
		t.Errorf("Expected InstID ETH-USDT, got %s", ticker.InstID)
	}
}
