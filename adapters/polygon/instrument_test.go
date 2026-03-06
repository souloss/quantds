package polygon

import (
	"testing"

	"github.com/souloss/quantds/clients/polygon"
	"github.com/souloss/quantds/domain"
)

func TestNewInstrumentAdapter(t *testing.T) {
	adapter := NewInstrumentAdapter(polygon.NewClient())
	if adapter == nil {
		t.Error("NewInstrumentAdapter returned nil")
	}
	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestInstrumentAdapter_CanHandle(t *testing.T) {
	adapter := NewInstrumentAdapter(polygon.NewClient())
	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"AAPL.US", true},
		{"000001.SZ", false},
	}
	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			if got := adapter.CanHandle(tt.symbol); got != tt.canHandle {
				t.Errorf("CanHandle(%s) = %v, want %v", tt.symbol, got, tt.canHandle)
			}
		})
	}
}

func TestInstrumentAdapter_SupportedMarkets(t *testing.T) {
	adapter := NewInstrumentAdapter(polygon.NewClient())
	if len(adapter.SupportedMarkets()) != 1 || adapter.SupportedMarkets()[0] != domain.MarketUS {
		t.Errorf("Expected [%s], got %v", domain.MarketUS, adapter.SupportedMarkets())
	}
}
