package polygon

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

func TestToTimespan(t *testing.T) {
	tests := []struct {
		input            string
		expectedTimespan string
		expectedMult     int
	}{
		{"1m", TimespanMinute, 1},
		{"5m", TimespanMinute, 5},
		{"15m", TimespanMinute, 15},
		{"1d", TimespanDay, 1},
		{"1w", TimespanWeek, 1},
		{"1M", TimespanMonth, 1},
		{"", TimespanDay, 1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ts, mult := ToTimespan(tt.input)
			if ts != tt.expectedTimespan {
				t.Errorf("ToTimespan(%s) timespan = %s, want %s", tt.input, ts, tt.expectedTimespan)
			}
			if mult != tt.expectedMult {
				t.Errorf("ToTimespan(%s) multiplier = %d, want %d", tt.input, mult, tt.expectedMult)
			}
		})
	}
}

func TestClient_GetAggregates(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetAggregates(context.Background(), &AggregateParams{
		Symbol:     "AAPL",
		Multiplier: 1,
		Timespan:   TimespanDay,
		From:       "2024-01-01",
		To:         "2024-01-31",
		Limit:      50,
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("AAPL aggregates: %d bars", result.Count)

	if len(result.Bars) == 0 {
		t.Fatal("Expected aggregate bars, got 0")
	}

	bar := result.Bars[0]
	t.Logf("First: Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.0f",
		bar.Open, bar.High, bar.Low, bar.Close, bar.Volume)
}
