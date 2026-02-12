package facade

import (
	"context"
	"testing"
	"time"

	"github.com/souloss/quantds/domain/kline"
)

func TestService_GetKline(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	tests := []struct {
		name string
		req  kline.Request
	}{
		{
			name: "get daily kline",
			req: kline.Request{
				Symbol:    "000001.SZ",
				StartTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
				EndTime:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.Local),
				Timeframe: kline.Timeframe1d,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := svc.GetKline(ctx, tt.req)
			if err != nil {
				t.Fatalf("GetKline() error = %v", err)
			}

			if result.Source == "" {
				t.Error("result.Source is empty")
			}

			if len(result.Bars) == 0 {
				t.Log("Warning: no data returned")
				return
			}

			t.Logf("Got %d bars from %s", len(result.Bars), result.Source)
		})
	}
}

func TestService_MarketRouting(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	// 测试 A 股市场路由
	ctx := context.Background()
	req := kline.Request{
		Symbol:    "000001.SZ",
		StartTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		EndTime:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.Local),
		Timeframe: kline.Timeframe1d,
	}

	result, err := svc.GetKline(ctx, req)
	if err != nil {
		t.Fatalf("GetKline() error = %v", err)
	}

	if result.Source == "" {
		t.Error("result.Source is empty")
	}

	t.Logf("Got %d bars from %s", len(result.Bars), result.Source)
}
