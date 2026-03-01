package alphavantage

import (
	"context"
	"testing"
)

func TestClient_SearchSymbol(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.SearchSymbol(context.Background(), &SearchParams{Keywords: "Apple"})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Search 'Apple': %d matches", result.Count)

	if len(result.Matches) == 0 {
		t.Fatal("Expected search results, got 0")
	}

	m := result.Matches[0]
	t.Logf("First: Symbol=%s, Name=%s, Type=%s, Region=%s",
		m.Symbol, m.Name, m.Type, m.Region)
}
