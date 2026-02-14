package bse

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetInstruments(t *testing.T) {
	client := NewClient()
	defer client.Close()

	tests := []struct {
		name   string
		params *InstrumentParams
	}{
		{"all stocks (nil params)", nil},
		{"all stocks (empty params)", &InstrumentParams{}},
		{"page 1", &InstrumentParams{Page: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			rows, records, err := client.GetInstruments(ctx, tt.params)
			if err != nil {
				t.Skipf("GetInstruments() error = %v (API may be unavailable)", err)
				return
			}

			if len(records) == 0 {
				t.Fatal("no records returned")
			}

			t.Logf("Got %d stocks in %d requests", len(rows), len(records))

			if len(rows) == 0 {
				t.Skip("no stocks returned (API may have no data)")
				return
			}

			for i, s := range rows[:min(3, len(rows))] {
				t.Logf("Stock[%d]: code=%s, name=%s", i, s.StockCode, s.StockName)
			}

			// Test ToInstrumentData conversion
			if len(rows) > 0 {
				inst := rows[0].ToInstrumentData()
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

func TestClient_GetInstrumentsPage(t *testing.T) {
	client := NewClient()
	defer client.Close()

	tests := []struct {
		name   string
		params *InstrumentParams
	}{
		{"page 1", &InstrumentParams{Page: 1}},
		{"page 1 with type", &InstrumentParams{Page: 1, Typejb: "T", Xxfcbj: "2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, record, err := client.GetInstrumentsPage(ctx, tt.params)
			if err != nil {
				t.Skipf("GetInstrumentsPage() error = %v (API may be unavailable)", err)
				return
			}

			if record == nil {
				t.Fatal("record is nil")
			}

			if result == nil {
				t.Fatal("result is nil")
			}

			t.Logf("Got %d stocks, total pages: %d", len(result.Data), result.TotalPages)

			if len(result.Data) == 0 {
				t.Skip("no stocks returned (API may have no data)")
				return
			}

			for i, s := range result.Data[:min(3, len(result.Data))] {
				t.Logf("Stock[%d]: code=%s, name=%s", i, s.StockCode, s.StockName)
			}
		})
	}
}

func TestStockRow_ToInstrumentData(t *testing.T) {
	row := StockRow{
		StockCode:   "835305",
		StockName:   "云创数据",
		ListDate:    "2021-08-09",
		TotalShares: 13237600,
		FloatShares: 5000000,
		Industry:    "软件和信息技术服务业",
	}

	inst := row.ToInstrumentData()

	if inst.Code != "835305" {
		t.Errorf("Code = %v, want %v", inst.Code, "835305")
	}
	if inst.Name != "云创数据" {
		t.Errorf("Name = %v, want %v", inst.Name, "云创数据")
	}
	if inst.ListDate != "2021-08-09" {
		t.Errorf("ListDate = %v, want %v", inst.ListDate, "2021-08-09")
	}
	if inst.TotalShares != 13237600 {
		t.Errorf("TotalShares = %v, want %v", inst.TotalShares, 13237600)
	}
	if inst.FloatShares != 5000000 {
		t.Errorf("FloatShares = %v, want %v", inst.FloatShares, 5000000)
	}
	if inst.Industry != "软件和信息技术服务业" {
		t.Errorf("Industry = %v, want %v", inst.Industry, "软件和信息技术服务业")
	}
}
