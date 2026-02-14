package eastmoney

import (
	"testing"

	"github.com/souloss/quantds/clients/eastmoney"
	"github.com/souloss/quantds/domain"
)

func TestNewAnnouncementAdapter(t *testing.T) {
	client := eastmoney.NewClient()
	adapter := NewAnnouncementAdapter(client)

	if adapter == nil {
		t.Error("NewAnnouncementAdapter returned nil")
	}

	if adapter.Name() != "eastmoney" {
		t.Errorf("Expected name 'eastmoney', got '%s'", adapter.Name())
	}
}

func TestAnnouncementAdapter_Name(t *testing.T) {
	client := eastmoney.NewClient()
	adapter := NewAnnouncementAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestAnnouncementAdapter_SupportedMarkets(t *testing.T) {
	client := eastmoney.NewClient()
	adapter := NewAnnouncementAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCN {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCN, markets)
	}
}

func TestAnnouncementAdapter_CanHandle(t *testing.T) {
	client := eastmoney.NewClient()
	adapter := NewAnnouncementAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"000001.SZ", true},
		{"600001.SH", true},
		{"430001.BJ", true},
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
