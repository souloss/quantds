package twelvedata

import (
	"context"
	"testing"
)

func TestClient_GetStocksList(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetStocksList(context.Background(), &ListParams{Country: "United States"})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("US stocks: %d", result.Count)
	if len(result.Data) > 0 {
		t.Logf("First: Symbol=%s, Name=%s, Exchange=%s", result.Data[0].Symbol, result.Data[0].Name, result.Data[0].Exchange)
	}
}

func TestClient_GetForexPairsList(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetForexPairsList(context.Background())
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Forex pairs: %d", result.Count)
	if len(result.Data) > 0 {
		t.Logf("First: Symbol=%s, Name=%s", result.Data[0].Symbol, result.Data[0].Name)
	}
}

func TestClient_GetBondsList(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetBondsList(context.Background(), nil)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Bonds: %d", result.Count)
	if len(result.Data) > 0 {
		t.Logf("First: Symbol=%s, Name=%s, Country=%s", result.Data[0].Symbol, result.Data[0].Name, result.Data[0].Country)
	}
}

func TestClient_GetFundsList(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetFundsList(context.Background(), nil)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Funds: %d", result.Count)
	if len(result.Data) > 0 {
		t.Logf("First: Symbol=%s, Name=%s, Country=%s", result.Data[0].Symbol, result.Data[0].Name, result.Data[0].Country)
	}
}

func TestClient_GetCryptocurrenciesList(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetCryptocurrenciesList(context.Background())
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Cryptocurrencies: %d", result.Count)
	if len(result.Data) > 0 {
		t.Logf("First: Symbol=%s, Name=%s", result.Data[0].Symbol, result.Data[0].Name)
	}
}

func TestClient_GetETFList(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetETFList(context.Background(), nil)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("ETFs: %d", result.Count)
	if len(result.Data) > 0 {
		t.Logf("First: Symbol=%s, Name=%s", result.Data[0].Symbol, result.Data[0].Name)
	}
}
