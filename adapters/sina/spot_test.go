package sina

import (
	"testing"

	"github.com/souloss/quantds/clients/sina"
)

func TestNewSpotAdapter(t *testing.T) {
	client := sina.NewClient(nil)
	adapter := NewSpotAdapter(client)

	if adapter == nil {
		t.Error("NewSpotAdapter returned nil")
	}

	if adapter.Name() != "sina" {
		t.Errorf("Expected name 'sina', got '%s'", adapter.Name())
	}
}

func TestSpotAdapter_Name(t *testing.T) {
	client := sina.NewClient(nil)
	adapter := NewSpotAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}
