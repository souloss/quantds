package tushare

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetFundBasic(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetFundBasic(ctx, &FundBasicParams{
		Market: "E",
		Status: "L",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetFundBasic() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d fund_basic rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: code=%s, name=%s, type=%s, management=%s, market=%s",
			r.TSCode, r.Name, r.FundType, r.Management, r.Market)
	}
}

func TestClient_GetFundNAV(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetFundNAV(ctx, &FundNAVParams{
		TSCode: "159901.SZ",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetFundNAV() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d fund_nav rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: code=%s, end_date=%s, unit_nav=%.4f, accum_nav=%.4f",
			r.TSCode, r.EndDate, r.UnitNAV, r.AccumNAV)
	}
}
