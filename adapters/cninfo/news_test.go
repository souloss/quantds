package cninfo

import (
	"testing"

	"github.com/souloss/quantds/clients/cninfo"
)

func TestNewAnnouncementAdapter(t *testing.T) {
	client := cninfo.NewClient()
	adapter := NewAnnouncementAdapter(client)

	if adapter == nil {
		t.Error("NewAnnouncementAdapter returned nil")
	}

	if adapter.Name() != "cninfo" {
		t.Errorf("Expected name 'cninfo', got '%s'", adapter.Name())
	}
}

func TestAnnouncementAdapter_Name(t *testing.T) {
	client := cninfo.NewClient()
	adapter := NewAnnouncementAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestAnnouncementAdapter_CanHandle(t *testing.T) {
	client := cninfo.NewClient()
	adapter := NewAnnouncementAdapter(client)

	tests := []struct {
		symbol    string
		canHandle bool
	}{
		{"", true},
		{"000001.SZ", true},
		{"600001.SH", true},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			result := adapter.CanHandle(tt.symbol)
			if result != tt.canHandle {
				t.Errorf("CanHandle(%s) = %v, want %v", tt.symbol, result, tt.canHandle)
			}
		})
	}
}
