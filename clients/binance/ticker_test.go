package binance

import (
	"context"
	"testing"
)

func TestClient_GetTicker24hr(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	params := &TickerParams{
		Symbol: "BTCUSDT",
	}

	result, _, err := client.GetTicker24hr(ctx, params)
	if err != nil {
		checkAPIError(t, err)
	}

	if len(result.Tickers) != 1 {
		t.Fatalf("Expected 1 ticker, got %d", len(result.Tickers))
	}

	ticker := result.Tickers[0]
	if ticker.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", ticker.Symbol)
	}
	if ticker.LastPrice <= 0 {
		t.Errorf("Expected LastPrice > 0, got %f", ticker.LastPrice)
	}
	if ticker.Volume <= 0 {
		t.Errorf("Expected Volume > 0, got %f", ticker.Volume)
	}
}

func TestClient_GetTicker24hr_All(t *testing.T) {
	client := NewClient(nil)
	ctx := context.Background()

	params := &TickerParams{
		Symbol: "", // Empty symbol means all tickers
	}

	result, _, err := client.GetTicker24hr(ctx, params)
	if err != nil {
		checkAPIError(t, err)
	}

	if len(result.Tickers) < 2 {
		t.Fatalf("Expected at least 2 tickers, got %d", len(result.Tickers))
	}
	
	// Check if common pairs are present
	foundBTC := false
	foundETH := false
	for _, t := range result.Tickers {
		if t.Symbol == "BTCUSDT" {
			foundBTC = true
		}
		if t.Symbol == "ETHUSDT" {
			foundETH = true
		}
	}
	
	if !foundBTC {
		t.Error("BTCUSDT ticker not found in all tickers response")
	}
	if !foundETH {
		t.Error("ETHUSDT ticker not found in all tickers response")
	}
}
