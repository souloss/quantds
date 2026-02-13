package szse

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
		{"default params", nil},
		{"with catalog", &InstrumentParams{CatalogID: "1110", TabKey: "tab1"}},
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

			t.Logf("Got %d rows of data", len(result.Data))

			if len(result.Data) == 0 {
				t.Skip("no data returned (API may have no data)")
				return
			}

			// Log header row
			if len(result.Data) > 0 {
				t.Logf("Header: %v", result.Data[0])
			}

			// Log first few data rows
			for i, row := range result.Data[1:min(4, len(result.Data))] {
				t.Logf("Row[%d]: %v", i, row)
			}
		})
	}
}

func TestParseInstruments(t *testing.T) {
	// Test with sample data matching SZSE Excel format
	// Need at least 17 columns for Province(14), City(15), Industry(16) to be parsed
	// Columns: 板块, 公司代码, 公司简称, 公司全称, A股代码, A股简称, A股上市日期, B股代码, B股简称, B股上市日期, 地区, 省份, 城市, 行业, 17列之后的数据
	rows := [][]string{
		{"板块", "公司代码", "公司简称", "公司全称", "A股代码", "A股简称", "A股上市日期", "B股代码", "B股简称", "B股上市日期", "地区", "省份", "城市", "行业", "col14", "col15", "col16"},
		{"主板", "000001", "平安银行股份", "平安银行股份有限公司", "000001", "平安银行", "1991-04-03", "", "", "", "华南", "col11", "col12", "col13", "广东省", "深圳市", "货币金融服务"},
		{"主板", "000002", "万科企业股份", "万科企业股份有限公司", "000002", "万科A", "1991-01-29", "", "", "", "华南", "col11", "col12", "col13", "广东省", "深圳市", "房地产业"},
	}

	instruments := ParseInstruments(rows)

	if len(instruments) != 2 {
		t.Errorf("ParseInstruments() returned %d instruments, want 2", len(instruments))
	}

	if len(instruments) > 0 {
		inst := instruments[0]
		if inst.Code != "000001" {
			t.Errorf("First instrument Code = %v, want %v", inst.Code, "000001")
		}
		if inst.Name != "平安银行" {
			t.Errorf("First instrument Name = %v, want %v", inst.Name, "平安银行")
		}
		if inst.ListDate != "1991-04-03" {
			t.Errorf("First instrument ListDate = %v, want %v", inst.ListDate, "1991-04-03")
		}
		if inst.Province != "广东省" {
			t.Errorf("First instrument Province = %v, want %v", inst.Province, "广东省")
		}
		if inst.City != "深圳市" {
			t.Errorf("First instrument City = %v, want %v", inst.City, "深圳市")
		}
		if inst.Industry != "货币金融服务" {
			t.Errorf("First instrument Industry = %v, want %v", inst.Industry, "货币金融服务")
		}
	}
}

func TestParseInstruments_Empty(t *testing.T) {
	// Test with empty data
	instruments := ParseInstruments(nil)
	if instruments != nil {
		t.Error("ParseInstruments(nil) should return nil")
	}

	// Test with header only
	instruments = ParseInstruments([][]string{{"header1", "header2"}})
	if instruments != nil {
		t.Error("ParseInstruments with header only should return nil")
	}
}

func TestParseInstruments_Malformed(t *testing.T) {
	// Test with short rows
	rows := [][]string{
		{"header1", "header2"},
		{"short"}, // too short, should be skipped
		{"主板", "000001", "平安银行", "公司全称", "000001", "平安银行", "1991-04-03"},
	}

	instruments := ParseInstruments(rows)

	// Should have 1 valid instrument (short row skipped)
	if len(instruments) != 1 {
		t.Errorf("ParseInstruments() returned %d instruments, want 1", len(instruments))
	}
}
