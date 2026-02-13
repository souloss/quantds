package sse

import (
	"context"
	"testing"
	"time"

	"github.com/souloss/quantds/request"
)

func TestClient_GetInstruments(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	tests := []struct {
		name   string
		params *InstrumentParams
	}{
		{"all stocks", nil},
		{"with page size", &InstrumentParams{PageSize: "100"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, record, err := client.GetInstruments(ctx, tt.params)
			if err != nil {
				t.Skipf("GetInstruments() error = %v (API may be unavailable)", err)
				return
			}

			if record == nil {
				t.Fatal("record is nil")
			}

			if result == nil {
				t.Fatal("result is nil")
			}

			t.Logf("Got %d stocks (total: %d)", len(result.Data), result.Total)

			if len(result.Data) == 0 {
				t.Skip("no stocks returned (API may have no data)")
				return
			}

			for i, s := range result.Data[:min(3, len(result.Data))] {
				t.Logf("Stock[%d]: code=%s, name=%s", i, s.CompanyCode, s.CompanyAbbr)
			}

			// Test ToInstrumentData conversion
			if len(result.Data) > 0 {
				inst := result.Data[0].ToInstrumentData()
				if inst.Code == "" {
					t.Error("ToInstrumentData() returned empty code")
				}
				if inst.Name == "" {
					t.Error("ToInstrumentData() returned empty name")
				}
				t.Logf("Converted instrument: code=%s, name=%s", inst.Code, inst.Name)
			}
		})
	}
}

func TestStockRow_ToInstrumentData(t *testing.T) {
	row := StockRow{
		CompanyCode: "600000",
		CompanyAbbr: "浦发银行",
		SecNameCn:   "上海浦东发展银行股份有限公司",
		ListDate:    "1999-11-10",
		TotalShares: "2935200",
		FloatShares: "2935200",
		Industry:    "银行业",
	}

	inst := row.ToInstrumentData()

	if inst.Code != "600000" {
		t.Errorf("Code = %v, want %v", inst.Code, "600000")
	}
	if inst.Name != "浦发银行" {
		t.Errorf("Name = %v, want %v", inst.Name, "浦发银行")
	}
	if inst.FullName != "上海浦东发展银行股份有限公司" {
		t.Errorf("FullName = %v, want %v", inst.FullName, "上海浦东发展银行股份有限公司")
	}
	if inst.ListDate != "1999-11-10" {
		t.Errorf("ListDate = %v, want %v", inst.ListDate, "1999-11-10")
	}
	if inst.Industry != "银行业" {
		t.Errorf("Industry = %v, want %v", inst.Industry, "银行业")
	}
}
