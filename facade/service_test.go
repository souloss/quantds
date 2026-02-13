package facade

import (
	"context"
	"testing"
	"time"

	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/domain/spot"
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

// TestService_USMarket tests US stock market (Yahoo)
func TestService_USMarket(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx := context.Background()

	// Test US stock K-line
	t.Run("US stock kline", func(t *testing.T) {
		req := kline.Request{
			Symbol:    "AAPL.US",
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
	})

	// Test US stock spot
	t.Run("US stock spot", func(t *testing.T) {
		req := spot.Request{
			Symbols: []string{"AAPL.US", "MSFT.US"},
		}

		result, err := svc.GetSpot(ctx, req)
		if err != nil {
			t.Fatalf("GetSpot() error = %v", err)
		}

		if result.Source == "" {
			t.Error("result.Source is empty")
		}

		t.Logf("Got %d quotes from %s", len(result.Quotes), result.Source)
	})
}

// TestService_HKMarket tests Hong Kong stock market (EastMoney HK)
func TestService_HKMarket(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx := context.Background()

	// Test HK stock K-line
	t.Run("HK stock kline", func(t *testing.T) {
		req := kline.Request{
			Symbol:    "00700.HK",
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
	})

	// Test HK stock spot
	t.Run("HK stock spot", func(t *testing.T) {
		req := spot.Request{
			Symbols: []string{"00700.HK", "00941.HK"},
		}

		result, err := svc.GetSpot(ctx, req)
		if err != nil {
			t.Fatalf("GetSpot() error = %v", err)
		}

		if result.Source == "" {
			t.Error("result.Source is empty")
		}

		t.Logf("Got %d quotes from %s", len(result.Quotes), result.Source)
	})
}

// TestService_USInstruments tests US stock instrument list
func TestService_USInstruments(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx := context.Background()

	t.Run("US stock instruments", func(t *testing.T) {
		req := instrument.Request{
			Market:   "US",
			PageSize: 50,
		}

		result, err := svc.GetInstruments(ctx, req)
		if err != nil {
			t.Fatalf("GetInstruments() error = %v", err)
		}

		if result.Source == "" {
			t.Error("result.Source is empty")
		}

		t.Logf("Got %d instruments from %s", len(result.Data), result.Source)
	})
}

// TestService_HKInstruments tests HK stock instrument list
func TestService_HKInstruments(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx := context.Background()

	t.Run("HK stock instruments", func(t *testing.T) {
		req := instrument.Request{
			Market:   "HK",
			PageSize: 50,
		}

		result, err := svc.GetInstruments(ctx, req)
		if err != nil {
			t.Fatalf("GetInstruments() error = %v", err)
		}

		if result.Source == "" {
			t.Error("result.Source is empty")
		}

		t.Logf("Got %d instruments from %s", len(result.Data), result.Source)
	})
}

// TestService_CryptoInstruments tests cryptocurrency instrument list
func TestService_CryptoInstruments(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx := context.Background()

	t.Run("crypto instruments", func(t *testing.T) {
		req := instrument.Request{
			Market:   "USDT",
			PageSize: 50,
		}

		result, err := svc.GetInstruments(ctx, req)
		if err != nil {
			t.Fatalf("GetInstruments() error = %v", err)
		}

		if result.Source == "" {
			t.Error("result.Source is empty")
		}

		t.Logf("Got %d instruments from %s", len(result.Data), result.Source)
	})
}

// TestService_CryptoMarket tests cryptocurrency market (Binance)
func TestService_CryptoMarket(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx := context.Background()

	// Test crypto K-line
	t.Run("crypto kline", func(t *testing.T) {
		req := kline.Request{
			Symbol:    "BTCUSDT",
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
	})

	// Test crypto spot
	t.Run("crypto spot", func(t *testing.T) {
		req := spot.Request{
			Symbols: []string{"BTCUSDT", "ETHUSDT"},
		}

		result, err := svc.GetSpot(ctx, req)
		if err != nil {
			t.Fatalf("GetSpot() error = %v", err)
		}

		if result.Source == "" {
			t.Error("result.Source is empty")
		}

		t.Logf("Got %d quotes from %s", len(result.Quotes), result.Source)
	})
}
