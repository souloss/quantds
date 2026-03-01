package eodhd

import (
	"context"
	"testing"
)

func TestClient_GetRealTimeQuote(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetRealTimeQuote(context.Background(), &RealTimeParams{Symbol: "AAPL.US"})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("AAPL: Close=%.2f, Open=%.2f, High=%.2f, Low=%.2f, Volume=%.0f",
		result.Close, result.Open, result.High, result.Low, result.Volume)
}
