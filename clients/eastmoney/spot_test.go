//go:build integration
// +build integration

package eastmoney

// Spot 实时行情 API 测试

import (
	"context"
	"testing"
	"time"

	"github.com/souloss/quantds/request"
)

func TestClient_GetSpot(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	tests := []struct {
		name   string
		params *SpotParams
	}{
		{
			name: "get all stocks",
			params: &SpotParams{
				PageNumber: 1,
				PageSize:   20,
			},
		},
		{
			name: "get SH stocks",
			params: &SpotParams{
				PageNumber: 1,
				PageSize:   10,
				Market:     "SH",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, record, err := client.GetSpot(ctx, tt.params)
			if err != nil {
				t.Fatalf("GetSpot() error = %v", err)
			}

			if record == nil {
				t.Fatal("record is nil")
			}

			if !record.IsSuccess() {
				t.Errorf("record should be success, got error: %v", record.Error)
			}

			if result == nil {
				t.Fatal("result is nil")
			}

			if len(result.Data) == 0 {
				t.Fatal("no data returned")
			}

			for i, q := range result.Data {
				if q.Code == "" {
					t.Errorf("quote[%d].Code is empty", i)
				}
				if q.Name == "" {
					t.Errorf("quote[%d].Name is empty", i)
				}
			}

			t.Logf("Got %d quotes (total: %d)", len(result.Data), result.Total)
		})
	}
}
