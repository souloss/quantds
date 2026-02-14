package xueqiu

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetSpot tests retrieving real-time spot data
// Note: Xueqiu API may require authentication cookie for some endpoints
// Geo-Restriction: May be blocked outside China
func TestClient_GetSpot(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &SpotParams{
		Symbols: []string{"000001.SZ"}, // Ping An Bank
	}

	result, record, err := client.GetSpot(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Spot Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d quotes", result.Count)

	if result.Count == 0 {
		t.Log("Warning: No spot data returned (may require authentication)")
		return
	}

	for i, q := range result.Data {
		t.Logf("Quote[%d]: symbol=%s, latest=%.2f, change=%.2f (%.2f%%)",
			i, q.Symbol, q.Latest, q.Change, q.ChangeRate)
	}
}

// TestClient_GetSpot_Multiple tests retrieving spot for multiple symbols
func TestClient_GetSpot_Multiple(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &SpotParams{
		Symbols: []string{"600519.SH", "000001.SZ", "601318.SH"},
	}

	result, _, err := client.GetSpot(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d quotes", result.Count)
	for _, q := range result.Data {
		t.Logf("  %s: %.2f (%.2f%%)", q.Symbol, q.Latest, q.ChangeRate)
	}
}

// TestParseSpotResponse tests the spot response parsing
func TestParseSpotResponse(t *testing.T) {
	body := []byte(`{"data":[{"symbol":"SZ000001","current":12.5,"open":12.3,"high":12.8,"low":12.1,"last_close":12.2,"chg":0.3,"percent":2.46,"volume":100000,"amount":1250000,"time":1700000000000}]}`)

	result, err := parseSpotResponse(body)
	if err != nil {
		t.Fatalf("parseSpotResponse() error = %v", err)
	}

	if result.Count != 1 {
		t.Errorf("Count = %d, want 1", result.Count)
	}

	if result.Data[0].Symbol != "SZ000001" {
		t.Errorf("Symbol = %s, want SZ000001", result.Data[0].Symbol)
	}

	if result.Data[0].Latest != 12.5 {
		t.Errorf("Latest = %v, want 12.5", result.Data[0].Latest)
	}
}
