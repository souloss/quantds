package eastmoney

import (
	"testing"

	"github.com/souloss/quantds/clients/eastmoney"
	"github.com/souloss/quantds/domain"
)

func TestNewKlineAdapter(t *testing.T) {
	client := eastmoney.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}

	if adapter.Name() != "eastmoney" {
		t.Errorf("Expected name 'eastmoney', got '%s'", adapter.Name())
	}
}

func TestKlineAdapter_Name(t *testing.T) {
	client := eastmoney.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_SupportedMarkets(t *testing.T) {
	client := eastmoney.NewClient()
	adapter := NewKlineAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCN {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCN, markets)
	}
}

func TestKlineAdapter_CanHandle(t *testing.T) {
	client := eastmoney.NewClient()
	adapter := NewKlineAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"000001.SZ", true},
		{"600001.SH", true},
		{"00700.HK", false},
		{"AAPL.US", false},
		{"BTCUSDT", false},
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
