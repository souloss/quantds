package sina

import (
	"testing"

	"github.com/souloss/quantds/clients/sina"
)

func TestNewKlineAdapter(t *testing.T) {
	client := sina.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter == nil {
		t.Error("NewKlineAdapter returned nil")
	}

	if adapter.Name() != "sina" {
		t.Errorf("Expected name 'sina', got '%s'", adapter.Name())
	}
}

func TestKlineAdapter_Name(t *testing.T) {
	client := sina.NewClient()
	adapter := NewKlineAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}
