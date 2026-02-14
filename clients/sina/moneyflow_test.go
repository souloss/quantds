package sina

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetMoneyFlow tests retrieving money flow data
// API Rule: No authentication required
// Geo-Restriction: Sina Finance may be blocked outside China
func TestClient_GetMoneyFlow(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &MoneyFlowParams{
		Symbol: "600519.SH", // Kweichow Moutai
		Count:  5,
	}

	result, record, err := client.GetMoneyFlow(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Money Flow Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d days of money flow data", len(result))

	if len(result) == 0 {
		t.Log("Warning: No money flow data returned (market may be closed or symbol invalid)")
		return
	}

	for i, item := range result {
		t.Logf("Day[%d]: date=%s, close=%.2f, change=%.2f%%, net_inflow=%.2f, ratio=%.2f%%",
			i, item.Date, item.Trade, item.ChangeRatio, item.NetAmount, item.RatioAmount)
		t.Logf("         super_large=%.2f, large=%.2f, medium=%.2f, small=%.2f",
			item.R0Net, item.R1Net, item.R2Net, item.R3Net)
	}
}

// TestClient_GetMoneyFlow_SZ tests money flow for Shenzhen stocks
func TestClient_GetMoneyFlow_SZ(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &MoneyFlowParams{
		Symbol: "000001.SZ", // Ping An Bank
		Count:  3,
	}

	result, _, err := client.GetMoneyFlow(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Ping An Bank Money Flow: %d days", len(result))
	for _, item := range result {
		t.Logf("  %s: close=%.2f, main_net=%.2f, ratio=%.2f%%",
			item.Date, item.Trade, item.NetAmount, item.RatioAmount)
	}
}

// TestClient_GetMoneyFlow_Recent tests recent money flow data
func TestClient_GetMoneyFlow_Recent(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &MoneyFlowParams{
		Symbol: "601318.SH", // Ping An Insurance
		Count:  10,
	}

	result, _, err := client.GetMoneyFlow(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Ping An Insurance Money Flow: %d days", len(result))

	// Calculate total net inflow
	var totalNetInflow float64
	for _, item := range result {
		totalNetInflow += item.NetAmount
	}
	t.Logf("Total net inflow over %d days: %.2f", len(result), totalNetInflow)

	// Find the day with largest main net inflow
	if len(result) > 0 {
		maxInflow := result[0]
		for _, item := range result {
			if item.NetAmount > maxInflow.NetAmount {
				maxInflow = item
			}
		}
		t.Logf("Largest main net inflow: %s with %.2f (%.2f%%)",
			maxInflow.Date, maxInflow.NetAmount, maxInflow.RatioAmount)
	}
}
