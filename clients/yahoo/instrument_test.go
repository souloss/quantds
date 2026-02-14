package yahoo

import (
	"context"
	"testing"
)

// TestClient_GetInstruments tests searching for US stock instruments
// API Rule: No authentication required
// Geo-Restriction: Yahoo Finance may be blocked in some regions
func TestClient_GetInstruments(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &InstrumentParams{
		Query: "Apple",
		Limit: 10,
	}

	result, record, err := client.GetInstruments(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Instrument Search Response Status: %d", record.Response.StatusCode)
	t.Logf("Found %d instruments", result.Total)

	if result.Total == 0 {
		t.Fatal("Expected instruments, got 0")
	}

	for i, inst := range result.Instruments {
		t.Logf("Instrument[%d]: symbol=%s, name=%s, exchange=%s, type=%s",
			i, inst.Symbol, inst.Name, inst.Exchange, inst.AssetType)
	}

	// Check if Apple is in the results
	found := false
	for _, inst := range result.Instruments {
		if inst.Symbol == "AAPL" {
			found = true
			break
		}
	}
	if !found {
		t.Log("Warning: AAPL not found in search results for 'Apple'")
	}
}

// TestClient_GetAllUSStocks tests retrieving all major US stocks
// This returns a hardcoded list of major stocks, no API call required
func TestClient_GetAllUSStocks(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetAllUSStocks(ctx)
	if err != nil {
		t.Fatalf("GetAllUSStocks failed: %v", err)
	}

	t.Logf("Got %d major US stocks", result.Total)

	if result.Total == 0 {
		t.Fatal("Expected stocks, got 0")
	}

	// Verify some well-known stocks are in the list
	expectedStocks := []string{"AAPL", "MSFT", "GOOGL", "AMZN", "TSLA"}
	for _, expected := range expectedStocks {
		found := false
		for _, inst := range result.Instruments {
			if inst.Symbol == expected {
				found = true
				t.Logf("Found expected stock: %s (%s)", inst.Symbol, inst.Name)
				break
			}
		}
		if !found {
			t.Errorf("Expected stock %s not found in list", expected)
		}
	}

	// Print sample of stocks
	t.Log("Sample of major US stocks:")
	for i := 0; i < 5 && i < len(result.Instruments); i++ {
		inst := result.Instruments[i]
		t.Logf("  %s - %s (%s)", inst.Symbol, inst.Name, inst.Exchange)
	}
}

// TestClient_GetInstrumentsByExchange tests filtering by exchange
func TestClient_GetInstrumentsByExchange(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	// Search with NASDAQ filter
	result, record, err := client.GetInstrumentsByExchange(ctx, "NASDAQ")
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Exchange Filter Response Status: %d", record.Response.StatusCode)
	t.Logf("Found %d NASDAQ instruments", result.Total)

	if result.Total > 0 {
		t.Logf("First result: %s - %s", result.Instruments[0].Symbol, result.Instruments[0].Name)
	}
}
