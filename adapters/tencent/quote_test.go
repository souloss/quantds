package tencent

import (
	"testing"

	"github.com/souloss/quantds/clients/tencent"
)

func TestNewQuoteAdapter(t *testing.T) {
	client := tencent.NewClient(nil)
	adapter := NewQuoteAdapter(client)

	if adapter == nil {
		t.Error("NewQuoteAdapter returned nil")
	}

	if adapter.Name() != "tencent" {
		t.Errorf("Expected name 'tencent', got '%s'", adapter.Name())
	}
}

func TestQuoteAdapter_Name(t *testing.T) {
	client := tencent.NewClient(nil)
	adapter := NewQuoteAdapter(client)

	if adapter.Name() != Name {
		t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
	}
}

func TestQuoteAdapter_SupportedMarkets(t *testing.T) {
	client := tencent.NewClient(nil)
	adapter := NewQuoteAdapter(client)

	markets := adapter.SupportedMarkets()
	if len(markets) != 1 {
		t.Errorf("Expected 1 supported market, got %d", len(markets))
	}

	if markets[0] != "CN" {
		t.Errorf("Expected market 'CN', got '%s'", markets[0])
	}
}

func TestQuoteAdapter_CanHandle(t *testing.T) {
	client := tencent.NewClient(nil)
	adapter := NewQuoteAdapter(client)

	tests := []struct {
		symbol string
		want   bool
	}{
		{"000001.SZ", true},
		{"600001.SH", true},
		{"430001.BJ", true},
		{"AAPL.US", false},
		{"0700.HK", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			got := adapter.CanHandle(tt.symbol)
			if got != tt.want {
				t.Errorf("CanHandle(%s) = %v, want %v", tt.symbol, got, tt.want)
			}
		})
	}
}
