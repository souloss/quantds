package polygon

import (
	"context"
	"testing"
)

func TestClient_GetTickers(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetTickers(context.Background(), &TickerParams{
		Market: "stocks",
		Limit:  20,
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Stock tickers: %d", result.Count)

	if len(result.Tickers) == 0 {
		t.Fatal("Expected tickers, got 0")
	}

	tk := result.Tickers[0]
	t.Logf("First: Ticker=%s, Name=%s, Market=%s, Exchange=%s",
		tk.Ticker, tk.Name, tk.Market, tk.PrimaryExchange)
}
