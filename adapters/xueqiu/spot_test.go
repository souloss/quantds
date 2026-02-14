package xueqiu

import (
	"testing"

	"github.com/souloss/quantds/clients/xueqiu"
	"github.com/souloss/quantds/domain"
)

func TestNewSpotAdapter(t *testing.T) {
	client := xueqiu.NewClient()
	adapter := NewSpotAdapter(client)

	if adapter == nil {
		t.Error("NewSpotAdapter returned nil")
	}

	if adapter.Name() != "xueqiu" {
		t.Errorf("Expected name 'xueqiu', got '%s'", adapter.Name())
	}
}

func TestSpotAdapter_Name(t *testing.T) {
	client := xueqiu.NewClient()
	adapter := NewSpotAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestSpotAdapter_SupportedMarkets(t *testing.T) {
	client := xueqiu.NewClient()
	adapter := NewSpotAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCN {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCN, markets)
	}
}

func TestSpotAdapter_CanHandle(t *testing.T) {
	client := xueqiu.NewClient()
	adapter := NewSpotAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"000001.SZ", true},
		{"600001.SH", true},
		{"00700.HK", false},
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
