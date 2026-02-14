package sina

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetSpot_Batch tests retrieving real-time quotes for multiple symbols
// API Rule: No authentication required, but has rate limiting
// Geo-Restriction: Sina Finance may be blocked outside China
func TestClient_GetSpot_Batch(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &SpotParams{
		Symbols: []string{"600519.SH", "000001.SZ", "000858.SZ"}, // Moutai, Ping An Bank, Wuliangye
	}

	result, record, err := client.GetSpot(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Spot Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d quotes", len(result.Data))

	if len(result.Data) == 0 {
		t.Fatal("Expected quotes, got 0")
	}

	for i, q := range result.Data {
		t.Logf("Quote[%d]: symbol=%s, name=%s, open=%.2f, latest=%.2f, high=%.2f, low=%.2f, volume=%.0f",
			i, q.Symbol, q.Name, q.Open, q.Latest, q.High, q.Low, q.Volume)

		if q.Symbol == "" {
			t.Error("Expected non-empty symbol")
		}
		if q.Name == "" {
			t.Logf("Warning: Empty name for symbol %s (market may be closed)", q.Symbol)
		}
	}
}

// TestClient_GetSpot_Single tests retrieving quote for a single symbol
func TestClient_GetSpot_Single(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &SpotParams{
		Symbols: []string{"600519.SH"}, // Kweichow Moutai
	}

	result, record, err := client.GetSpot(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Single Quote Response Status: %d", record.Response.StatusCode)

	if len(result.Data) == 0 {
		t.Log("Warning: No quote data returned (market may be closed)")
		return
	}

	quote := result.Data[0]
	t.Logf("Moutai (600519.SH): name=%s, latest=%.2f, open=%.2f, high=%.2f, low=%.2f",
		quote.Name, quote.Latest, quote.Open, quote.High, quote.Low)

	if quote.Latest <= 0 {
		t.Logf("Warning: Latest price is 0 or negative (market may be closed)")
	}
}

// TestClient_GetSpot_SHMarket tests Shanghai market stocks
func TestClient_GetSpot_SHMarket(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &SpotParams{
		Symbols: []string{"600036.SH", "601318.SH"}, // China Merchants Bank, Ping An Insurance
	}

	result, _, err := client.GetSpot(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Shanghai stocks: %d quotes", len(result.Data))
	for _, q := range result.Data {
		t.Logf("  %s (%s): %.2f", q.Symbol, q.Name, q.Latest)
	}
}

// TestClient_GetSpot_SZMarket tests Shenzhen market stocks
func TestClient_GetSpot_SZMarket(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &SpotParams{
		Symbols: []string{"000333.SZ", "002594.SZ"}, // Midea, BYD
	}

	result, _, err := client.GetSpot(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Shenzhen stocks: %d quotes", len(result.Data))
	for _, q := range result.Data {
		t.Logf("  %s (%s): %.2f", q.Symbol, q.Name, q.Latest)
	}
}
