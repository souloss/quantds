package eodhd

import (
	"context"
	"testing"
)

func TestClient_GetExchangeSymbolList(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetExchangeSymbolList(context.Background(), &ExchangeSymbolsParams{Exchange: "US"})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("US exchange symbols: %d", result.Count)

	if len(result.Symbols) == 0 {
		t.Fatal("Expected symbols, got 0")
	}

	sym := result.Symbols[0]
	t.Logf("First: Code=%s, Name=%s, Type=%s, Currency=%s",
		sym.Code, sym.Name, sym.Type, sym.Currency)
}
