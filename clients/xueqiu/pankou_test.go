package xueqiu

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetPankou tests retrieving order book (pankou) data
// Note: Xueqiu API may require authentication
func TestClient_GetPankou(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, record, err := client.GetPankou(ctx, "000001.SZ")
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Pankou Response Status: %d", record.Response.StatusCode)
	t.Logf("Symbol: %s", result.Symbol)
	t.Logf("Buy 1: price=%.2f, volume=%.0f", result.Bp1, result.Bc1)
	t.Logf("Buy 2: price=%.2f, volume=%.0f", result.Bp2, result.Bc2)
	t.Logf("Buy 3: price=%.2f, volume=%.0f", result.Bp3, result.Bc3)
	t.Logf("Sell 1: price=%.2f, volume=%.0f", result.Sp1, result.Sc1)
	t.Logf("Sell 2: price=%.2f, volume=%.0f", result.Sp2, result.Sc2)
	t.Logf("Sell 3: price=%.2f, volume=%.0f", result.Sp3, result.Sc3)

	if result.Bp1 == 0 && result.Sp1 == 0 {
		t.Log("Warning: Order book prices are zero (market may be closed or auth required)")
	}
}

// TestClient_GetPankou_SH tests pankou for Shanghai stocks
func TestClient_GetPankou_SH(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetPankou(ctx, "600519.SH") // Moutai
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Moutai Pankou:")
	t.Logf("  Bid: %.2f (%.0f), %.2f (%.0f), %.2f (%.0f)",
		result.Bp1, result.Bc1, result.Bp2, result.Bc2, result.Bp3, result.Bc3)
	t.Logf("  Ask: %.2f (%.0f), %.2f (%.0f), %.2f (%.0f)",
		result.Sp1, result.Sc1, result.Sp2, result.Sc2, result.Sp3, result.Sc3)
}
