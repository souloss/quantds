package binance

import (
	"context"
	"testing"
)

func TestClient_GetExchangeInfo(t *testing.T) {
	client := NewClient(nil)
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

	// Check if BTCUSDT exists
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
	client := NewClient(nil)
	ctx := context.Background()

	// Test filter by Quote Asset
	params := &InstrumentParams{
		Quote: "USDT",
	}

	result, _, err := client.GetInstruments(ctx, params)
	if err != nil {
		t.Fatalf("GetInstruments failed: %v", err)
	}

	if len(result.Instruments) == 0 {
		t.Error("Expected instruments with quote USDT, got 0")
	}

	for _, inst := range result.Instruments {
		if inst.QuoteAsset != "USDT" {
			t.Errorf("Expected quote asset USDT, got %s", inst.QuoteAsset)
		}
	}
}
