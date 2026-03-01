package finnhub

import (
	"context"
	"testing"
)

func TestClient_GetQuote(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetQuote(context.Background(), &QuoteParams{Symbol: "AAPL"})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("AAPL: Current=%.2f, Open=%.2f, High=%.2f, Low=%.2f, PrevClose=%.2f",
		result.Current, result.Open, result.High, result.Low, result.PreviousClose)

	if result.Current <= 0 {
		t.Errorf("Expected Current > 0, got %f", result.Current)
	}
}
