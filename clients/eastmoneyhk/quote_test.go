package eastmoneyhk

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetQuote tests retrieving HK stock quotes
// API Rule: No authentication required
// Geo-Restriction: EastMoney HK API may be blocked in some regions
func TestClient_GetQuote(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &QuoteParams{
		PageSize: 10,
		PageNo:   0,
	}

	result, record, err := client.GetQuote(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Quote Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d HK stocks (total: %d)", len(result.Quotes), result.Total)

	if len(result.Quotes) == 0 {
		t.Fatal("Expected quotes, got 0")
	}

	for i, q := range result.Quotes {
		t.Logf("Quote[%d]: code=%s, name=%s, latest=%.2f, change=%.2f (%.2f%%)",
			i, q.Code, q.Name, q.Latest, q.Change, q.ChangeRate)
	}
}

// TestClient_GetQuotesBySymbols tests retrieving quotes for specific symbols
func TestClient_GetQuotesBySymbols(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	symbols := []string{"00700.HK", "00941.HK", "09988.HK"} // Tencent, China Mobile, Alibaba
	result, _, err := client.GetQuotesBySymbols(ctx, symbols)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d quotes for specific symbols", len(result.Quotes))
	for _, q := range result.Quotes {
		t.Logf("  %s (%s): %.2f HKD (%.2f%%)", q.Code, q.Name, q.Latest, q.ChangeRate)
	}
}

// TestClient_GetSpot tests the Spot alias
func TestClient_GetSpot(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetSpot(ctx, &SpotParams{PageSize: 5})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d HK stocks", len(result.Quotes))
	for _, q := range result.Quotes {
		t.Logf("  %s: %.2f", q.Code, q.Latest)
	}
}
