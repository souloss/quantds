package tencent

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetQuotes tests retrieving real-time quotes
// API Rule: No authentication required
// Geo-Restriction: Tencent Finance may be blocked outside China
func TestClient_GetQuotes(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &QuoteParams{
		Symbols: []string{"600519.SH", "000001.SZ"}, // Moutai, Ping An Bank
	}

	result, record, err := client.GetQuotes(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Quote Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d quotes", result.Count)

	if result.Count == 0 {
		t.Fatal("Expected quotes, got 0")
	}

	for i, q := range result.Data {
		t.Logf("Quote[%d]: symbol=%s, name=%s, latest=%.2f, open=%.2f, high=%.2f, low=%.2f",
			i, q.Symbol, q.Name, q.Latest, q.Open, q.High, q.Low)
		t.Logf("         change=%.2f, changeRate=%.2f%%, volume=%.0f, time=%s",
			q.Change, q.ChangeRate, q.Volume, q.Time)
	}
}

// TestClient_GetQuotes_Single tests retrieving quote for a single symbol
func TestClient_GetQuotes_Single(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &QuoteParams{
		Symbols: []string{"600519.SH"}, // Kweichow Moutai
	}

	result, _, err := client.GetQuotes(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	if result.Count != 1 {
		t.Fatalf("Expected 1 quote, got %d", result.Count)
	}

	quote := result.Data[0]
	t.Logf("Moutai (600519): name=%s, latest=%.2f, change=%.2f (%.2f%%)",
		quote.Name, quote.Latest, quote.Change, quote.ChangeRate)
}

// TestClient_GetQuotes_SHMarket tests Shanghai market stocks
func TestClient_GetQuotes_SHMarket(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &QuoteParams{
		Symbols: []string{"600036.SH", "601318.SH"}, // China Merchants Bank, Ping An Insurance
	}

	result, _, err := client.GetQuotes(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Shanghai stocks: %d quotes", result.Count)
	for _, q := range result.Data {
		t.Logf("  %s (%s): %.2f (%.2f%%)", q.Symbol, q.Name, q.Latest, q.ChangeRate)
	}
}

// TestClient_GetQuotes_SZMarket tests Shenzhen market stocks
func TestClient_GetQuotes_SZMarket(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &QuoteParams{
		Symbols: []string{"000333.SZ", "002594.SZ"}, // Midea, BYD
	}

	result, _, err := client.GetQuotes(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Shenzhen stocks: %d quotes", result.Count)
	for _, q := range result.Data {
		t.Logf("  %s (%s): %.2f (%.2f%%)", q.Symbol, q.Name, q.Latest, q.ChangeRate)
	}
}
