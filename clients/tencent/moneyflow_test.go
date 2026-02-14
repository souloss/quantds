package tencent

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetMoneyFlow tests retrieving real-time money flow
// API Rule: No authentication required
// Geo-Restriction: Tencent Finance may be blocked outside China
func TestClient_GetMoneyFlow(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	params := &MoneyFlowParams{
		Symbol: "600519.SH", // Kweichow Moutai
	}

	result, record, err := client.GetMoneyFlow(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Money Flow Response Status: %d", record.Response.StatusCode)
	t.Logf("Money Flow for %s (%s):", result.Code, result.Name)
	t.Logf("  Main: in=%.2f, out=%.2f, net=%.2f, ratio=%.2f%%",
		result.MainIn, result.MainOut, result.MainNet, result.MainRatio)
	t.Logf("  Super Large: net=%.2f, ratio=%.2f%%", result.SuperNet, result.SuperRatio)
	t.Logf("  Large: net=%.2f, ratio=%.2f%%", result.LargeNet, result.LargeRatio)
	t.Logf("  Medium: net=%.2f, ratio=%.2f%%", result.MediumNet, result.MediumRatio)
	t.Logf("  Small: net=%.2f, ratio=%.2f%%", result.SmallNet, result.SmallRatio)
	t.Logf("  Date: %s", result.Date)

	if result.MainNet == 0 && result.SuperNet == 0 {
		t.Log("Warning: All money flow values are zero (market may be closed)")
	}
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
	t.Logf("  Main Net: %.2f (%.2f%%)", result.MainNet, result.MainRatio)
	t.Logf("  Super Large Net: %.2f", result.SuperNet)
	t.Logf("  Large Net: %.2f", result.LargeNet)
	t.Logf("  Medium Net: %.2f", result.MediumNet)
	t.Logf("  Small Net: %.2f", result.SmallNet)
}

// TestClient_GetMoneyFlow_Multiple tests money flow for multiple stocks
func TestClient_GetMoneyFlow_Multiple(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	symbols := []string{"600519.SH", "000001.SZ", "601318.SH"}

	for _, symbol := range symbols {
		result, _, err := client.GetMoneyFlow(ctx, &MoneyFlowParams{Symbol: symbol})
		if err != nil {
			t.Logf("Error for %s: %v", symbol, err)
			continue
		}

		t.Logf("%s (%s): main_net=%.2f, main_ratio=%.2f%%",
			result.Code, result.Name, result.MainNet, result.MainRatio)
	}
}
