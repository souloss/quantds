package eastmoneyhk

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetInstruments tests retrieving HK stock instruments
// API Rule: No authentication required
// Geo-Restriction: EastMoney HK API may be blocked in some regions
func TestClient_GetInstruments(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &InstrumentParams{
		PageSize:   20,
		PageNumber: 0,
	}

	result, record, err := client.GetInstruments(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Instrument Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d HK instruments (total: %d)", len(result.Instruments), result.Total)

	if len(result.Instruments) == 0 {
		t.Fatal("Expected instruments, got 0")
	}

	for i, inst := range result.Instruments {
		t.Logf("Instrument[%d]: code=%s, symbol=%s, name=%s, latest=%.2f, change=%.2f%%",
			i, inst.Code, inst.Symbol, inst.Name, inst.LatestPrice, inst.ChangeRate)
	}
}

// TestClient_GetInstrumentsByCode tests retrieving instruments by specific codes
func TestClient_GetInstrumentsByCode(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	codes := []string{"00700.HK", "00941.HK"} // Tencent, China Mobile
	result, _, err := client.GetInstrumentsByCode(ctx, codes)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d instruments by code", len(result.Instruments))
	for _, inst := range result.Instruments {
		t.Logf("  %s (%s): %.2f HKD", inst.Code, inst.Name, inst.LatestPrice)
	}
}

// TestClient_GetAllHKStocks tests retrieving all HK stocks
func TestClient_GetAllHKStocks(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetAllHKStocks(ctx)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d total HK stocks", result.Total)

	// Show sample of stocks
	for i := 0; i < 5 && i < len(result.Instruments); i++ {
		inst := result.Instruments[i]
		t.Logf("Sample[%d]: %s (%s) - %.2f HKD", i, inst.Code, inst.Name, inst.LatestPrice)
	}
}
