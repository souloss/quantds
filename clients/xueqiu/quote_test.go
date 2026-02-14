package xueqiu

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetQuoteDetail tests retrieving quote detail
// Note: Xueqiu API may require authentication
func TestClient_GetQuoteDetail(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, record, err := client.GetQuoteDetail(ctx, &QuoteDetailParams{
		Symbol: "000001.SZ",
		Extend: true,
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Quote Detail Response Status: %d", record.Response.StatusCode)
	t.Logf("Symbol: %s", result.Symbol)
	t.Logf("Name: %s", result.Name)
	t.Logf("Latest: %.2f", result.Current)
	t.Logf("Change: %.2f (%.2f%%)", result.Change, result.Percent)
	t.Logf("Market Capitalization: %.2f", result.TotalMarketCap)
}

// TestClient_GetProfile tests retrieving company profile
func TestClient_GetProfile(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetProfile(ctx, &QuoteDetailParams{
		Symbol: "600519.SH",
		Extend: true,
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Moutai Profile:")
	t.Logf("  Symbol: %s", result.Symbol)
	t.Logf("  Name: %s", result.Name)
	t.Logf("  Industry: %s", result.Industry)
	t.Logf("  Market Cap: %.2f", result.TotalMarketCap)
}
