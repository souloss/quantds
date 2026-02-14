package xueqiu

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetStockList tests retrieving stock list
// Note: Xueqiu API may require authentication
func TestClient_GetStockList(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &StockListParams{
		Market: "CN", // China market
		Size:   10,
	}

	result, record, err := client.GetStockList(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Stock List Response Status: %d", record.Response.StatusCode)
	t.Logf("Total: %d, Got: %d items", result.Total, len(result.Items))

	if len(result.Items) == 0 {
		t.Log("Warning: No stocks returned (may require authentication)")
		return
	}

	for i, stock := range result.Items {
		t.Logf("Stock[%d]: code=%s, name=%s", i, stock.Symbol, stock.Name)
	}
}

// TestClient_GetStockList_SmallPage tests retrieving a small page of stock list
func TestClient_GetStockList_SmallPage(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, _, err := client.GetStockList(ctx, &StockListParams{
		Market: "CN",
		Size:   5,
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d items", len(result.Items))
	for _, item := range result.Items {
		t.Logf("  %s: %s", item.Symbol, item.Name)
	}
}
