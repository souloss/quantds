package okx

import (
	"context"
	"testing"
)

// TestClient_GetInstruments tests retrieving list of instruments
// API Rule: No authentication required
// Geo-Restriction: OKX API may be blocked in some regions
func TestClient_GetInstruments(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &InstrumentParams{
		InstType: "SPOT",
	}

	result, record, err := client.GetInstruments(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Instruments Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d spot instruments", result.Total)

	if result.Total == 0 {
		t.Fatal("Expected instruments, got 0")
	}

	for i, inst := range result.Instruments {
		if i >= 5 {
			break
		}
		t.Logf("Instrument[%d]: id=%s, base=%s, quote=%s, state=%s",
			i, inst.InstID, inst.BaseCcy, inst.QuoteCcy, inst.State)
	}
}

// TestClient_GetSpotInstruments tests retrieving spot instruments
func TestClient_GetSpotInstruments(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetSpotInstruments(ctx)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d spot instruments", result.Total)

	// Find BTC-USDT
	for _, inst := range result.Instruments {
		if inst.InstID == "BTC-USDT" {
			t.Logf("Found BTC-USDT: base=%s, quote=%s, minSz=%s, lotSz=%s",
				inst.BaseCcy, inst.QuoteCcy, inst.MinSz, inst.LotSz)
			break
		}
	}
}

// TestClient_GetSwapInstruments tests retrieving swap instruments
func TestClient_GetSwapInstruments(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetSwapInstruments(ctx)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d swap (perpetual) instruments", result.Total)

	for i, inst := range result.Instruments {
		if i >= 3 {
			break
		}
		t.Logf("Swap[%d]: %s (%s-%s)", i, inst.InstID, inst.BaseCcy, inst.SettleCcy)
	}
}

// TestClient_GetFuturesInstruments tests retrieving futures instruments
func TestClient_GetFuturesInstruments(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetFuturesInstruments(ctx)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d futures instruments", result.Total)

	for i, inst := range result.Instruments {
		if i >= 3 {
			break
		}
		t.Logf("Futures[%d]: %s (expiry=%s)", i, inst.InstID, inst.ExpTime)
	}
}
