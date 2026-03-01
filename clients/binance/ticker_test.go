package binance

import (
	"context"
	"testing"
)

func TestClient_GetTicker24hr(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &TickerParams{
		Symbol: "BTCUSDT",
	}

	result, _, err := client.GetTicker24hr(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	if len(result.Tickers) != 1 {
		t.Fatalf("Expected 1 ticker, got %d", len(result.Tickers))
	}

	ticker := result.Tickers[0]
	t.Logf("BTCUSDT LastPrice: %f, Volume: %f", ticker.LastPrice, ticker.Volume)

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
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &TickerParams{
		Symbol: "",
	}

	result, _, err := client.GetTicker24hr(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	if len(result.Tickers) < 2 {
		t.Fatalf("Expected at least 2 tickers, got %d", len(result.Tickers))
	}

	t.Logf("Total tickers returned: %d", len(result.Tickers))

	foundBTC := false
	foundETH := false
	for _, tk := range result.Tickers {
		if tk.Symbol == "BTCUSDT" {
			foundBTC = true
		}
		if tk.Symbol == "ETHUSDT" {
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

func TestClient_GetPrice(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &PriceParams{
		Symbol: "BTCUSDT",
	}

	result, _, err := client.GetPrice(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	if len(result.Prices) != 1 {
		t.Fatalf("Expected 1 price, got %d", len(result.Prices))
	}

	price := result.Prices[0]
	t.Logf("BTCUSDT Price: %f", price.Price)

	if price.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", price.Symbol)
	}
	if price.Price <= 0 {
		t.Errorf("Expected Price > 0, got %f", price.Price)
	}
}

func TestClient_GetPrice_All(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &PriceParams{
		Symbol: "",
	}

	result, _, err := client.GetPrice(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	if len(result.Prices) < 2 {
		t.Fatalf("Expected at least 2 prices, got %d", len(result.Prices))
	}

	t.Logf("Total prices returned: %d", result.Count)

	foundBTC := false
	for _, p := range result.Prices {
		if p.Symbol == "BTCUSDT" {
			foundBTC = true
			t.Logf("BTCUSDT Price: %f", p.Price)
			if p.Price <= 0 {
				t.Errorf("Expected BTCUSDT price > 0, got %f", p.Price)
			}
			break
		}
	}

	if !foundBTC {
		t.Error("BTCUSDT not found in all prices response")
	}
}

func TestClient_GetSpot(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &SpotParams{
		Symbol: "ETHUSDT",
	}

	result, _, err := client.GetSpot(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	if len(result.Tickers) != 1 {
		t.Fatalf("Expected 1 spot ticker, got %d", len(result.Tickers))
	}

	spot := result.Tickers[0]
	t.Logf("ETHUSDT Spot: LastPrice=%f, Volume=%f, HighPrice=%f, LowPrice=%f",
		spot.LastPrice, spot.Volume, spot.HighPrice, spot.LowPrice)

	if spot.Symbol != "ETHUSDT" {
		t.Errorf("Expected symbol ETHUSDT, got %s", spot.Symbol)
	}
	if spot.LastPrice <= 0 {
		t.Errorf("Expected LastPrice > 0, got %f", spot.LastPrice)
	}
	if spot.HighPrice <= 0 {
		t.Errorf("Expected HighPrice > 0, got %f", spot.HighPrice)
	}
	if spot.LowPrice <= 0 {
		t.Errorf("Expected LowPrice > 0, got %f", spot.LowPrice)
	}
}
