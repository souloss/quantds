package twelvedata

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

func TestToInterval(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1m", Interval1min},
		{"5m", Interval5min},
		{"1d", Interval1day},
		{"1w", Interval1week},
		{"1M", Interval1month},
		{"", Interval1day},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ToInterval(tt.input); got != tt.expected {
				t.Errorf("ToInterval(%s) = %s, want %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestClient_GetTimeSeries(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetTimeSeries(context.Background(), &TimeSeriesParams{
		Symbol:     "AAPL",
		Interval:   Interval1day,
		OutputSize: 10,
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("AAPL time series: %d data points", result.Count)

	if len(result.Data) == 0 {
		t.Fatal("Expected data, got 0")
	}

	d := result.Data[0]
	t.Logf("First: Datetime=%s, Open=%.2f, Close=%.2f, Volume=%.0f",
		d.Datetime, d.Open, d.Close, d.Volume)
}

func TestClient_GetTimeSeries_Forex(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetTimeSeries(context.Background(), &TimeSeriesParams{
		Symbol:     "EUR/USD",
		Interval:   Interval1day,
		OutputSize: 10,
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("EUR/USD time series: %d data points", result.Count)

	if len(result.Data) == 0 {
		t.Log("Warning: no forex data returned")
		return
	}

	t.Logf("First: Datetime=%s, Open=%.5f, Close=%.5f",
		result.Data[0].Datetime, result.Data[0].Open, result.Data[0].Close)
}
