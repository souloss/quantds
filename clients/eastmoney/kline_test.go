package eastmoney

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetKline(t *testing.T) {
	client := NewClient()
	defer client.Close()

	tests := []struct {
		name   string
		params *KlineParams
	}{
		{
			name: "daily kline for SH stock",
			params: &KlineParams{
				Symbol:    "600001.SH",
				StartDate: "20240101",
				EndDate:   "20240131",
				Period:    "101",
				Adjust:    "0",
			},
		},
		{
			name: "daily kline for SZ stock",
			params: &KlineParams{
				Symbol:    "000001.SZ",
				StartDate: "20240101",
				EndDate:   "20240131",
				Period:    "101",
				Adjust:    "1",
			},
		},
		{
			name: "weekly kline",
			params: &KlineParams{
				Symbol:    "600519.SH",
				StartDate: "20230101",
				EndDate:   "20240131",
				Period:    "102",
				Adjust:    "0",
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

			if !record.IsSuccess() {
				t.Errorf("record should be success, got error: %v", record.Error)
			}

			if result == nil {
				t.Fatal("result is nil")
			}

			if len(result.Data) == 0 {
				t.Log("Warning: no data returned (might be holiday or future date)")
				return
			}

			for i, bar := range result.Data {
				if bar.Date == "" {
					t.Errorf("bar[%d].Date is empty", i)
				}
				if bar.Open <= 0 {
					t.Errorf("bar[%d].Open = %v, should be positive", i, bar.Open)
				}
				if bar.High < bar.Low {
					t.Errorf("bar[%d].High = %v < Low = %v", i, bar.High, bar.Low)
				}
				if bar.Volume < 0 {
					t.Errorf("bar[%d].Volume = %v, should not be negative", i, bar.Volume)
				}
			}

			t.Logf("Got %d bars from %s to %s",
				len(result.Data),
				result.Data[0].Date,
				result.Data[len(result.Data)-1].Date)
		})
	}
}

func TestClient_GetKline_InvalidSymbol(t *testing.T) {
	client := NewClient()
	defer client.Close()

	ctx := context.Background()
	_, _, err := client.GetKline(ctx, &KlineParams{
		Symbol:    "INVALID",
		StartDate: "20240101",
		EndDate:   "20240131",
	})

	if err == nil {
		t.Error("expected error for invalid symbol")
	}
}

func TestToEastMoneySecid(t *testing.T) {
	tests := []struct {
		symbol  string
		want    string
		wantErr bool
	}{
		{"600001.SH", "1.600001", false},
		{"000001.SZ", "0.000001", false},
		{"430001.BJ", "0.430001", false},
		{"INVALID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			got, err := toEastMoneySecid(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("toEastMoneySecid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toEastMoneySecid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPeriod(t *testing.T) {
	tests := []struct {
		timeframe string
		want      string
	}{
		{"1m", "1"},
		{"5m", "5"},
		{"15m", "15"},
		{"30m", "30"},
		{"60m", "60"},
		{"1d", "101"},
		{"", "101"},
		{"1w", "102"},
		{"1M", "103"},
		{"unknown", "101"},
	}

	for _, tt := range tests {
		t.Run(tt.timeframe, func(t *testing.T) {
			if got := ToPeriod(tt.timeframe); got != tt.want {
				t.Errorf("ToPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToAdjust(t *testing.T) {
	tests := []struct {
		adj  string
		want string
	}{
		{"qfq", "1"},
		{"hfq", "2"},
		{"", "0"},
		{"other", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.adj, func(t *testing.T) {
			if got := ToAdjust(tt.adj); got != tt.want {
				t.Errorf("ToAdjust() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseKlineLine(t *testing.T) {
	line := "2024-01-02,10.50,10.80,11.00,10.40,1000000,10500000,5.71,2.86,0.30,1.50"
	bar, err := parseKlineLine(line)
	if err != nil {
		t.Fatalf("parseKlineLine() error = %v", err)
	}

	if bar.Date != "2024-01-02" {
		t.Errorf("Date = %v, want 2024-01-02", bar.Date)
	}
	if bar.Open != 10.50 {
		t.Errorf("Open = %v, want 10.50", bar.Open)
	}
	if bar.Close != 10.80 {
		t.Errorf("Close = %v, want 10.80", bar.Close)
	}
	if bar.High != 11.00 {
		t.Errorf("High = %v, want 11.00", bar.High)
	}
	if bar.Low != 10.40 {
		t.Errorf("Low = %v, want 10.40", bar.Low)
	}
}
