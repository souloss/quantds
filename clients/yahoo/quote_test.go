package yahoo

import (
	"context"
	"testing"
)

// TestClient_GetQuote tests retrieving real-time quotes for US stocks
// API Rule: No authentication required, but may have rate limiting
// Geo-Restriction: Yahoo Finance may be blocked in some regions (e.g., China mainland)
func TestClient_GetQuote(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &QuoteParams{
		Symbols: []string{"AAPL", "MSFT", "GOOGL"},
	}

	result, record, err := client.GetQuote(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Quote Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d quotes", result.Count)

	if result.Count == 0 {
		t.Fatal("Expected quotes, got 0")
	}

	for i, q := range result.Quotes {
		t.Logf("Quote[%d]: symbol=%s, name=%s, latest=%.2f, change=%.2f, changeRate=%.2f%%",
			i, q.Symbol, q.Name, q.Latest, q.Change, q.ChangeRate)

		if q.Symbol == "" {
			t.Error("Expected non-empty symbol")
		}
		if q.Latest <= 0 {
			t.Errorf("Expected Latest > 0 for %s, got %f", q.Symbol, q.Latest)
		}
	}
}

// TestClient_GetQuote_Single tests retrieving quote for a single symbol
func TestClient_GetQuote_Single(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &QuoteParams{
		Symbols: []string{"AAPL"},
	}

	result, record, err := client.GetQuote(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Quote Response Status: %d", record.Response.StatusCode)

	if result.Count != 1 {
		t.Fatalf("Expected 1 quote, got %d", result.Count)
	}

	quote := result.Quotes[0]
	t.Logf("Apple (AAPL): Latest=%.2f, Open=%.2f, High=%.2f, Low=%.2f, Volume=%.0f",
		quote.Latest, quote.Open, quote.High, quote.Low, quote.Volume)

	if quote.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got %s", quote.Symbol)
	}
	if quote.Latest <= 0 {
		t.Errorf("Expected Latest > 0, got %f", quote.Latest)
	}
}

// TestClient_GetSpot tests the Spot alias for GetQuote
func TestClient_GetSpot(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &SpotParams{
		Symbols: []string{"TSLA"},
	}

	result, _, err := client.GetSpot(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Tesla (TSLA): Latest=%.2f", result.Quotes[0].Latest)

	if result.Count != 1 {
		t.Errorf("Expected 1 quote, got %d", result.Count)
	}
}
