package tencent

import (
	"testing"

	"github.com/souloss/quantds/clients/tencent"
)

func TestNewKlineAdapter(t *testing.T) {
	client := tencent.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}

	if adapter.Name() != "tencent" {
		t.Errorf("Expected name 'tencent', got '%s'", adapter.Name())
	}
}

func TestKlineAdapter_Name(t *testing.T) {
	client := tencent.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}
