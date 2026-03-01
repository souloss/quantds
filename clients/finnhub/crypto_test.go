package finnhub

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetCryptoCandles(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	now := time.Now()
	from := now.AddDate(0, -1, 0).Unix()
	to := now.Unix()

	result, _, err := client.GetCryptoCandles(context.Background(), &CandleParams{
		Symbol:     "BINANCE:BTCUSDT",
		Resolution: ResD,
		From:       from,
		To:         to,
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("BTC/USDT crypto candles: %d bars", result.Count)

	if len(result.Candles) == 0 {
		t.Log("Warning: no crypto candle data returned")
		return
	}

	candle := result.Candles[0]
	t.Logf("First: Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.0f",
		candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
}

func TestClient_GetCryptoSymbols(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetCryptoSymbols(context.Background(), "binance")
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Crypto symbols (binance): %d", result.Count)

	if len(result.Symbols) == 0 {
		t.Fatal("Expected crypto symbols, got 0")
	}

	t.Logf("First: Symbol=%s, Description=%s", result.Symbols[0].Symbol, result.Symbols[0].Description)
}
