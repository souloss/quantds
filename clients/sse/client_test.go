package sse

import (
	"context"
	"testing"
	"time"

	"github.com/souloss/quantds/request"
)

func TestClient_GetStockList(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, record, err := client.GetStockList(ctx, &StockListParams{PageSize: "50"})
	if err != nil {
		t.Fatalf("GetStockList() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	if result == nil {
		t.Fatal("result is nil")
	}

	t.Logf("Got %d SSE stocks", len(result.Data))

	if len(result.Data) == 0 {
		t.Fatal("no stocks returned")
	}

	for i, s := range result.Data[:min(3, len(result.Data))] {
		t.Logf("Stock[%d]: code=%s, name=%s, listDate=%s, industry=%s",
			i, s.CompanyCode, s.CompanyAbbr, s.ListDate, s.Industry)
	}

	if result.Data[0].CompanyCode == "" {
		t.Error("first stock code is empty")
	}
	if result.Data[0].CompanyAbbr == "" {
		t.Error("first stock name is empty")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
