package eastmoney

import (
	"context"
	"testing"
)

// TestClient_GetMoneyFlow tests retrieving real-time money flow
// API Rule: No authentication required
// Geo-Restriction: May be blocked in some regions
func TestClient_GetMoneyFlow(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	// Test with Ping An Bank (000001.SZ)
	params := &MoneyFlowParams{
		Symbol: "000001.SZ",
	}

	result, record, err := client.GetMoneyFlow(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Money Flow Response Status: %d", record.Response.StatusCode)
	t.Logf("Money Flow Data for 000001.SZ:")
	t.Logf("  Main Net Inflow: %.2f (%.2f%%)", result.MainNetInflow, result.MainNetInflowRatio)
	t.Logf("  Super Large: %.2f", result.SuperNetInflow)
	t.Logf("  Large: %.2f", result.LargeNetInflow)
	t.Logf("  Medium: %.2f", result.MediumNetInflow)
	t.Logf("  Small: %.2f", result.SmallNetInflow)

	if result.MainNetInflow == 0 && result.SuperNetInflow == 0 {
		t.Log("Warning: All money flow values are zero (market may be closed)")
	}
}

// TestClient_GetMoneyFlowHistory tests retrieving historical money flow
func TestClient_GetMoneyFlowHistory(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	// Test with Kweichow Moutai (600519.SH)
	params := &MoneyFlowHistoryParams{
		Symbol: "600519.SH",
		Limit:  5,
	}

	result, record, err := client.GetMoneyFlowHistory(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Money Flow History Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d days of money flow history", len(result))

	if len(result) == 0 {
		t.Log("Warning: No money flow history returned")
		return
	}

	for i, item := range result {
		t.Logf("Day[%d]: date=%s, close=%.2f, change=%.2f%%, main_net=%.2f, main_ratio=%.2f%%",
			i, item.Date, item.ClosePrice, item.ChangePercent, item.MainNetInflow, item.MainNetInflowRatio)
	}
}

// TestClient_GetMoneyFlow_SH tests money flow for Shanghai stocks
func TestClient_GetMoneyFlow_SH(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &MoneyFlowParams{
		Symbol: "600519.SH", // Kweichow Moutai
	}

	result, _, err := client.GetMoneyFlow(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Moutai (600519.SH) Money Flow:")
	t.Logf("  Main: %.2f (%.2f%%)", result.MainNetInflow, result.MainNetInflowRatio)
	t.Logf("  Super Large: %.2f", result.SuperNetInflow)
	t.Logf("  Large: %.2f", result.LargeNetInflow)
}

// TestClient_GetMoneyFlow_SZ tests money flow for Shenzhen stocks
func TestClient_GetMoneyFlow_SZ(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &MoneyFlowParams{
		Symbol: "000001.SZ", // Ping An Bank
	}

	result, _, err := client.GetMoneyFlow(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Ping An Bank (000001.SZ) Money Flow:")
	t.Logf("  Main: %.2f (%.2f%%)", result.MainNetInflow, result.MainNetInflowRatio)
	t.Logf("  Super Large: %.2f", result.SuperNetInflow)
	t.Logf("  Large: %.2f", result.LargeNetInflow)
}
