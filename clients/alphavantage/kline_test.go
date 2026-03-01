package alphavantage

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

func TestClient_GetDailyTimeSeries(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetDailyTimeSeries(context.Background(), &KlineParams{
		Symbol: "AAPL",
		Size:   "compact",
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("AAPL daily: %d data points", result.Count)

	if len(result.Data) == 0 {
		t.Fatal("Expected kline data, got 0")
	}

	d := result.Data[0]
	t.Logf("First: Date=%s, Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.0f",
		d.Date, d.Open, d.High, d.Low, d.Close, d.Volume)

	if d.Open <= 0 {
		t.Errorf("Expected Open > 0, got %f", d.Open)
	}
}
