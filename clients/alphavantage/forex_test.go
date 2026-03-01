package alphavantage

import (
	"context"
	"testing"
)

func TestClient_GetForexExchangeRate(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetForexExchangeRate(context.Background(), &ForexRateParams{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("USD/EUR: Rate=%.5f, Bid=%.5f, Ask=%.5f, LastRefreshed=%s",
		result.ExchangeRate, result.BidPrice, result.AskPrice, result.LastRefreshed)

	if result.ExchangeRate <= 0 {
		t.Errorf("Expected ExchangeRate > 0, got %f", result.ExchangeRate)
	}
}
