package binance

import (
	"context"
	"testing"
)

func TestClient_GetExchangeInfo(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetExchangeInfo(ctx)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	if result.Total == 0 {
		t.Error("Expected instruments, got 0")
	}

	if len(result.Instruments) == 0 {
		t.Fatal("Expected instruments, got 0")
	}

	t.Logf("Total instruments: %d", result.Total)

	found := false
	for _, inst := range result.Instruments {
		if inst.Symbol == "BTCUSDT" {
			found = true
			if inst.BaseAsset != "BTC" {
				t.Errorf("Expected base asset BTC, got %s", inst.BaseAsset)
			}
			if inst.QuoteAsset != "USDT" {
				t.Errorf("Expected quote asset USDT, got %s", inst.QuoteAsset)
			}
			break
		}
	}

	if !found {
		t.Error("BTCUSDT instrument not found")
	}
}

func TestClient_GetInstruments_Filter(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &InstrumentParams{
		Quote: "USDT",
	}

	result, _, err := client.GetInstruments(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	if len(result.Instruments) == 0 {
		t.Error("Expected instruments with quote USDT, got 0")
	}

	t.Logf("USDT pairs count: %d", len(result.Instruments))

	for _, inst := range result.Instruments {
		if inst.QuoteAsset != "USDT" {
			t.Errorf("Expected quote asset USDT, got %s for %s", inst.QuoteAsset, inst.Symbol)
		}
	}
}

func TestClient_GetAllSpotInstruments(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetAllSpotInstruments(ctx)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	if len(result.Instruments) == 0 {
		t.Fatal("Expected spot instruments, got 0")
	}

	t.Logf("Total spot instruments (TRADING): %d", result.Total)

	for _, inst := range result.Instruments {
		if inst.Status != "TRADING" {
			t.Errorf("Expected status TRADING, got %s for %s", inst.Status, inst.Symbol)
		}
	}
}

func TestClient_GetInstrumentsByQuote(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	tests := []struct {
		quote       string
		minExpected int
	}{
		{"USDT", 10},
		{"BTC", 5},
		{"ETH", 2},
	}

	for _, tt := range tests {
		t.Run(tt.quote, func(t *testing.T) {
			result, _, err := client.GetInstrumentsByQuote(ctx, tt.quote)
			if err != nil {
				checkAPIError(t, err)
				return
			}

			t.Logf("%s pairs (TRADING): %d", tt.quote, len(result.Instruments))

			if len(result.Instruments) < tt.minExpected {
				t.Errorf("Expected at least %d %s pairs, got %d", tt.minExpected, tt.quote, len(result.Instruments))
			}

			for _, inst := range result.Instruments {
				if inst.QuoteAsset != tt.quote {
					t.Errorf("Expected quote %s, got %s for %s", tt.quote, inst.QuoteAsset, inst.Symbol)
				}
			}
		})
	}
}
