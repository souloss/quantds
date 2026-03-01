package eodhd

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("NewClient returned nil")
	}
	defer client.Close()
}

func TestClient_GetEOD(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetEOD(context.Background(), &EODParams{
		Symbol: "AAPL.US",
		From:   "2024-01-01",
		To:     "2024-01-31",
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("AAPL.US EOD: %d data points", result.Count)

	if len(result.Data) == 0 {
		t.Fatal("Expected EOD data, got 0")
	}

	d := result.Data[0]
	t.Logf("First: Date=%s, Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.0f",
		d.Date, d.Open, d.High, d.Low, d.Close, d.Volume)
}
