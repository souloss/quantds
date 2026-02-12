package eastmoney

import (
	"testing"
	"time"

	"github.com/souloss/quantds/clients/eastmoney"
	"github.com/souloss/quantds/domain"
)

func TestNewFinancialAdapter(t *testing.T) {
	client := eastmoney.NewClient(nil)
	adapter := NewFinancialAdapter(client)

	if adapter == nil {
		t.Error("NewFinancialAdapter returned nil")
	}

	if adapter.Name() != "eastmoney" {
		t.Errorf("Expected name 'eastmoney', got '%s'", adapter.Name())
	}
}

func TestFinancialAdapter_Name(t *testing.T) {
	client := eastmoney.NewClient(nil)
	adapter := NewFinancialAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestFinancialAdapter_SupportedMarkets(t *testing.T) {
	client := eastmoney.NewClient(nil)
	adapter := NewFinancialAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCN {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCN, markets)
	}
}

func TestFinancialAdapter_CanHandle(t *testing.T) {
	client := eastmoney.NewClient(nil)
	adapter := NewFinancialAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"000001.SZ", true},
		{"600001.SH", true},
		{"430001.BJ", true},
		{"00700.HK", false},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			result := adapter.CanHandle(tt.symbol)
			if result != tt.canHandle {
				t.Errorf("CanHandle(%s) = %v, want %v", tt.symbol, result, tt.canHandle)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	data := map[string]interface{}{
		"str_field":   "hello",
		"float_field": float64(42),
		"nil_field":   nil,
	}

	tests := []struct {
		key    string
		expect string
	}{
		{"str_field", "hello"},
		{"float_field", "*"},
		{"nil_field", ""},
		{"missing_field", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := getString(data, tt.key)
			if tt.expect != "*" && result != tt.expect {
				t.Errorf("getString(%s) = %q, want %q", tt.key, result, tt.expect)
			}
		})
	}
}

func TestGetFloat(t *testing.T) {
	data := map[string]interface{}{
		"float_field":  float64(3.14),
		"str_field":    "2.718",
		"nil_field":    nil,
		"invalid_str":  "abc",
	}

	tests := []struct {
		key    string
		expect float64
	}{
		{"float_field", 3.14},
		{"str_field", 2.718},
		{"nil_field", 0},
		{"invalid_str", 0},
		{"missing_field", 0},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := getFloat(data, tt.key)
			if result != tt.expect {
				t.Errorf("getFloat(%s) = %f, want %f", tt.key, result, tt.expect)
			}
		})
	}
}

func TestParseReportDate(t *testing.T) {
	tests := []struct {
		input  string
		isZero bool
	}{
		{"2024-03-31 00:00:00", false},
		{"", true},
		{"invalid-date", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseReportDate(tt.input)
			if result.IsZero() != tt.isZero {
				t.Errorf("parseReportDate(%q).IsZero() = %v, want %v", tt.input, result.IsZero(), tt.isZero)
			}
			if !tt.isZero && tt.input == "2024-03-31 00:00:00" {
				expected := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
				if !result.Equal(expected) {
					t.Errorf("parseReportDate(%q) = %v, want %v", tt.input, result, expected)
				}
			}
		})
	}
}
