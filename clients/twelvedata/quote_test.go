package twelvedata

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

	t.Logf("AAPL: Close=%.2f, Open=%.2f, High=%.2f, Low=%.2f, Name=%s",
		result.Close, result.Open, result.High, result.Low, result.Name)

	if result.Close <= 0 {
		t.Errorf("Expected Close > 0, got %f", result.Close)
	}
}
