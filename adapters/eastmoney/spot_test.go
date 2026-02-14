package eastmoney

import (
	"testing"

	"github.com/souloss/quantds/clients/eastmoney"
)

func TestNewSpotAdapter(t *testing.T) {
	client := eastmoney.NewClient()
	adapter := NewSpotAdapter(client)

	if adapter == nil {
		t.Error("NewSpotAdapter returned nil")
	}

	if adapter.Name() != "eastmoney" {
		t.Errorf("Expected name 'eastmoney', got '%s'", adapter.Name())
	}
}

func TestSpotAdapter_Name(t *testing.T) {
	client := eastmoney.NewClient()
	adapter := NewSpotAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}
