package finnhub

import (
	"testing"

	"github.com/souloss/quantds/clients/finnhub"
	"github.com/souloss/quantds/domain"
)

func TestNewKlineAdapter(t *testing.T) {
	client := finnhub.NewClient()
	adapter := NewKlineAdapter(client)
	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}
	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestKlineAdapter_SupportedMarkets(t *testing.T) {
	adapter := NewKlineAdapter(finnhub.NewClient())
	markets := adapter.SupportedMarkets()
	found := map[domain.Market]bool{}
	for _, m := range markets {
		found[m] = true
	}
	for _, want := range []domain.Market{domain.MarketUS, domain.MarketForex, domain.MarketCrypto} {
		if !found[want] {
			t.Errorf("Expected market %s in supported markets", want)
		}
	}
}

func TestKlineAdapter_CanHandle(t *testing.T) {
	adapter := NewKlineAdapter(finnhub.NewClient())
	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"AAPL.US", true},
		{"MSFT.US.NASDAQ", true},
		{"000001.SZ", false},
		{"00700.HK.HKEX", false},
	}
	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			if got := adapter.CanHandle(tt.symbol); got != tt.canHandle {
				t.Errorf("CanHandle(%s) = %v, want %v", tt.symbol, got, tt.canHandle)
			}
		})
	}
}
