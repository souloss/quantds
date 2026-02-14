package eastmoney

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetStockDetail(t *testing.T) {
	client := NewClient()
	defer client.Close()

	tests := []struct {
		name   string
		params *DetailParams
	}{
		{"SZ stock", &DetailParams{Symbol: "000001.SZ"}},
		{"SH stock", &DetailParams{Symbol: "600519.SH"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, record, err := client.GetStockDetail(ctx, tt.params)
			if err != nil {
				checkAPIError(t, err)
				return
			}

			if record == nil {
				t.Fatal("record is nil")
			}

			if result == nil {
				t.Fatal("result is nil")
			}

			t.Logf("Code: %s, Name: %s", result.Code, result.Name)
			t.Logf("Latest: %.2f, Open: %.2f, High: %.2f, Low: %.2f",
				result.LatestPrice, result.Open, result.High, result.Low)
			t.Logf("TotalCap: %.2f亿, FloatCap: %.2f亿",
				result.TotalMarketCap/100000000, result.FloatMarketCap/100000000)

			if result.Code == "" {
				t.Error("code is empty")
			}
			if result.Name == "" {
				t.Error("name is empty")
			}
			if result.LatestPrice <= 0 {
				t.Error("latest price should be positive")
			}
		})
	}
}

func TestClient_GetStockDetail_InvalidSymbol(t *testing.T) {
	client := NewClient()
	defer client.Close()

	ctx := context.Background()
	_, _, err := client.GetStockDetail(ctx, &DetailParams{Symbol: "INVALID"})

	if err == nil {
		t.Error("expected error for invalid symbol")
	}
}
