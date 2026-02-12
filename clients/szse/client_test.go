package szse

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

	result, record, err := client.GetStockList(ctx, nil)
	if err != nil {
		t.Fatalf("GetStockList() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	if result == nil {
		t.Fatal("result is nil")
	}

	t.Logf("Got %d rows (including header)", len(result.Data))

	if len(result.Data) < 2 {
		t.Fatal("expected header + at least 1 data row")
	}

	t.Logf("Header: %v", result.Data[0])
	for i, row := range result.Data[1:min(4, len(result.Data))] {
		t.Logf("Row[%d]: %v", i, row)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
