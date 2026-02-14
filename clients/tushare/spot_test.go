package tushare

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetSpot tests retrieving spot data (based on daily API)
// API Rule: Requires Tushare token
func TestClient_GetSpot(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetSpot(ctx, &SpotParams{
		TsCode: "000001.SZ",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetSpot() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d spot rows for 000001.SZ", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: date=%s, open=%.2f, high=%.2f, low=%.2f, close=%.2f",
			r.TradeDate, r.Open, r.High, r.Low, r.Close)
		t.Logf("       volume=%.0f, amount=%.2f, change_pct=%.2f%%",
			r.Vol, r.Amount, r.PctChg)
	}
}

// TestClient_GetSpot_Multiple tests retrieving spot data for multiple stocks
func TestClient_GetSpot_Multiple(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := context.Background()

	codes := []string{"000001.SZ", "600519.SH", "000858.SZ"}
	for _, code := range codes {
		rows, _, err := client.GetSpot(ctx, &SpotParams{
			TsCode: code,
		})
		skipOnTokenError(t, err)
		if err != nil {
			t.Errorf("GetSpot(%s) error = %v", code, err)
			continue
		}

		if len(rows) > 0 {
			r := rows[0]
			t.Logf("%s: close=%.2f, change_pct=%.2f%%", code, r.Close, r.PctChg)
		}
	}
}
