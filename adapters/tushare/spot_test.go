package tushare

import (
	"context"
	"strings"
	"testing"

	"github.com/souloss/quantds/clients/tushare"
	"github.com/souloss/quantds/domain/spot"
)

func TestSpotAdapter_Fetch(t *testing.T) {
	// mock client logic is hard without network, so we use real network if token is present,
	// or skip if no token.
	// But unit tests should ideally mock the client response.
	// Given the complexity of mocking http client here without a mock framework,
	// we will rely on integration testing style if token exists, or just basic structure check.
	
	client := tushare.NewClient() // Will use env vars
	adapter := NewSpotAdapter(client)

	t.Run("CanHandle", func(t *testing.T) {
		if !adapter.CanHandle("000001.SZ") {
			t.Error("Should handle SZ stock")
		}
		if !adapter.CanHandle("600000.SH") {
			t.Error("Should handle SH stock")
		}
		if adapter.CanHandle("AAPL.US") {
			t.Error("Should not handle US stock")
		}
	})

	t.Run("Fetch", func(t *testing.T) {
		// Only run if TUSHARE_TOKEN is set
		if client.Token() == "" {
			t.Skip("TUSHARE_TOKEN not set, skipping integration test")
		}

		ctx := context.Background()
		req := spot.Request{
			Symbols: []string{"000001.SZ", "600000.SH"},
		}

		resp, trace, err := adapter.Fetch(ctx, nil, req)
		if err != nil {
			msg := err.Error()
			if strings.Contains(msg, "token") || strings.Contains(msg, "40101") || strings.Contains(msg, "-1") {
				t.Skipf("Token issue, skipping: %v", err)
			}
			t.Fatalf("Fetch failed: %v", err)
		}

		if trace == nil {
			t.Error("Trace should not be nil")
		}

		if len(resp.Quotes) == 0 {
			t.Log("No quotes returned (maybe market closed or no permission)")
		} else {
			for _, q := range resp.Quotes {
				t.Logf("Quote: %+v", q)
				if q.Symbol == "" {
					t.Error("Symbol should not be empty")
				}
				if q.Latest == 0 {
					t.Log("Latest price is 0 (suspension?)")
				}
			}
		}
	})
}
