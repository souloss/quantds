package finnhub

import (
	"context"
	"testing"
)

func TestClient_GetStockSymbols(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetStockSymbols(context.Background(), &SymbolParams{Exchange: "US"})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("US stock symbols: %d", result.Count)

	if len(result.Symbols) == 0 {
		t.Fatal("Expected symbols, got 0")
	}

	sym := result.Symbols[0]
	t.Logf("First: Symbol=%s, Description=%s, Type=%s", sym.Symbol, sym.Description, sym.Type)

	if sym.Symbol == "" {
		t.Error("Expected non-empty symbol")
	}
}
