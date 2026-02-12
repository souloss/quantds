package eastmoney

import (
	"context"
	"testing"
	"time"

	"github.com/souloss/quantds/request"
)

func TestClient_GetStockList(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	tests := []struct {
		name   string
		params *StockListParams
	}{
		{"all stocks", &StockListParams{PageSize: 100}},
		{"SH stocks", &StockListParams{Market: "SH", PageSize: 50}},
		{"SZ stocks", &StockListParams{Market: "SZ", PageSize: 50}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, record, err := client.GetStockList(ctx, tt.params)
			if err != nil {
				t.Fatalf("GetStockList() error = %v", err)
			}

			if record == nil {
				t.Fatal("record is nil")
			}

			if result == nil {
				t.Fatal("result is nil")
			}

			t.Logf("Got %d stocks (total: %d)", len(result.Data), result.Total)

			if len(result.Data) == 0 {
				t.Fatal("no stocks returned")
			}

			for i, s := range result.Data[:min(3, len(result.Data))] {
				t.Logf("Stock[%d]: code=%s, name=%s, marketID=%d", i, s.Code, s.Name, s.MarketID)
			}

			if result.Data[0].Code == "" {
				t.Error("first stock code is empty")
			}
			if result.Data[0].Name == "" {
				t.Error("first stock name is empty")
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
