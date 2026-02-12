package sina

import (
	"context"
	"testing"
	"time"

	"github.com/souloss/quantds/request"
)

func TestClient_GetKline(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	tests := []struct {
		name   string
		params *KlineParams
	}{
		{"daily SH", &KlineParams{Symbol: "600001.SH", Period: "d"}},
		{"daily SZ", &KlineParams{Symbol: "000001.SZ", Period: "d"}},
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

			t.Logf("Got %d bars for %s", result.Count, tt.params.Symbol)

			if result.Count > 0 {
				t.Logf("First: date=%s, open=%.2f, close=%.2f", result.Data[0].Date, result.Data[0].Open, result.Data[0].Close)
			}
		})
	}
}

func TestClient_GetSpot(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, record, err := client.GetSpot(ctx, &SpotParams{
		Symbols: []string{"000001.SZ", "600001.SH"},
	})
	if err != nil {
		t.Fatalf("GetSpot() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d quotes", len(result.Data))

	for i, q := range result.Data {
		t.Logf("Quote[%d]: symbol=%s, name=%s, latest=%.2f", i, q.Symbol, q.Name, q.Latest)
	}
}

func TestToSinaSymbol(t *testing.T) {
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
			got, err := toSinaSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("toSinaSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toSinaSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPeriod(t *testing.T) {
	tests := []struct {
		timeframe string
		want      string
	}{
		{"5m", "5"},
		{"15m", "15"},
		{"30m", "30"},
		{"60m", "60"},
		{"1d", "d"},
		{"", "d"},
		{"1w", "w"},
		{"1M", "m"},
	}

	for _, tt := range tests {
		t.Run(tt.timeframe, func(t *testing.T) {
			if got := ToPeriod(tt.timeframe); got != tt.want {
				t.Errorf("ToPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}
