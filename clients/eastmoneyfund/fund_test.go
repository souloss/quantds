package eastmoneyfund

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Error("NewClient returned nil")
	}
	defer client.Close()
}

func TestClient_GetFundList(t *testing.T) {
	client := NewClient()
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, _, err := client.GetFundList(ctx)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Total funds: %d", result.Count)

	if len(result.Funds) == 0 {
		t.Fatal("Expected funds, got 0")
	}

	fund := result.Funds[0]
	t.Logf("First: Code=%s, Name=%s, Type=%s", fund.Code, fund.Name, fund.Type)

	if fund.Code == "" {
		t.Error("Expected non-empty fund code")
	}
}

func TestClient_GetFundEstimate(t *testing.T) {
	client := NewClient()
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, _, err := client.GetFundEstimate(ctx, &FundEstimateParams{Code: "000001"})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Fund %s (%s): EstNAV=%s, EstChange=%s%%, Time=%s",
		result.Code, result.Name, result.EstNAV, result.EstChange, result.EstTime)

	if result.Code == "" {
		t.Error("Expected non-empty fund code")
	}
}
