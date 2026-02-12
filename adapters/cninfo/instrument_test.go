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

func TestInferExchangeFromCode(t *testing.T) {
	tests := []struct {
		code string
		want string
	}{
		{"600001", "SH"},
		{"000001", "SZ"},
		{"300001", "SZ"},
		{"430001", "BJ"},
		{"830001", "BJ"},
		{"680001", "SH"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got := inferExchangeFromCode(tt.code)
			if string(got) != tt.want {
				t.Errorf("inferExchangeFromCode(%s) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}
