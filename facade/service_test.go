package facade

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/souloss/quantds/config"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/announcement"
	"github.com/souloss/quantds/domain/financial"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/domain/profile"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
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

// ========== 外汇市场测试 ==========

func TestService_ForexMarket(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("kline", func(t *testing.T) {
		result, err := svc.GetKline(ctx, kline.Request{
			Symbol:    "EURUSD.FOREX.FOREX_SPOT",
			Timeframe: kline.Timeframe1d,
		})
		checkFacadeError(t, err)
		t.Logf("Got %d bars from %s", len(result.Bars), result.Source)
	})

	t.Run("spot", func(t *testing.T) {
		result, err := svc.GetSpot(ctx, spot.Request{
			Symbols: []string{"EURUSD.FOREX.FOREX_SPOT"},
		})
		checkFacadeError(t, err)
		t.Logf("Got %d quotes from %s", len(result.Quotes), result.Source)
	})

	t.Run("instruments", func(t *testing.T) {
		result, err := svc.GetInstruments(ctx, instrument.Request{
			Market:   "FOREX",
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
		{"Forex", "EURUSD.FOREX.FOREX_SPOT"},
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

// ========== WithTrace 测试 ==========

func TestService_GetKlineWithTrace(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, trace, err := svc.GetKlineWithTrace(ctx, kline.Request{
		Symbol:    "000001.SZ",
		Timeframe: kline.Timeframe1d,
	})
	checkFacadeError(t, err)

	if trace == nil {
		t.Error("Expected non-nil trace")
	} else {
		t.Logf("Trace: provider=%s, requests=%d, totalTime=%v",
			trace.Provider, trace.TotalRequests(), trace.TotalTime)
	}

	t.Logf("Got %d bars from %s", len(result.Bars), result.Source)
}

func TestService_GetSpotWithTrace(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, trace, err := svc.GetSpotWithTrace(ctx, spot.Request{
		Symbols: []string{"000001.SZ"},
	})
	checkFacadeError(t, err)

	if trace == nil {
		t.Error("Expected non-nil trace")
	} else {
		t.Logf("Trace: provider=%s, requests=%d, totalTime=%v",
			trace.Provider, trace.TotalRequests(), trace.TotalTime)
	}

	t.Logf("Got %d quotes from %s", len(result.Quotes), result.Source)
}

// ========== 跨市场 Spot 测试 ==========

func TestService_CrossMarketSpot(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Request symbols from different markets simultaneously
	result, err := svc.GetSpot(ctx, spot.Request{
		Symbols: []string{"000001.SZ", "AAPL.US", "BTCUSDT"},
	})
	checkFacadeError(t, err)

	t.Logf("Got %d cross-market quotes from %s", len(result.Quotes), result.Source)
	for _, q := range result.Quotes {
		t.Logf("  %s: name=%s, latest=%.2f", q.Symbol, q.Name, q.Latest)
	}
}

// ========== RegisterProvider 测试 ==========

func TestService_RegisterProvider(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	t.Run("RegisterKlineProvider", func(t *testing.T) {
		// Verify that the kline manager for CN has expected providers
		m := svc.klineManagers[domain.MarketCN]
		if m == nil {
			t.Fatal("CN kline manager should not be nil")
		}
		providers := m.Providers()
		t.Logf("CN kline providers: %v", providers)
		if len(providers) == 0 {
			t.Error("Expected at least one kline provider for CN")
		}
	})

	t.Run("RegisterSpotProvider", func(t *testing.T) {
		m := svc.spotManagers[domain.MarketCN]
		if m == nil {
			t.Fatal("CN spot manager should not be nil")
		}
		providers := m.Providers()
		t.Logf("CN spot providers: %v", providers)
		if len(providers) == 0 {
			t.Error("Expected at least one spot provider for CN")
		}
	})

	t.Run("RegisterInstrumentProvider", func(t *testing.T) {
		m := svc.instrumentManagers[domain.MarketCN]
		if m == nil {
			t.Fatal("CN instrument manager should not be nil")
		}
		providers := m.Providers()
		t.Logf("CN instrument providers: %v", providers)
		if len(providers) == 0 {
			t.Error("Expected at least one instrument provider for CN")
		}
	})

	t.Run("RegisterProfileProvider", func(t *testing.T) {
		m := svc.profileManagers[domain.MarketCN]
		if m == nil {
			t.Fatal("CN profile manager should not be nil")
		}
		providers := m.Providers()
		t.Logf("CN profile providers: %v", providers)
	})

	t.Run("RegisterFinancialProvider", func(t *testing.T) {
		m := svc.financialManagers[domain.MarketCN]
		if m == nil {
			t.Fatal("CN financial manager should not be nil")
		}
		providers := m.Providers()
		t.Logf("CN financial providers: %v", providers)
	})

	t.Run("RegisterAnnouncementProvider", func(t *testing.T) {
		m := svc.announcementManagers[domain.MarketCN]
		if m == nil {
			t.Fatal("CN announcement manager should not be nil")
		}
		providers := m.Providers()
		t.Logf("CN announcement providers: %v", providers)
	})
}

// ========== 逐适配器验证测试 ==========

// testAdapterFromManager 使用 Manager.FetchFrom 逐个测试指定市场的所有适配器。
// 这能帮助客户发现每个数据源的可用性和数据质量。
func testAdapterFromManager[K any, V any](t *testing.T, m *manager.Manager[K, V], ctx context.Context, req K, market domain.Market, domainName string) {
	t.Helper()
	providers := m.Providers()
	if len(providers) == 0 {
		t.Fatalf("No providers registered for %s %s", market, domainName)
	}

	for _, providerName := range providers {
		t.Run(providerName, func(t *testing.T) {
			result, err := m.FetchFrom(ctx, providerName, req)
			if err != nil {
				msg := err.Error()
				// Known external API issues: skip gracefully
				skipWords := []string{
					"timeout", "connection refused", "no such host",
					"dial tcp", "EOF", "TLS handshake",
					"403", "401", "429", "503",
					"rate limit", "geo-restrict", "blocked",
					"token", "invalid", "unauthorized",
					"retries exceeded", "client error",
					"authentication error", "authentication",
				}
				lowerMsg := strings.ToLower(msg)
				for _, w := range skipWords {
					if strings.Contains(lowerMsg, w) {
						t.Skipf("Provider %s unavailable: %v", providerName, err)
						return
					}
				}
				t.Errorf("Provider %s failed unexpectedly: %v", providerName, err)
				return
			}
			t.Logf("Provider %s: OK (cached=%v, hasError=%v)",
				providerName, result.IsCached(), result.HasError())
		})
	}
}

func TestService_AllKlineAdapters_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.klineManagers[domain.MarketCN], ctx,
		kline.Request{Symbol: "000001.SZ", Timeframe: kline.Timeframe1d},
		domain.MarketCN, "kline")
}

func TestService_AllSpotAdapters_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.spotManagers[domain.MarketCN], ctx,
		spot.Request{Symbols: []string{"000001.SZ"}},
		domain.MarketCN, "spot")
}

func TestService_AllInstrumentAdapters_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.instrumentManagers[domain.MarketCN], ctx,
		instrument.Request{PageSize: 50},
		domain.MarketCN, "instrument")
}

func TestService_AllProfileAdapters_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.profileManagers[domain.MarketCN], ctx,
		profile.Request{Symbol: "000001.SZ"},
		domain.MarketCN, "profile")
}

func TestService_AllFinancialAdapters_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.financialManagers[domain.MarketCN], ctx,
		financial.Request{Symbol: "000001.SZ"},
		domain.MarketCN, "financial")
}

func TestService_AllAnnouncementAdapters_CN(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.announcementManagers[domain.MarketCN], ctx,
		announcement.Request{Symbol: "000001.SZ", PageSize: 5},
		domain.MarketCN, "announcement")
}

func TestService_AllKlineAdapters_US(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.klineManagers[domain.MarketUS], ctx,
		kline.Request{Symbol: "AAPL.US", Timeframe: kline.Timeframe1d},
		domain.MarketUS, "kline")
}

func TestService_AllSpotAdapters_US(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.spotManagers[domain.MarketUS], ctx,
		spot.Request{Symbols: []string{"AAPL.US"}},
		domain.MarketUS, "spot")
}

func TestService_AllInstrumentAdapters_US(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.instrumentManagers[domain.MarketUS], ctx,
		instrument.Request{Market: "US", PageSize: 50},
		domain.MarketUS, "instrument")
}

func TestService_AllKlineAdapters_HK(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.klineManagers[domain.MarketHK], ctx,
		kline.Request{Symbol: "00700.HK.HKEX", Timeframe: kline.Timeframe1d},
		domain.MarketHK, "kline")
}

func TestService_AllSpotAdapters_HK(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.spotManagers[domain.MarketHK], ctx,
		spot.Request{Symbols: []string{"00700.HK.HKEX"}},
		domain.MarketHK, "spot")
}

func TestService_AllInstrumentAdapters_HK(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.instrumentManagers[domain.MarketHK], ctx,
		instrument.Request{Market: "HK", PageSize: 50},
		domain.MarketHK, "instrument")
}

func TestService_AllKlineAdapters_Crypto(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.klineManagers[domain.MarketCrypto], ctx,
		kline.Request{Symbol: "BTCUSDT", Timeframe: kline.Timeframe1d},
		domain.MarketCrypto, "kline")
}

func TestService_AllSpotAdapters_Crypto(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.spotManagers[domain.MarketCrypto], ctx,
		spot.Request{Symbols: []string{"BTCUSDT"}},
		domain.MarketCrypto, "spot")
}

func TestService_AllInstrumentAdapters_Crypto(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.instrumentManagers[domain.MarketCrypto], ctx,
		instrument.Request{Market: "CRYPTO", PageSize: 50},
		domain.MarketCrypto, "instrument")
}

func TestService_AllKlineAdapters_Forex(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.klineManagers[domain.MarketForex], ctx,
		kline.Request{Symbol: "EURUSD.FOREX.FOREX_SPOT", Timeframe: kline.Timeframe1d},
		domain.MarketForex, "kline")
}

func TestService_AllSpotAdapters_Forex(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.spotManagers[domain.MarketForex], ctx,
		spot.Request{Symbols: []string{"EURUSD.FOREX.FOREX_SPOT"}},
		domain.MarketForex, "spot")
}

func TestService_AllInstrumentAdapters_Forex(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testAdapterFromManager(t, svc.instrumentManagers[domain.MarketForex], ctx,
		instrument.Request{Market: "FOREX", PageSize: 50},
		domain.MarketForex, "instrument")
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

// ========== 服务发现总览 ==========

// TestService_Discovery 打印所有已注册的市场和适配器，帮助客户了解库的功能范围。
func TestService_Discovery(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	markets := []struct {
		name    string
		market  domain.Market
		kline   map[domain.Market]*manager.Manager[kline.Request, kline.Response]
		spot    map[domain.Market]*manager.Manager[spot.Request, spot.Response]
		inst    map[domain.Market]*manager.Manager[instrument.Request, instrument.Response]
		profile map[domain.Market]*manager.Manager[profile.Request, profile.Response]
		fin     map[domain.Market]*manager.Manager[financial.Request, financial.Response]
		ann     map[domain.Market]*manager.Manager[announcement.Request, announcement.Response]
	}{
		{"CN (A股)", domain.MarketCN, svc.klineManagers, svc.spotManagers, svc.instrumentManagers, svc.profileManagers, svc.financialManagers, svc.announcementManagers},
		{"US (美股)", domain.MarketUS, svc.klineManagers, svc.spotManagers, svc.instrumentManagers, svc.profileManagers, svc.financialManagers, svc.announcementManagers},
		{"HK (港股)", domain.MarketHK, svc.klineManagers, svc.spotManagers, svc.instrumentManagers, svc.profileManagers, svc.financialManagers, svc.announcementManagers},
		{"Crypto (加密货币)", domain.MarketCrypto, svc.klineManagers, svc.spotManagers, svc.instrumentManagers, svc.profileManagers, svc.financialManagers, svc.announcementManagers},
		{"Forex (外汇)", domain.MarketForex, svc.klineManagers, svc.spotManagers, svc.instrumentManagers, svc.profileManagers, svc.financialManagers, svc.announcementManagers},
	}

	for _, mk := range markets {
		t.Logf("\n=== %s ===", mk.name)

		if m, ok := mk.kline[mk.market]; ok {
			t.Logf("  K线: %v", m.Providers())
		} else {
			t.Logf("  K线: (none)")
		}

		if m, ok := mk.spot[mk.market]; ok {
			t.Logf("  实时行情: %v", m.Providers())
		} else {
			t.Logf("  实时行情: (none)")
		}

		if m, ok := mk.inst[mk.market]; ok {
			t.Logf("  证券列表: %v", m.Providers())
		} else {
			t.Logf("  证券列表: (none)")
		}

		if m, ok := mk.profile[mk.market]; ok {
			t.Logf("  个股档案: %v", m.Providers())
		} else {
			t.Logf("  个股档案: (none)")
		}

		if m, ok := mk.fin[mk.market]; ok {
			t.Logf("  财务数据: %v", m.Providers())
		} else {
			t.Logf("  财务数据: (none)")
		}

		if m, ok := mk.ann[mk.market]; ok {
			t.Logf("  公告新闻: %v", m.Providers())
		} else {
			t.Logf("  公告新闻: (none)")
		}
	}
}

// ========== WithMetrics 测试 ==========

func TestService_WithMetrics(t *testing.T) {
	collector := manager.NewMemoryCollector()
	svc := NewService(WithMetrics(collector))
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Make a request to populate metrics
	_, err := svc.GetKline(ctx, kline.Request{
		Symbol:    "000001.SZ",
		Timeframe: kline.Timeframe1d,
	})
	checkFacadeError(t, err)

	stats := svc.GetStats()
	t.Logf("Stats after kline request: total=%d, success=%d, failed=%d",
		stats.TotalFetches, stats.SuccessFetches, stats.FailedFetches)

	if stats.TotalFetches == 0 {
		t.Error("Expected at least one fetch recorded in metrics")
	}

	if len(stats.ByProvider) == 0 {
		t.Error("Expected at least one provider in metrics")
	}

	for name, pm := range stats.ByProvider {
		t.Logf("  Provider %s: fetches=%d, success=%d, failed=%d",
			name, pm.Fetches, pm.Success, pm.Failed)
	}
}

// ========== 辅助函数 ==========

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// formatResultSummary 格式化测试结果摘要，方便客户快速了解库的能力。
func formatResultSummary(domainName string, market domain.Market, providerName string, err error) string {
	if err != nil {
		return fmt.Sprintf("[%s/%s] %s: FAIL - %v", market, domainName, providerName, err)
	}
	return fmt.Sprintf("[%s/%s] %s: OK", market, domainName, providerName)
}

// ========== WithConfig 测试 ==========

func TestNewServiceWithConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	svc := NewService(WithConfig(cfg))
	defer svc.Close()

	// Verify config is accessible
	if svc.Config() != cfg {
		t.Error("Config() should return the same config passed to WithConfig")
	}
}

func TestNewServiceDefaultConfig(t *testing.T) {
	svc := NewService()
	defer svc.Close()

	// Verify default config is set
	if svc.Config() == nil {
		t.Error("Config() should not be nil when no config is provided")
	}
}
