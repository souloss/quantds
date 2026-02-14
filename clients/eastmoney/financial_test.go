package eastmoney

import (
	"context"
	"testing"
)

// TestClient_GetFinancials tests retrieving financial report data
// API Rule: No authentication required
// Geo-Restriction: May be blocked in some regions
func TestClient_GetFinancials(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	// Test with Ping An Bank (000001.SZ)
	params := &FinancialParams{
		ReportName: ReportBalanceSheet,
		Code:       "000001",
		PageNumber: 1,
		PageSize:   5,
	}

	result, record, err := client.GetFinancials(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Financial Response Status: %d", record.Response.StatusCode)
	t.Logf("Success: %v, Code: %d, Message: %s", result.Success, result.Code, result.Message)
	t.Logf("Got %d financial records", len(result.Data))

	if !result.Success {
		t.Logf("Warning: API returned success=false, code=%d, message=%s", result.Code, result.Message)
		return
	}

	if len(result.Data) == 0 {
		t.Log("Warning: No financial data returned")
		return
	}

	// Log first record
	firstRecord := result.Data[0]
	t.Logf("First record keys: %v", getKeys(firstRecord))
	if reportDate, ok := firstRecord["REPORT_DATE"]; ok {
		t.Logf("Report Date: %v", reportDate)
	}
}

// TestClient_GetBalanceSheet tests retrieving balance sheet data
func TestClient_GetBalanceSheet(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, record, err := client.GetBalanceSheet(ctx, "000001", 1, 3)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Balance Sheet Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d balance sheet records", len(result.Data))

	if len(result.Data) > 0 {
		t.Logf("First record: %+v", result.Data[0])
	}
}

// TestClient_GetIncomeStatement tests retrieving income statement data
func TestClient_GetIncomeStatement(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, record, err := client.GetIncomeStatement(ctx, "000001", 1, 3)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Income Statement Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d income statement records", len(result.Data))

	if len(result.Data) > 0 {
		t.Logf("Sample record keys: %v", getKeys(result.Data[0]))
	}
}

// TestClient_GetCashflowStatement tests retrieving cash flow statement data
func TestClient_GetCashflowStatement(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, record, err := client.GetCashflowStatement(ctx, "000001", 1, 3)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Cash Flow Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d cash flow records", len(result.Data))
}

// TestClient_GetFinancialIndicator tests retrieving financial indicator data
func TestClient_GetFinancialIndicator(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	result, record, err := client.GetFinancialIndicator(ctx, "000001", 1, 3)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Financial Indicator Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d financial indicator records", len(result.Data))

	if len(result.Data) > 0 {
		t.Logf("Sample record keys: %v", getKeys(result.Data[0]))
	}
}

// Helper function to get keys from a map
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
