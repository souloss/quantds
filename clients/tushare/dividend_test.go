package tushare

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetDividend tests retrieving dividend data
// API Rule: Requires Tushare token
func TestClient_GetDividend(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetDividend(ctx, &DividendParams{
		TSCode: "000001.SZ", // Ping An Bank
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetDividend() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d dividend rows for 000001.SZ", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: ann_date=%s, ex_date=%s, cash_div=%.4f, stk_div=%.4f",
			r.AnnDate, r.ExDate, r.CashDiv, r.StkDiv)
		t.Logf("       stk_bo_rate=%.4f, stk_co_rate=%.4f, record_date=%s",
			r.StkBoRate, r.StkCoRate, r.RecordDate)
	}
}

// TestClient_GetDividend_ByYear tests dividend data for a specific year
func TestClient_GetDividend_ByYear(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Get dividend for 2023
	rows, _, err := client.GetDividend(ctx, &DividendParams{
		TSCode:  "600519.SH", // Kweichow Moutai
		AnnDate: "20230101",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetDividend() error = %v", err)
	}

	t.Logf("Got %d dividend rows for Moutai", len(rows))
	for _, r := range rows {
		t.Logf("  ann_date=%s, cash_div=%.4f, ex_date=%s",
			r.AnnDate, r.CashDiv, r.ExDate)
	}
}
