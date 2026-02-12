package tushare

import (
	"testing"

	"github.com/souloss/quantds/clients/tushare"
)

func TestNewKlineAdapter(t *testing.T) {
	client := tushare.NewClient(nil)
	adapter := NewKlineAdapter(client)

	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}

	if adapter.Name() != "tushare" {
		t.Errorf("Expected name 'tushare', got '%s'", adapter.Name())
	}
}

func TestKlineAdapter_Name(t *testing.T) {
	client := tushare.NewClient(nil)
	adapter := NewKlineAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}
