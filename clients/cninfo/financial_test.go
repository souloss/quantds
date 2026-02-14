package cninfo

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetFinancialData tests retrieving financial data
// API Rule: No authentication required, but needs correct headers
// Geo-Restriction: Cninfo API may have restrictions
func TestClient_GetFinancialData(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &FinancialParams{
		StockCode: "000001", // Ping An Bank
		PageSize:  5,
	}

	result, record, err := client.GetFinancialData(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Financial Data Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d financial records (total: %d)", len(result.Data), result.Total)

	if len(result.Data) == 0 {
		t.Log("Warning: No financial data returned")
		return
	}

	for i, row := range result.Data {
		t.Logf("Record[%d]: code=%s, name=%s, report_date=%s",
			i, row.StockCode, row.StockName, row.ReportDate)
		t.Logf("         revenue=%.2f, net_profit=%.2f, total_assets=%.2f",
			row.Revenue, row.NetProfit, row.TotalAssets)
		t.Logf("         EPS=%.4f, operating_CF=%.2f",
			row.BasicEPS, row.OperatingCashFlow)
	}
}

// TestClient_GetProfile tests retrieving company profile
func TestClient_GetProfile(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &ProfileParams{
		StockCode: "000001", // Ping An Bank
	}

	result, _, err := client.GetProfile(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Company Profile for %s (%s):", result.StockCode, result.StockName)
	t.Logf("  Full Name: %s", result.FullName)
	t.Logf("  English Name: %s", result.EnglishName)
	t.Logf("  Industry: %s", result.Industry)
	t.Logf("  Province: %s, City: %s", result.Province, result.City)
	t.Logf("  Chairman: %s", result.Chairman)
	t.Logf("  Manager: %s", result.Manager)
	t.Logf("  List Date: %s", result.ListDate)
	t.Logf("  Reg Capital: %s", result.RegCapital)
	t.Logf("  Website: %s", result.Website)
	t.Logf("  Main Business: %s", truncate(result.MainBusiness, 100))
}

// TestClient_GetFinancialData_Multiple tests financial data for multiple stocks
func TestClient_GetFinancialData_Multiple(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	codes := []string{"000001", "600519", "000858"}

	for _, code := range codes {
		result, _, err := client.GetFinancialData(ctx, &FinancialParams{
			StockCode: code,
			PageSize:  1,
		})
		if err != nil {
			t.Logf("Error for %s: %v", code, err)
			continue
		}

		if len(result.Data) > 0 {
			row := result.Data[0]
			t.Logf("%s (%s): revenue=%.2f, net_profit=%.2f",
				row.StockCode, row.StockName, row.Revenue, row.NetProfit)
		}
	}
}

// Helper function to truncate string
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
