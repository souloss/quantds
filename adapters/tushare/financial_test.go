package tushare

import (
	"testing"

	"github.com/souloss/quantds/clients/tushare"
	"github.com/souloss/quantds/domain"
)

func TestNewFinancialAdapter(t *testing.T) {
	client := tushare.NewClient()
	adapter := NewFinancialAdapter(client)

	if adapter == nil {
		t.Error("NewFinancialAdapter returned nil")
	}

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestFinancialAdapter_SupportedMarkets(t *testing.T) {
	client := tushare.NewClient()
	adapter := NewFinancialAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCN {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCN, markets)
	}
}

func TestFinancialAdapter_CanHandle(t *testing.T) {
	client := tushare.NewClient()
	adapter := NewFinancialAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		// CN symbols should match for CN adapters
		{"000001.SZ", true},
		{"600519.SH", true},
		// Non-CN symbols should NOT match
		{"BTCUSDT", false},
		{"AAPL.US", false},
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
