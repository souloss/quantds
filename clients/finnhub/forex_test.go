package finnhub

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetForexRates(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetForexRates(context.Background(), &ForexRatesParams{Base: "USD"})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Base: %s, Rates count: %d", result.Base, len(result.Rates))

	if len(result.Rates) == 0 {
		t.Fatal("Expected rates, got 0")
	}

	if eur, ok := result.Rates["EUR"]; ok {
		t.Logf("USD/EUR: %f", eur)
		if eur <= 0 {
			t.Errorf("Expected EUR rate > 0, got %f", eur)
		}
	}
}

func TestClient_GetForexCandles(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	now := time.Now()
	from := now.AddDate(0, -1, 0).Unix()
	to := now.Unix()

	result, _, err := client.GetForexCandles(context.Background(), &CandleParams{
		Symbol:     "OANDA:EUR_USD",
		Resolution: ResD,
		From:       from,
		To:         to,
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("EUR/USD forex candles: %d bars", result.Count)

	if len(result.Candles) == 0 {
		t.Log("Warning: no forex candle data returned")
		return
	}

	candle := result.Candles[0]
	t.Logf("First: Open=%.5f, High=%.5f, Low=%.5f, Close=%.5f",
		candle.Open, candle.High, candle.Low, candle.Close)
}

func TestClient_GetForexSymbols(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetForexSymbols(context.Background(), "oanda")
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Forex symbols (oanda): %d", result.Count)

	if len(result.Symbols) == 0 {
		t.Fatal("Expected forex symbols, got 0")
	}

	t.Logf("First: Symbol=%s, Description=%s", result.Symbols[0].Symbol, result.Symbols[0].Description)
}
