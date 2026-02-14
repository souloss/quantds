package tushare

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestClient_GetKline(t *testing.T) {
	token := os.Getenv("TUSHARE_TOKEN")
	if token == "" {
		t.Skip("TUSHARE_TOKEN not set")
	}

	client := NewClient(WithToken(token))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, record, err := client.GetKline(ctx, &KlineParams{
		Symbol:    "000001.SZ",
		StartDate: "20240101",
		EndDate:   "20240131",
		Period:    "daily",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetKline() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d bars", result.Count)
	if result.Count > 0 {
		t.Logf("First: date=%s, open=%.2f, close=%.2f", result.Data[0].Date, result.Data[0].Open, result.Data[0].Close)
	}
}

func TestClient_GetStockBasic(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, record, err := client.GetStockBasic(ctx, &StockBasicParams{
		Exchange: "SSE",
		Limit:    10,
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetStockBasic() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d stocks", len(result))
	for i, s := range result {
		t.Logf("Stock[%d]: code=%s, name=%s, industry=%s", i, s.Symbol, s.Name, s.Industry)
	}
}

func TestClient_GetTradeCal(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	now := time.Now()
	result, record, err := client.GetTradeCal(ctx, &TradeCalParams{
		Exchange:  "SSE",
		StartDate: now.AddDate(0, 0, -7).Format("20060102"),
		EndDate:   now.Format("20060102"),
		IsOpen:    "1",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetTradeCal() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d trading days", len(result))
}

func TestToTushareSymbol(t *testing.T) {
	tests := []struct {
		symbol  string
		want    string
		wantErr bool
	}{
		{"600001.SH", "600001.SH", false},
		{"000001.SZ", "000001.SZ", false},
		{"INVALID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			got, err := ToTushareSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToTushareSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToTushareSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPeriod(t *testing.T) {
	tests := []struct {
		timeframe string
		want      string
	}{
		{"1d", "daily"},
		{"", "daily"},
		{"1w", "weekly"},
		{"1M", "monthly"},
	}

	for _, tt := range tests {
		t.Run(tt.timeframe, func(t *testing.T) {
			if got := ToPeriod(tt.timeframe); got != tt.want {
				t.Errorf("ToPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}
