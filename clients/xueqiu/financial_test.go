package xueqiu

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetIncome tests retrieving income statement
// Note: Xueqiu API may require authentication
func TestClient_GetIncome(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &FinancialParams{
		Symbol: "000001.SZ",
		Type:   "Q4", // Annual report
		Count:  3,
	}

	result, record, err := client.GetIncome(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Income Statement Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d income records", len(result))

	if len(result) == 0 {
		t.Log("Warning: No income data returned (may require authentication)")
		return
	}

	for i, item := range result {
		t.Logf("Income[%d]: report=%s, date=%d", i, item.ReportName, item.ReportDate)
		t.Logf("  Values: %v", item.Values)
	}
}

// TestClient_GetBalance tests retrieving balance sheet
func TestClient_GetBalance(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &FinancialParams{
		Symbol: "600519.SH",
		Count:  3,
	}

	result, _, err := client.GetBalance(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d balance sheet records", len(result))
	for i, item := range result {
		t.Logf("Balance[%d]: report=%s", i, item.ReportName)
	}
}

// TestClient_GetCashFlow tests retrieving cash flow statement
func TestClient_GetCashFlow(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &FinancialParams{
		Symbol: "000001.SZ",
		Count:  3,
	}

	result, _, err := client.GetCashFlow(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d cash flow records", len(result))
	for i, item := range result {
		t.Logf("CashFlow[%d]: report=%s", i, item.ReportName)
	}
}

// TestClient_GetIndicator tests retrieving financial indicators
func TestClient_GetIndicator(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &FinancialParams{
		Symbol: "600519.SH",
		Count:  3,
	}

	result, _, err := client.GetIndicator(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Got %d indicator records", len(result))
	for i, item := range result {
		t.Logf("Indicator[%d]: report=%s", i, item.ReportName)
	}
}
