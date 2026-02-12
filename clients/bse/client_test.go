package bse

import (
	"context"
	"testing"
	"time"

	"github.com/souloss/quantds/request"
)

func TestClient_GetStockListPage(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, record, err := client.GetStockListPage(ctx, &StockListParams{Page: 1})
	if err != nil {
		t.Fatalf("GetStockListPage() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("BSE page 1: %d stocks, total pages: %d", len(result.Data), result.TotalPages)

	if len(result.Data) == 0 {
		t.Fatal("no stocks returned")
	}

	for i, s := range result.Data[:min(3, len(result.Data))] {
		t.Logf("Stock[%d]: code=%s, name=%s, listDate=%s",
			i, s.StockCode, s.StockName, s.ListDate)
	}

	if result.Data[0].StockCode == "" {
		t.Error("first stock code is empty")
	}
}

func TestClient_GetStockList(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rows, records, err := client.GetStockList(ctx)
	if err != nil {
		t.Fatalf("GetStockList() error = %v", err)
	}

	t.Logf("BSE all stocks: %d, requests: %d", len(rows), len(records))

	if len(rows) == 0 {
		t.Fatal("no stocks returned")
	}

	for i, s := range rows[:min(3, len(rows))] {
		t.Logf("Stock[%d]: code=%s, name=%s", i, s.StockCode, s.StockName)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
