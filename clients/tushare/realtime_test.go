package tushare

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetRealtimeQuote tests retrieving realtime tick quotes
// API Rule: Requires Tushare token (0 credits - fully open)
// Note: Data has delay, not real-time
func TestClient_GetRealtimeQuote(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetRealtimeQuote(ctx, &RealtimeQuoteParams{
		TSCode:   "000001.SZ",
		PageSize: 5,
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetRealtimeQuote() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d realtime quote rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: code=%s, name=%s, price=%.2f, change=%.2f, change_pct=%.2f%%",
			r.TSCode, r.Name, r.Trade, r.PriceChange, r.ChangePercent)
		t.Logf("       open=%.2f, high=%.2f, low=%.2f, pre_close=%.2f",
			r.Open, r.High, r.Low, r.PreClose)
		t.Logf("       volume=%.0f, amount=%.2f, tick_time=%s",
			r.Volume, r.Amount, r.TickTime)
	}
}

// TestClient_GetRealtimeQuote_Wildcard tests retrieving quotes with wildcard
func TestClient_GetRealtimeQuote_Wildcard(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := context.Background()

	rows, _, err := client.GetRealtimeQuote(ctx, &RealtimeQuoteParams{
		TSCode:   "6*.SH", // Shanghai main board
		PageSize: 5,
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetRealtimeQuote() error = %v", err)
	}

	t.Logf("Got %d Shanghai main board quotes", len(rows))
	for i, r := range rows {
		if i >= 3 {
			break
		}
		t.Logf("  %s (%s): %.2f (%.2f%%)", r.TSCode, r.Name, r.Trade, r.ChangePercent)
	}
}

// TestClient_GetRtK tests retrieving realtime daily K-line data
func TestClient_GetRtK(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetRtK(ctx, &RtKParams{
		TSCode: "000001.SZ",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetRtK() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d realtime K-line rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: code=%s, date=%s, open=%.2f, high=%.2f, low=%.2f, close=%.2f",
			r.TSCode, r.TradeDate, r.Open, r.High, r.Low, r.Close)
		t.Logf("       pre_close=%.2f, change=%.2f, pct_chg=%.2f%%, vol=%.0f, amount=%.2f",
			r.PreClose, r.Change, r.PctChg, r.Vol, r.Amount)
	}
}

// TestClient_GetRtK_Wildcard tests retrieving K-lines with wildcard
func TestClient_GetRtK_Wildcard(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := context.Background()

	rows, _, err := client.GetRtK(ctx, &RtKParams{
		TSCode: "0*.SZ", // Shenzhen main board
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetRtK() error = %v", err)
	}

	t.Logf("Got %d Shenzhen main board K-lines", len(rows))
	for i, r := range rows {
		if i >= 3 {
			break
		}
		t.Logf("  %s: %.2f (%.2f%%)", r.TSCode, r.Close, r.PctChg)
	}
}
