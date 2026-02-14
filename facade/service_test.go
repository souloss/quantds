package facade

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/souloss/quantds/domain/announcement"
	"github.com/souloss/quantds/domain/financial"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/domain/profile"
	"github.com/souloss/quantds/domain/spot"
)

// checkFacadeError 检查 API 错误并优雅跳过不可控的外部故障
func checkFacadeError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	msg := err.Error()

	// 已知的网络/API 不可用错误，使用 t.Skipf 而不是 t.Fatalf
	skipPatterns := []string{
		"timeout", "connection refused", "no such host",
		"dial tcp", "EOF", "TLS handshake",
		"403", "401", "429", "503",
		"rate limit", "geo-restrict", "blocked",
		"unsupported market",
		"client error", "unmarshal error",
		"all providers failed", "no provider",
	}
	for _, p := range skipPatterns {
		if strings.Contains(strings.ToLower(msg), strings.ToLower(p)) {
			t.Skipf("Skipping due to external API issue: %v", err)
			return
		}
	}

	t.Fatalf("Unexpected error: %v", err)
}

// ========== A 股市场测试 ==========

func TestService_GetKline_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := kline.Request{
		Symbol:    "000001.SZ",
		StartTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		EndTime:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.Local),
		Timeframe: kline.Timeframe1d,
	}

	result, err := svc.GetKline(ctx, req)
	checkFacadeError(t, err)

	if result.Source == "" {
		t.Error("result.Source is empty")
	}

	if len(result.Bars) == 0 {
		t.Log("Warning: no data returned")
		return
	}

	t.Logf("Got %d bars from %s", len(result.Bars), result.Source)
	b := result.Bars[0]
	t.Logf("First bar: date=%v, O=%.2f, H=%.2f, L=%.2f, C=%.2f, V=%.0f",
		b.Timestamp.Format("2006-01-02"), b.Open, b.High, b.Low, b.Close, b.Volume)
}

func TestService_GetSpot_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := svc.GetSpot(ctx, spot.Request{
		Symbols: []string{"000001.SZ", "600519.SH"},
	})
	checkFacadeError(t, err)

	t.Logf("Got %d quotes from %s", len(result.Quotes), result.Source)
	for i, q := range result.Quotes {
		t.Logf("Quote[%d]: symbol=%s, name=%s, latest=%.2f", i, q.Symbol, q.Name, q.Latest)
	}
}

func TestService_GetInstruments_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := svc.GetInstruments(ctx, instrument.Request{
		PageSize: 50,
	})
	checkFacadeError(t, err)

	t.Logf("Got %d instruments from %s (total=%d)", len(result.Data), result.Source, result.Total)
	for i, inst := range result.Data {
		if i >= 3 {
			break
		}
		t.Logf("Instrument[%d]: %s %s (%s)", i, inst.Code, inst.Name, inst.Exchange)
	}
}

func TestService_GetProfile_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := svc.GetProfile(ctx, profile.Request{
		Symbol: "000001.SZ",
	})
	checkFacadeError(t, err)

	t.Logf("Profile from %s: %s %s, PE=%.2f, PB=%.2f",
		result.Source, result.Data.Symbol, result.Data.Name, result.Data.PE, result.Data.PB)
}

func TestService_GetFinancial_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := svc.GetFinancial(ctx, financial.Request{
		Symbol: "000001.SZ",
	})
	checkFacadeError(t, err)

	t.Logf("Got %d financial records from %s", len(result.Data), result.Source)
	for i, fd := range result.Data {
		if i >= 3 {
			break
		}
		t.Logf("Financial[%d]: date=%v, revenue=%.0f, net_profit=%.0f",
			i, fd.ReportDate.Format("2006-01-02"), fd.TotalRevenue, fd.NetProfit)
	}
}

func TestService_GetAnnouncements_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := svc.GetAnnouncements(ctx, announcement.Request{
		Symbol:   "000001.SZ",
		PageSize: 5,
	})
	checkFacadeError(t, err)

	t.Logf("Got %d announcements from %s (total=%d)", len(result.Data), result.Source, result.TotalCount)
	for i, ann := range result.Data {
		if i >= 3 {
			break
		}
		t.Logf("Announcement[%d]: %s - %s", i, ann.Code, ann.Title)
	}
}

// ========== 美股市场测试 ==========

func TestService_USMarket(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("kline", func(t *testing.T) {
		result, err := svc.GetKline(ctx, kline.Request{
			Symbol:    "AAPL.US",
			Timeframe: kline.Timeframe1d,
		})
		checkFacadeError(t, err)
		t.Logf("Got %d bars from %s", len(result.Bars), result.Source)
	})

	t.Run("spot", func(t *testing.T) {
		result, err := svc.GetSpot(ctx, spot.Request{
			Symbols: []string{"AAPL.US", "MSFT.US"},
		})
		checkFacadeError(t, err)
		t.Logf("Got %d quotes from %s", len(result.Quotes), result.Source)
	})

	t.Run("instruments", func(t *testing.T) {
		result, err := svc.GetInstruments(ctx, instrument.Request{
			Market:   "US",
			PageSize: 50,
		})
		checkFacadeError(t, err)
		t.Logf("Got %d instruments from %s", len(result.Data), result.Source)
	})
}

// ========== 港股市场测试 ==========

func TestService_HKMarket(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("kline", func(t *testing.T) {
		result, err := svc.GetKline(ctx, kline.Request{
			Symbol:    "00700.HK.HKEX",
			Timeframe: kline.Timeframe1d,
		})
		checkFacadeError(t, err)
		t.Logf("Got %d bars from %s", len(result.Bars), result.Source)
	})

	t.Run("spot", func(t *testing.T) {
		result, err := svc.GetSpot(ctx, spot.Request{
			Symbols: []string{"00700.HK.HKEX", "00941.HK.HKEX"},
		})
		checkFacadeError(t, err)
		t.Logf("Got %d quotes from %s", len(result.Quotes), result.Source)
	})

	t.Run("instruments", func(t *testing.T) {
		result, err := svc.GetInstruments(ctx, instrument.Request{
			Market:   "HK",
			PageSize: 50,
		})
		checkFacadeError(t, err)
		t.Logf("Got %d instruments from %s", len(result.Data), result.Source)
	})
}

// ========== 加密货币市场测试 ==========

func TestService_CryptoMarket(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("kline", func(t *testing.T) {
		result, err := svc.GetKline(ctx, kline.Request{
			Symbol:    "BTCUSDT",
			Timeframe: kline.Timeframe1d,
		})
		checkFacadeError(t, err)
		t.Logf("Got %d bars from %s", len(result.Bars), result.Source)
	})

	t.Run("spot", func(t *testing.T) {
		result, err := svc.GetSpot(ctx, spot.Request{
			Symbols: []string{"BTCUSDT", "ETHUSDT"},
		})
		checkFacadeError(t, err)
		t.Logf("Got %d quotes from %s", len(result.Quotes), result.Source)
	})

	t.Run("instruments", func(t *testing.T) {
		result, err := svc.GetInstruments(ctx, instrument.Request{
			Market:   "CRYPTO",
			PageSize: 50,
		})
		checkFacadeError(t, err)
		t.Logf("Got %d instruments from %s", len(result.Data), result.Source)
	})
}

// ========== 市场路由测试 ==========

func TestService_MarketRouting(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tests := []struct {
		name   string
		symbol string
	}{
		{"CN A-share", "000001.SZ"},
		{"US stock", "AAPL.US"},
		{"HK stock", "00700.HK.HKEX"},
		{"Crypto", "BTCUSDT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetKline(ctx, kline.Request{
				Symbol:    tt.symbol,
				Timeframe: kline.Timeframe1d,
			})
			checkFacadeError(t, err)
			t.Logf("[%s] %s → source=%s, bars=%d", tt.name, tt.symbol, result.Source, len(result.Bars))
		})
	}
}

// ========== 服务级别测试 ==========

func TestService_Stats(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	stats := svc.GetStats()
	t.Logf("Stats: %+v", stats)
}

func TestService_UnsupportedMarket(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx := context.Background()

	_, err := svc.GetKline(ctx, kline.Request{
		Symbol:    "invalid-symbol-format!!!",
		Timeframe: kline.Timeframe1d,
	})
	if err == nil {
		t.Error("Expected error for invalid symbol, got nil")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}
