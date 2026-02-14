package tencent

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetSpot(t *testing.T) {
	client := NewClient()
	defer client.Close()

	tests := []struct {
		name   string
		params *SpotParams
	}{
		{
			name: "single stock",
			params: &SpotParams{
				Symbols: []string{"000001.SZ"},
			},
		},
		{
			name: "multiple stocks",
			params: &SpotParams{
				Symbols: []string{"000001.SZ", "600001.SH"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, record, err := client.GetSpot(ctx, tt.params)
			if err != nil {
				t.Fatalf("GetSpot() error = %v", err)
			}

			if record == nil {
				t.Fatal("record is nil")
			}

			if result == nil {
				t.Fatal("result is nil")
			}

			t.Logf("Got %d quotes", result.Count)

			for i, q := range result.Data {
				t.Logf("Quote[%d]: symbol=%s, name=%s, latest=%.2f, change=%.2f, changeRate=%.2f%%",
					i, q.Symbol, q.Name, q.Latest, q.Change, q.ChangeRate)

				if q.Symbol == "" {
					t.Errorf("quote[%d].Symbol is empty", i)
				}
				if q.Name == "" {
					t.Errorf("quote[%d].Name is empty", i)
				}
			}
		})
	}
}

func TestClient_GetKline(t *testing.T) {
	client := NewClient()
	defer client.Close()

	tests := []struct {
		name   string
		params *KlineParams
	}{
		{
			name: "daily kline for SZ stock",
			params: &KlineParams{
				Symbol: "000001.SZ",
				Period: "day",
				Count:  30,
			},
		},
		{
			name: "weekly kline for SH stock",
			params: &KlineParams{
				Symbol: "600519.SH",
				Period: "week",
				Count:  20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, record, err := client.GetKline(ctx, tt.params)
			if err != nil {
				t.Fatalf("GetKline() error = %v", err)
			}

			if record == nil {
				t.Fatal("record is nil")
			}

			if result == nil {
				t.Fatal("result is nil")
			}

			t.Logf("Got %d bars for %s", result.Count, tt.params.Symbol)

			if result.Count > 0 {
				first := result.Data[0]
				t.Logf("First bar: date=%s, open=%.2f, close=%.2f, high=%.2f, low=%.2f, vol=%.0f",
					first.Date, first.Open, first.Close, first.High, first.Low, first.Volume)

				if first.Date == "" {
					t.Error("first bar date is empty")
				}
			}
		})
	}
}

func TestToTencentSymbol(t *testing.T) {
	tests := []struct {
		symbol  string
		want    string
		wantErr bool
	}{
		{"600001.SH", "sh600001", false},
		{"000001.SZ", "sz000001", false},
		{"430001.BJ", "bj430001", false},
		{"INVALID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			got, err := toTencentSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("toTencentSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toTencentSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPeriod(t *testing.T) {
	tests := []struct {
		timeframe string
		want      string
	}{
		{"1m", "m1"},
		{"5m", "m5"},
		{"15m", "m15"},
		{"30m", "m30"},
		{"60m", "m60"},
		{"1d", "day"},
		{"", "day"},
		{"1w", "week"},
		{"1M", "month"},
		{"unknown", "day"},
	}

	for _, tt := range tests {
		t.Run(tt.timeframe, func(t *testing.T) {
			if got := ToPeriod(tt.timeframe); got != tt.want {
				t.Errorf("ToPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}
