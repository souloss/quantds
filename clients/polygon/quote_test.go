package polygon

import (
	"context"
	"testing"
)

func TestClient_GetSnapshot(t *testing.T) {
	skipIfNoAPIKey(t)
	client := NewClient()
	defer client.Close()

	result, _, err := client.GetSnapshot(context.Background(), &SnapshotParams{
		Tickers: []string{"AAPL"},
	})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Snapshot tickers: %d", result.Count)

	if len(result.Tickers) == 0 {
		t.Log("Warning: no snapshot data (may require paid plan)")
		return
	}

	snap := result.Tickers[0]
	t.Logf("AAPL: Day Open=%.2f, Close=%.2f, Change=%.2f", snap.Day.Open, snap.Day.Close, snap.Change)
}
