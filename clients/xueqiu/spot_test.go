package xueqiu

import (
	"testing"
)

func TestClient_GetSpot(t *testing.T) {
	t.Skip("Xueqiu requires authentication, skipping")
}

func TestParseSpotResponse(t *testing.T) {
	body := []byte(`{"data":[{"symbol":"SZ000001","current":12.5,"open":12.3,"high":12.8,"low":12.1,"last_close":12.2,"chg":0.3,"percent":2.46,"volume":100000,"amount":1250000,"time":1700000000000}]}`)

	result, err := parseSpotResponse(body)
	if err != nil {
		t.Fatalf("parseSpotResponse() error = %v", err)
	}

	if result.Count != 1 {
		t.Errorf("Count = %d, want 1", result.Count)
	}

	if result.Data[0].Symbol != "SZ000001" {
		t.Errorf("Symbol = %s, want SZ000001", result.Data[0].Symbol)
	}

	if result.Data[0].Latest != 12.5 {
		t.Errorf("Latest = %v, want 12.5", result.Data[0].Latest)
	}
}
