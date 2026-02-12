package tencent

import (
	"testing"

	"github.com/souloss/quantds/clients/tencent"
)

func TestNewSpotAdapter(t *testing.T) {
	client := tencent.NewClient(nil)
	adapter := NewSpotAdapter(client)

	if adapter == nil {
		t.Error("NewSpotAdapter returned nil")
	}

	if adapter.Name() != "tencent" {
		t.Errorf("Expected name 'tencent', got '%s'", adapter.Name())
	}
}

func TestSpotAdapter_Name(t *testing.T) {
	client := tencent.NewClient(nil)
	adapter := NewSpotAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}
