package cninfo

import (
	"testing"

	"github.com/souloss/quantds/clients/cninfo"
	"github.com/souloss/quantds/domain"
)

func TestNewInstrumentAdapter(t *testing.T) {
	client := cninfo.NewClient(nil)
	adapter := NewInstrumentAdapter(client)

	if adapter == nil {
		t.Error("NewInstrumentAdapter returned nil")
	}

	if adapter.Name() != "cninfo" {
		t.Errorf("Expected name 'cninfo', got '%s'", adapter.Name())
	}
}

func TestInstrumentAdapter_Name(t *testing.T) {
	client := cninfo.NewClient(nil)
	adapter := NewInstrumentAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestInstrumentAdapter_SupportedMarkets(t *testing.T) {
	client := cninfo.NewClient(nil)
	adapter := NewInstrumentAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 || markets[0] != domain.MarketCN {
		t.Errorf("Expected supported markets [%s], got %v", domain.MarketCN, markets)
	}
}
