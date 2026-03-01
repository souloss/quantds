package tushare

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetIndexBasic(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetIndexBasic(ctx, &IndexBasicParams{
		Market: "SSE",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetIndexBasic() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d index_basic rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: code=%s, name=%s, publisher=%s, market=%s",
			r.TSCode, r.Name, r.Publisher, r.Market)
	}
}

func TestClient_GetIndexDaily(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetIndexDaily(ctx, &IndexDailyParams{
		TSCode:    "000001.SH",
		StartDate: "20240101",
		EndDate:   "20240131",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetIndexDaily() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d index_daily rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: code=%s, date=%s, open=%.2f, close=%.2f, vol=%.2f",
			r.TSCode, r.TradeDate, r.Open, r.Close, r.Vol)
	}
}
