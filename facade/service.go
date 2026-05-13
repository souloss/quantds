package facade

import (
	"context"
	"fmt"
	"strings"
	"time"

	alphavantageadapter "github.com/souloss/quantds/adapters/alphavantage"
	binanceadapter "github.com/souloss/quantds/adapters/binance"
	bseadapter "github.com/souloss/quantds/adapters/bse"
	cninfoadapter "github.com/souloss/quantds/adapters/cninfo"
	coingeckoadapter "github.com/souloss/quantds/adapters/coingecko"
	eastmoneyadapter "github.com/souloss/quantds/adapters/eastmoney"
	eastmoneyfundadapter "github.com/souloss/quantds/adapters/eastmoneyfund"
	eastmoneyhkadapter "github.com/souloss/quantds/adapters/eastmoneyhk"
	eodhadadapter "github.com/souloss/quantds/adapters/eodhd"
	finnhubadapter "github.com/souloss/quantds/adapters/finnhub"
	okxadapter "github.com/souloss/quantds/adapters/okx"
	polygonadapter "github.com/souloss/quantds/adapters/polygon"
	sinaadapter "github.com/souloss/quantds/adapters/sina"
	sseadapter "github.com/souloss/quantds/adapters/sse"
	szseadapter "github.com/souloss/quantds/adapters/szse"
	tencentadapter "github.com/souloss/quantds/adapters/tencent"
	tushareadapter "github.com/souloss/quantds/adapters/tushare"
	twelvedataadapter "github.com/souloss/quantds/adapters/twelvedata"
	xueqiuadapter "github.com/souloss/quantds/adapters/xueqiu"
	yahooadapter "github.com/souloss/quantds/adapters/yahoo"
	alphavantageclient "github.com/souloss/quantds/clients/alphavantage"
	binanceclient "github.com/souloss/quantds/clients/binance"
	bseclient "github.com/souloss/quantds/clients/bse"
	cninfoclient "github.com/souloss/quantds/clients/cninfo"
	coingeckoclient "github.com/souloss/quantds/clients/coingecko"
	eastmoneyclient "github.com/souloss/quantds/clients/eastmoney"
	eastmoneyfundclient "github.com/souloss/quantds/clients/eastmoneyfund"
	eastmoneyhkclient "github.com/souloss/quantds/clients/eastmoneyhk"
	eodhdclient "github.com/souloss/quantds/clients/eodhd"
	finnhubclient "github.com/souloss/quantds/clients/finnhub"
	okxclient "github.com/souloss/quantds/clients/okx"
	polygonclient "github.com/souloss/quantds/clients/polygon"
	sinaclient "github.com/souloss/quantds/clients/sina"
	sseclient "github.com/souloss/quantds/clients/sse"
	szseclient "github.com/souloss/quantds/clients/szse"
	tencentclient "github.com/souloss/quantds/clients/tencent"
	tushareclient "github.com/souloss/quantds/clients/tushare"
	twelvedataclient "github.com/souloss/quantds/clients/twelvedata"
	xueqiuclient "github.com/souloss/quantds/clients/xueqiu"
	yahooclient "github.com/souloss/quantds/clients/yahoo"
	"github.com/souloss/quantds/config"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/announcement"
	"github.com/souloss/quantds/domain/financial"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/domain/profile"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	mw "github.com/souloss/quantds/manager/middleware"
	"github.com/souloss/quantds/request"
)

const (
	PriorityHighest = 100
	PriorityHigh    = 75
	PriorityMedium  = 50
	PriorityLow     = 25
	PriorityLowest  = 1

	CacheTTLKline = 5 * time.Minute
	CacheTTLSpot  = 10 * time.Second
	CacheTTLList  = 1 * time.Hour
)

// Service 多市场数据服务门面，统一编排各数据源提供商。
type Service struct {
	klineManagers        map[domain.Market]*manager.Manager[kline.Request, kline.Response]
	spotManagers         map[domain.Market]*manager.Manager[spot.Request, spot.Response]
	instrumentManagers   map[domain.Market]*manager.Manager[instrument.Request, instrument.Response]
	profileManagers      map[domain.Market]*manager.Manager[profile.Request, profile.Response]
	financialManagers    map[domain.Market]*manager.Manager[financial.Request, financial.Response]
	announcementManagers map[domain.Market]*manager.Manager[announcement.Request, announcement.Response]

	httpClient request.Client
	metrics    manager.Collector
	cfg        *config.Config
}

// ServiceOption defines the option for Service.
type ServiceOption func(*Service)

// WithMetrics enables metrics collection.
func WithMetrics(collector manager.Collector) ServiceOption {
	return func(s *Service) {
		s.metrics = collector
	}
}

// WithConfig sets the configuration for the service.
func WithConfig(cfg *config.Config) ServiceOption {
	return func(s *Service) {
		s.cfg = cfg
	}
}

// NewService 创建新的多市场数据服务。
func NewService(opts ...ServiceOption) *Service {
	s := &Service{
		httpClient:           request.NewClient(request.DefaultConfig()),
		klineManagers:        make(map[domain.Market]*manager.Manager[kline.Request, kline.Response]),
		spotManagers:         make(map[domain.Market]*manager.Manager[spot.Request, spot.Response]),
		instrumentManagers:   make(map[domain.Market]*manager.Manager[instrument.Request, instrument.Response]),
		profileManagers:      make(map[domain.Market]*manager.Manager[profile.Request, profile.Response]),
		financialManagers:    make(map[domain.Market]*manager.Manager[financial.Request, financial.Response]),
		announcementManagers: make(map[domain.Market]*manager.Manager[announcement.Request, announcement.Response]),
		metrics:              manager.NewNoopCollector(),
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.cfg == nil {
		s.cfg = config.DefaultConfig()
	}
	s.initManagers()
	return s
}

// GetStats returns the metrics statistics.
func (s *Service) GetStats() manager.Stats {
	return s.metrics.GetStats()
}

// providerEnabled checks if a provider is enabled in the config.
// Returns true if the provider is not found in the config (default: enabled).
func (s *Service) providerEnabled(name string) bool {
	pc, err := s.cfg.GetProvider(name)
	if err != nil {
		return true // not configured → enabled by default
	}
	return pc.Enabled
}

// applyMiddlewareFromConfig wraps a provider with rate limiter and circuit breaker
// middleware if configured for the given provider name.
func applyMiddlewareFromConfig[Req, Resp any](s *Service, name string, p manager.Provider[Req, Resp]) manager.Provider[Req, Resp] {
	pc, err := s.cfg.GetProvider(name)
	if err != nil {
		return p
	}
	// Apply rate limiter
	if pc.RateLimit.Enabled {
		p = mw.RateLimiterFromConfig[Req, Resp](mw.RateLimiterConfig{
			Enabled:           pc.RateLimit.Enabled,
			RequestsPerMinute: pc.RateLimit.RequestsPerMinute,
			Burst:             pc.RateLimit.Burst,
			Jitter:            pc.RateLimit.Jitter,
		})(p)
	}
	// Apply circuit breaker
	if pc.CircuitBreaker.Enabled {
		p = mw.CircuitBreakerFromConfig[Req, Resp](mw.CircuitBreakerConfig{
			Enabled:          pc.CircuitBreaker.Enabled,
			FailureThreshold: pc.CircuitBreaker.FailureThreshold,
			SuccessThreshold: pc.CircuitBreaker.SuccessThreshold,
			Timeout:          pc.CircuitBreaker.Timeout,
		})(p)
	}
	return p
}

func (s *Service) initManagers() {
	// ========== K 线数据 ==========
	// A股 (CN) - 支持 eastmoney, sina, tencent, tushare, xueqiu
	s.klineManagers[domain.MarketCN] = manager.NewManager[kline.Request, kline.Response](
		manager.WithTwoLevelCache[kline.Request, kline.Response](time.Minute, CacheTTLKline),
		manager.WithMetrics[kline.Request, kline.Response](s.metrics),
		manager.WithProvider[kline.Request, kline.Response](
			applyMiddlewareFromConfig(s, "eastmoney", eastmoneyadapter.NewKlineAdapter(eastmoneyclient.NewClient(eastmoneyclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[kline.Request, kline.Response](
			sinaadapter.NewKlineAdapter(sinaclient.NewClient(sinaclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[kline.Request, kline.Response](
			tencentadapter.NewKlineAdapter(tencentclient.NewClient(tencentclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityMedium),
		),
		manager.WithProvider[kline.Request, kline.Response](
			applyMiddlewareFromConfig(s, "tushare", tushareadapter.NewKlineAdapter(tushareclient.NewClient(tushareclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[kline.Request, kline.Response](
			applyMiddlewareFromConfig(s, "xueqiu", xueqiuadapter.NewKlineAdapter(xueqiuclient.NewClient(xueqiuclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityLowest),
		),
	)

	// ========== 实时行情 ==========
	// A股 (CN) - 支持 sina, tencent, eastmoney, xueqiu, eastmoneyfund
	s.spotManagers[domain.MarketCN] = manager.NewManager[spot.Request, spot.Response](
		manager.WithTwoLevelCache[spot.Request, spot.Response](time.Minute, CacheTTLSpot),
		manager.WithMetrics[spot.Request, spot.Response](s.metrics),
		manager.WithProvider[spot.Request, spot.Response](
			sinaadapter.NewSpotAdapter(sinaclient.NewClient(sinaclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[spot.Request, spot.Response](
			tencentadapter.NewSpotAdapter(tencentclient.NewClient(tencentclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[spot.Request, spot.Response](
			applyMiddlewareFromConfig(s, "eastmoney", eastmoneyadapter.NewSpotAdapter(eastmoneyclient.NewClient(eastmoneyclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityMedium),
		),
		manager.WithProvider[spot.Request, spot.Response](
			applyMiddlewareFromConfig(s, "xueqiu", xueqiuadapter.NewSpotAdapter(xueqiuclient.NewClient(xueqiuclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[spot.Request, spot.Response](
			eastmoneyfundadapter.NewSpotAdapter(eastmoneyfundclient.NewClient(eastmoneyfundclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityLowest),
		),
	)

	// ========== 证券列表 ==========
	// A股 (CN) - 支持 eastmoney, tushare, cninfo, sse, szse, bse, eastmoneyfund
	s.instrumentManagers[domain.MarketCN] = manager.NewManager[instrument.Request, instrument.Response](
		manager.WithTwoLevelCache[instrument.Request, instrument.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[instrument.Request, instrument.Response](s.metrics),
		manager.WithProvider[instrument.Request, instrument.Response](
			applyMiddlewareFromConfig(s, "eastmoney", eastmoneyadapter.NewInstrumentAdapter(eastmoneyclient.NewClient(eastmoneyclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			applyMiddlewareFromConfig(s, "tushare", tushareadapter.NewInstrumentAdapter(tushareclient.NewClient(tushareclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			applyMiddlewareFromConfig(s, "cninfo", cninfoadapter.NewInstrumentAdapter(cninfoclient.NewClient(cninfoclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityMedium),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			sseadapter.NewInstrumentAdapter(sseclient.NewClient(sseclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			szseadapter.NewInstrumentAdapter(szseclient.NewClient(szseclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			bseadapter.NewInstrumentAdapter(bseclient.NewClient(bseclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityLowest),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			eastmoneyfundadapter.NewInstrumentAdapter(eastmoneyfundclient.NewClient(eastmoneyfundclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityLowest),
		),
	)

	// ========== 个股档案 ==========
	// A股 (CN) - 支持 eastmoney, tushare
	s.profileManagers[domain.MarketCN] = manager.NewManager[profile.Request, profile.Response](
		manager.WithTwoLevelCache[profile.Request, profile.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[profile.Request, profile.Response](s.metrics),
		manager.WithProvider[profile.Request, profile.Response](
			applyMiddlewareFromConfig(s, "eastmoney", eastmoneyadapter.NewProfileAdapter(eastmoneyclient.NewClient(eastmoneyclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[profile.Request, profile.Response](
			applyMiddlewareFromConfig(s, "tushare", tushareadapter.NewProfileAdapter(tushareclient.NewClient(tushareclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityMedium),
		),
	)

	// ========== 财务数据 ==========
	// A股 (CN) - 支持 eastmoney, tushare
	s.financialManagers[domain.MarketCN] = manager.NewManager[financial.Request, financial.Response](
		manager.WithTwoLevelCache[financial.Request, financial.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[financial.Request, financial.Response](s.metrics),
		manager.WithProvider[financial.Request, financial.Response](
			applyMiddlewareFromConfig(s, "eastmoney", eastmoneyadapter.NewFinancialAdapter(eastmoneyclient.NewClient(eastmoneyclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[financial.Request, financial.Response](
			tushareadapter.NewFinancialAdapter(tushareclient.NewClient(tushareclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityMedium),
		),
	)

	// ========== 公告新闻 ==========
	// A股 (CN) - 支持 eastmoney, cninfo
	s.announcementManagers[domain.MarketCN] = manager.NewManager[announcement.Request, announcement.Response](
		manager.WithTwoLevelCache[announcement.Request, announcement.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[announcement.Request, announcement.Response](s.metrics),
		manager.WithProvider[announcement.Request, announcement.Response](
			applyMiddlewareFromConfig(s, "eastmoney", eastmoneyadapter.NewAnnouncementAdapter(eastmoneyclient.NewClient(eastmoneyclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[announcement.Request, announcement.Response](
			applyMiddlewareFromConfig(s, "cninfo", cninfoadapter.NewAnnouncementAdapter(cninfoclient.NewClient(cninfoclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityHigh),
		),
	)

	// ========== 美股 (US) ==========
	// K线 - 支持 yahoo, finnhub, polygon, alphavantage, twelvedata, eodhd
	s.klineManagers[domain.MarketUS] = manager.NewManager[kline.Request, kline.Response](
		manager.WithTwoLevelCache[kline.Request, kline.Response](time.Minute, CacheTTLKline),
		manager.WithMetrics[kline.Request, kline.Response](s.metrics),
		manager.WithProvider[kline.Request, kline.Response](
			yahooadapter.NewKlineAdapter(yahooclient.NewClient(yahooclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[kline.Request, kline.Response](
			finnhubadapter.NewKlineAdapter(finnhubclient.NewClient()),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[kline.Request, kline.Response](
			polygonadapter.NewKlineAdapter(polygonclient.NewClient()),
			manager.WithPriority(PriorityMedium),
		),
		manager.WithProvider[kline.Request, kline.Response](
			alphavantageadapter.NewKlineAdapter(alphavantageclient.NewClient()),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[kline.Request, kline.Response](
			twelvedataadapter.NewKlineAdapter(twelvedataclient.NewClient()),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[kline.Request, kline.Response](
			eodhadadapter.NewKlineAdapter(eodhdclient.NewClient()),
			manager.WithPriority(PriorityLowest),
		),
	)

	// 实时行情 - 支持 yahoo, finnhub, polygon, alphavantage, twelvedata, eodhd
	s.spotManagers[domain.MarketUS] = manager.NewManager[spot.Request, spot.Response](
		manager.WithTwoLevelCache[spot.Request, spot.Response](time.Minute, CacheTTLSpot),
		manager.WithMetrics[spot.Request, spot.Response](s.metrics),
		manager.WithProvider[spot.Request, spot.Response](
			yahooadapter.NewSpotAdapter(yahooclient.NewClient(yahooclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[spot.Request, spot.Response](
			finnhubadapter.NewSpotAdapter(finnhubclient.NewClient()),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[spot.Request, spot.Response](
			polygonadapter.NewSpotAdapter(polygonclient.NewClient()),
			manager.WithPriority(PriorityMedium),
		),
		manager.WithProvider[spot.Request, spot.Response](
			alphavantageadapter.NewSpotAdapter(alphavantageclient.NewClient()),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[spot.Request, spot.Response](
			twelvedataadapter.NewSpotAdapter(twelvedataclient.NewClient()),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[spot.Request, spot.Response](
			eodhadadapter.NewSpotAdapter(eodhdclient.NewClient()),
			manager.WithPriority(PriorityLowest),
		),
	)

	// 证券列表 - 支持 yahoo, finnhub, polygon, twelvedata, eodhd
	s.instrumentManagers[domain.MarketUS] = manager.NewManager[instrument.Request, instrument.Response](
		manager.WithTwoLevelCache[instrument.Request, instrument.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[instrument.Request, instrument.Response](s.metrics),
		manager.WithProvider[instrument.Request, instrument.Response](
			yahooadapter.NewInstrumentAdapter(yahooclient.NewClient(yahooclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			finnhubadapter.NewInstrumentAdapter(finnhubclient.NewClient()),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			polygonadapter.NewInstrumentAdapter(polygonclient.NewClient()),
			manager.WithPriority(PriorityMedium),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			twelvedataadapter.NewInstrumentAdapter(twelvedataclient.NewClient()),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			eodhadadapter.NewInstrumentAdapter(eodhdclient.NewClient()),
			manager.WithPriority(PriorityLowest),
		),
	)

	// ========== 港股 (HK) ==========
	// K线 - 支持 eastmoneyhk
	s.klineManagers[domain.MarketHK] = manager.NewManager[kline.Request, kline.Response](
		manager.WithTwoLevelCache[kline.Request, kline.Response](time.Minute, CacheTTLKline),
		manager.WithMetrics[kline.Request, kline.Response](s.metrics),
		manager.WithProvider[kline.Request, kline.Response](
			applyMiddlewareFromConfig(s, "eastmoneyhk", eastmoneyhkadapter.NewKlineAdapter(eastmoneyhkclient.NewClient(eastmoneyhkclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityHighest),
		),
	)

	// 实时行情 - 支持 eastmoneyhk
	s.spotManagers[domain.MarketHK] = manager.NewManager[spot.Request, spot.Response](
		manager.WithTwoLevelCache[spot.Request, spot.Response](time.Minute, CacheTTLSpot),
		manager.WithMetrics[spot.Request, spot.Response](s.metrics),
		manager.WithProvider[spot.Request, spot.Response](
			applyMiddlewareFromConfig(s, "eastmoneyhk", eastmoneyhkadapter.NewSpotAdapter(eastmoneyhkclient.NewClient(eastmoneyhkclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityHighest),
		),
	)

	// 证券列表 - 支持 eastmoneyhk
	s.instrumentManagers[domain.MarketHK] = manager.NewManager[instrument.Request, instrument.Response](
		manager.WithTwoLevelCache[instrument.Request, instrument.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[instrument.Request, instrument.Response](s.metrics),
		manager.WithProvider[instrument.Request, instrument.Response](
			applyMiddlewareFromConfig(s, "eastmoneyhk", eastmoneyhkadapter.NewInstrumentAdapter(eastmoneyhkclient.NewClient(eastmoneyhkclient.WithHTTPClient(s.httpClient)))),
			manager.WithPriority(PriorityHighest),
		),
	)

	// ========== 加密货币 (Crypto) ==========
	// K线 - 支持 binance, okx, coingecko
	s.klineManagers[domain.MarketCrypto] = manager.NewManager[kline.Request, kline.Response](
		manager.WithTwoLevelCache[kline.Request, kline.Response](time.Minute, CacheTTLKline),
		manager.WithMetrics[kline.Request, kline.Response](s.metrics),
		manager.WithProvider[kline.Request, kline.Response](
			binanceadapter.NewKlineAdapter(binanceclient.NewClient(binanceclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[kline.Request, kline.Response](
			okxadapter.NewKlineAdapter(okxclient.NewClient(okxclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[kline.Request, kline.Response](
			coingeckoadapter.NewKlineAdapter(coingeckoclient.NewClient(coingeckoclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityLowest),
		),
	)

	// 实时行情 - 支持 binance, okx, coingecko
	s.spotManagers[domain.MarketCrypto] = manager.NewManager[spot.Request, spot.Response](
		manager.WithTwoLevelCache[spot.Request, spot.Response](time.Minute, CacheTTLSpot),
		manager.WithMetrics[spot.Request, spot.Response](s.metrics),
		manager.WithProvider[spot.Request, spot.Response](
			binanceadapter.NewSpotAdapter(binanceclient.NewClient(binanceclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[spot.Request, spot.Response](
			okxadapter.NewSpotAdapter(okxclient.NewClient(okxclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[spot.Request, spot.Response](
			coingeckoadapter.NewSpotAdapter(coingeckoclient.NewClient(coingeckoclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityLowest),
		),
	)

	// 证券列表 - 支持 binance, okx, coingecko
	s.instrumentManagers[domain.MarketCrypto] = manager.NewManager[instrument.Request, instrument.Response](
		manager.WithTwoLevelCache[instrument.Request, instrument.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[instrument.Request, instrument.Response](s.metrics),
		manager.WithProvider[instrument.Request, instrument.Response](
			binanceadapter.NewInstrumentAdapter(binanceclient.NewClient(binanceclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			okxadapter.NewInstrumentAdapter(okxclient.NewClient(okxclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			coingeckoadapter.NewInstrumentAdapter(coingeckoclient.NewClient(coingeckoclient.WithHTTPClient(s.httpClient))),
			manager.WithPriority(PriorityLowest),
		),
	)

	// ========== 外汇 (Forex) ==========
	// K线 - 支持 finnhub, alphavantage, twelvedata
	s.klineManagers[domain.MarketForex] = manager.NewManager[kline.Request, kline.Response](
		manager.WithTwoLevelCache[kline.Request, kline.Response](time.Minute, CacheTTLKline),
		manager.WithMetrics[kline.Request, kline.Response](s.metrics),
		manager.WithProvider[kline.Request, kline.Response](
			finnhubadapter.NewKlineAdapter(finnhubclient.NewClient()),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[kline.Request, kline.Response](
			alphavantageadapter.NewKlineAdapter(alphavantageclient.NewClient()),
			manager.WithPriority(PriorityMedium),
		),
		manager.WithProvider[kline.Request, kline.Response](
			twelvedataadapter.NewKlineAdapter(twelvedataclient.NewClient()),
			manager.WithPriority(PriorityLow),
		),
	)

	// 实时行情 - 支持 finnhub, twelvedata
	s.spotManagers[domain.MarketForex] = manager.NewManager[spot.Request, spot.Response](
		manager.WithTwoLevelCache[spot.Request, spot.Response](time.Minute, CacheTTLSpot),
		manager.WithMetrics[spot.Request, spot.Response](s.metrics),
		manager.WithProvider[spot.Request, spot.Response](
			finnhubadapter.NewSpotAdapter(finnhubclient.NewClient()),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[spot.Request, spot.Response](
			twelvedataadapter.NewSpotAdapter(twelvedataclient.NewClient()),
			manager.WithPriority(PriorityMedium),
		),
	)

	// 证券列表 - 支持 finnhub, twelvedata
	s.instrumentManagers[domain.MarketForex] = manager.NewManager[instrument.Request, instrument.Response](
		manager.WithTwoLevelCache[instrument.Request, instrument.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[instrument.Request, instrument.Response](s.metrics),
		manager.WithProvider[instrument.Request, instrument.Response](
			finnhubadapter.NewInstrumentAdapter(finnhubclient.NewClient()),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			twelvedataadapter.NewInstrumentAdapter(twelvedataclient.NewClient()),
			manager.WithPriority(PriorityMedium),
		),
	)
}

// getMarketFromSymbol 从 symbol 解析市场。
func (s *Service) getMarketFromSymbol(symbol string) (domain.Market, error) {
	var sym domain.Symbol
	if err := sym.Parse(symbol); err != nil {
		return "", fmt.Errorf("invalid symbol %s: %w", symbol, err)
	}
	return sym.Market, nil
}

// GetKline 获取 K 线数据。
func (s *Service) GetKline(ctx context.Context, req kline.Request) (kline.Response, error) {
	resp, _, err := s.GetKlineWithTrace(ctx, req)
	return resp, err
}

// GetKlineWithTrace 获取 K 线数据并返回请求追踪信息。
func (s *Service) GetKlineWithTrace(ctx context.Context, req kline.Request) (kline.Response, *manager.RequestTrace, error) {
	market, err := s.getMarketFromSymbol(req.Symbol)
	if err != nil {
		return kline.Response{}, nil, err
	}
	m, ok := s.klineManagers[market]
	if !ok {
		return kline.Response{}, nil, fmt.Errorf("unsupported market for kline: %s", market)
	}
	result, err := m.Fetch(ctx, req)
	if err != nil {
		return kline.Response{}, nil, err
	}
	result.Data.Normalize()
	return result.Data, result.Trace, nil
}

// GetSpot 获取实时行情。
func (s *Service) GetSpot(ctx context.Context, req spot.Request) (spot.Response, error) {
	resp, _, err := s.GetSpotWithTrace(ctx, req)
	return resp, err
}

// GetSpotWithTrace 获取实时行情并返回请求追踪信息。
// 支持跨市场查询：将 symbols 按市场分组，分别路由到对应 manager，合并结果。
func (s *Service) GetSpotWithTrace(ctx context.Context, req spot.Request) (spot.Response, *manager.RequestTrace, error) {
	if len(req.Symbols) == 0 {
		return spot.Response{}, nil, nil
	}

	// Group symbols by market
	marketSymbols := make(map[domain.Market][]string)
	for _, symbol := range req.Symbols {
		market, err := s.getMarketFromSymbol(symbol)
		if err != nil {
			return spot.Response{}, nil, err
		}
		marketSymbols[market] = append(marketSymbols[market], symbol)
	}

	// Single market: fast path
	if len(marketSymbols) == 1 {
		for market, symbols := range marketSymbols {
			m, ok := s.spotManagers[market]
			if !ok {
				return spot.Response{}, nil, fmt.Errorf("unsupported market for spot: %s", market)
			}
			result, err := m.Fetch(ctx, spot.Request{Symbols: symbols})
			if err != nil {
				return spot.Response{}, nil, err
			}
			result.Data.Normalize()
			return result.Data, result.Trace, nil
		}
	}

	// Multi-market: fetch from each manager and merge
	var allQuotes []spot.Quote
	var firstTrace *manager.RequestTrace
	var firstSource string

	for market, symbols := range marketSymbols {
		m, ok := s.spotManagers[market]
		if !ok {
			continue
		}
		result, err := m.Fetch(ctx, spot.Request{Symbols: symbols})
		if err != nil {
			return spot.Response{}, nil, fmt.Errorf("spot fetch failed for market %s: %w", market, err)
		}
		allQuotes = append(allQuotes, result.Data.Quotes...)
		if firstTrace == nil {
			firstTrace = result.Trace
			firstSource = result.Data.Source
		}
	}

	if len(allQuotes) == 0 {
		return spot.Response{}, nil, fmt.Errorf("no spot data returned for any market")
	}

	return spot.Response{
		Quotes: allQuotes,
		Total:  len(allQuotes),
		Source: firstSource,
	}, firstTrace, nil
}

// GetInstruments 获取证券列表。
func (s *Service) GetInstruments(ctx context.Context, req instrument.Request) (instrument.Response, error) {
	market := domain.Market(domain.MarketCN)

	// Determine market from the Market field using domain.Market type
	if req.Market != "" {
		m := domain.Market(strings.ToUpper(req.Market))
		if _, ok := domain.MarketConfigs[m]; ok {
			market = m
		} else {
			// Fallback: try common aliases
			switch strings.ToUpper(req.Market) {
			case "NASDAQ", "NYSE", "AMEX":
				market = domain.MarketUS
			case "HKEX":
				market = domain.MarketHK
			case "BINANCE", "OKX", "COINBASE":
				market = domain.MarketCrypto
			case "FX":
				market = domain.MarketForex
			}
		}
	}

	if req.Exchange != "" {
		var sym domain.Symbol
		if err := sym.Parse("000001." + string(req.Exchange)); err == nil {
			market = sym.Market
		}
	}

	m, ok := s.instrumentManagers[market]
	if !ok {
		return instrument.Response{}, fmt.Errorf("unsupported market for instruments: %s", market)
	}
	result, err := m.Fetch(ctx, req)
	if err != nil {
		return instrument.Response{}, err
	}
	return result.Data, nil
}

// GetProfile 获取个股档案。
func (s *Service) GetProfile(ctx context.Context, req profile.Request) (profile.Response, error) {
	market, err := s.getMarketFromSymbol(req.Symbol)
	if err != nil {
		return profile.Response{}, err
	}
	m, ok := s.profileManagers[market]
	if !ok {
		return profile.Response{}, fmt.Errorf("unsupported market for profile: %s", market)
	}
	result, err := m.Fetch(ctx, req)
	if err != nil {
		return profile.Response{}, err
	}
	return result.Data, nil
}

// GetFinancial 获取财务数据。
func (s *Service) GetFinancial(ctx context.Context, req financial.Request) (financial.Response, error) {
	market, err := s.getMarketFromSymbol(req.Symbol)
	if err != nil {
		return financial.Response{}, err
	}
	m, ok := s.financialManagers[market]
	if !ok {
		return financial.Response{}, fmt.Errorf("unsupported market for financial: %s", market)
	}
	result, err := m.Fetch(ctx, req)
	if err != nil {
		return financial.Response{}, err
	}
	return result.Data, nil
}

// GetAnnouncements 获取公告新闻。
func (s *Service) GetAnnouncements(ctx context.Context, req announcement.Request) (announcement.Response, error) {
	market, err := s.getMarketFromSymbol(req.Symbol)
	if err != nil {
		return announcement.Response{}, err
	}
	m, ok := s.announcementManagers[market]
	if !ok {
		return announcement.Response{}, fmt.Errorf("unsupported market for announcements: %s", market)
	}
	result, err := m.Fetch(ctx, req)
	if err != nil {
		return announcement.Response{}, err
	}
	return result.Data, nil
}

// Stats 返回统计信息，聚合所有 manager 的统计。
func (s *Service) Stats() manager.Stats {
	// Use the shared metrics collector which already aggregates across all managers
	return s.metrics.GetStats()
}

// RegisterKlineProvider registers a kline provider for a specific market.
func (s *Service) RegisterKlineProvider(market domain.Market, p manager.Provider[kline.Request, kline.Response], opts ...manager.ProviderOption) {
	if _, ok := s.klineManagers[market]; !ok {
		s.klineManagers[market] = manager.NewManager[kline.Request, kline.Response](
			manager.WithTwoLevelCache[kline.Request, kline.Response](time.Minute, CacheTTLKline),
			manager.WithMetrics[kline.Request, kline.Response](s.metrics),
		)
	}
	s.klineManagers[market].Register(p, opts...)
}

// RegisterSpotProvider registers a spot provider for a specific market.
func (s *Service) RegisterSpotProvider(market domain.Market, p manager.Provider[spot.Request, spot.Response], opts ...manager.ProviderOption) {
	if _, ok := s.spotManagers[market]; !ok {
		s.spotManagers[market] = manager.NewManager[spot.Request, spot.Response](
			manager.WithTwoLevelCache[spot.Request, spot.Response](time.Minute, CacheTTLSpot),
			manager.WithMetrics[spot.Request, spot.Response](s.metrics),
		)
	}
	s.spotManagers[market].Register(p, opts...)
}

// RegisterInstrumentProvider registers an instrument provider for a specific market.
func (s *Service) RegisterInstrumentProvider(market domain.Market, p manager.Provider[instrument.Request, instrument.Response], opts ...manager.ProviderOption) {
	if _, ok := s.instrumentManagers[market]; !ok {
		s.instrumentManagers[market] = manager.NewManager[instrument.Request, instrument.Response](
			manager.WithTwoLevelCache[instrument.Request, instrument.Response](time.Minute, CacheTTLList),
			manager.WithMetrics[instrument.Request, instrument.Response](s.metrics),
		)
	}
	s.instrumentManagers[market].Register(p, opts...)
}

// RegisterProfileProvider registers a profile provider for a specific market.
func (s *Service) RegisterProfileProvider(market domain.Market, p manager.Provider[profile.Request, profile.Response], opts ...manager.ProviderOption) {
	if _, ok := s.profileManagers[market]; !ok {
		s.profileManagers[market] = manager.NewManager[profile.Request, profile.Response](
			manager.WithTwoLevelCache[profile.Request, profile.Response](time.Minute, CacheTTLList),
			manager.WithMetrics[profile.Request, profile.Response](s.metrics),
		)
	}
	s.profileManagers[market].Register(p, opts...)
}

// RegisterFinancialProvider registers a financial provider for a specific market.
func (s *Service) RegisterFinancialProvider(market domain.Market, p manager.Provider[financial.Request, financial.Response], opts ...manager.ProviderOption) {
	if _, ok := s.financialManagers[market]; !ok {
		s.financialManagers[market] = manager.NewManager[financial.Request, financial.Response](
			manager.WithTwoLevelCache[financial.Request, financial.Response](time.Minute, CacheTTLList),
			manager.WithMetrics[financial.Request, financial.Response](s.metrics),
		)
	}
	s.financialManagers[market].Register(p, opts...)
}

// RegisterAnnouncementProvider registers an announcement provider for a specific market.
func (s *Service) RegisterAnnouncementProvider(market domain.Market, p manager.Provider[announcement.Request, announcement.Response], opts ...manager.ProviderOption) {
	if _, ok := s.announcementManagers[market]; !ok {
		s.announcementManagers[market] = manager.NewManager[announcement.Request, announcement.Response](
			manager.WithTwoLevelCache[announcement.Request, announcement.Response](time.Minute, CacheTTLList),
			manager.WithMetrics[announcement.Request, announcement.Response](s.metrics),
		)
	}
	s.announcementManagers[market].Register(p, opts...)
}

// Close 释放资源。
// Config returns the service configuration.
func (s *Service) Config() *config.Config {
	return s.cfg
}

// ProviderStatusEntry holds the status of a single provider.
type ProviderStatusEntry struct {
	Name               string
	Market             domain.Market
	Domain             string // "kline", "spot", "instrument", etc.
	CircuitBreakerState string // "closed", "open", "half-open"
}

// ProviderStatus returns the circuit breaker state for all registered providers.
func (s *Service) ProviderStatus() []ProviderStatusEntry {
	var entries []ProviderStatusEntry

	// Kline providers
	for market, m := range s.klineManagers {
		for _, entry := range m.ProviderEntries() {
			state := "closed"
			if sp, ok := entry.Provider.(mw.CircuitBreakerStateProvider); ok {
				state = sp.CircuitBreakerState().String()
			}
			entries = append(entries, ProviderStatusEntry{
				Name:               entry.Name,
				Market:             market,
				Domain:             "kline",
				CircuitBreakerState: state,
			})
		}
	}

	// Spot providers
	for market, m := range s.spotManagers {
		for _, entry := range m.ProviderEntries() {
			state := "closed"
			if sp, ok := entry.Provider.(mw.CircuitBreakerStateProvider); ok {
				state = sp.CircuitBreakerState().String()
			}
			entries = append(entries, ProviderStatusEntry{
				Name:               entry.Name,
				Market:             market,
				Domain:             "spot",
				CircuitBreakerState: state,
			})
		}
	}

	// Instrument providers
	for market, m := range s.instrumentManagers {
		for _, entry := range m.ProviderEntries() {
			state := "closed"
			if sp, ok := entry.Provider.(mw.CircuitBreakerStateProvider); ok {
				state = sp.CircuitBreakerState().String()
			}
			entries = append(entries, ProviderStatusEntry{
				Name:               entry.Name,
				Market:             market,
				Domain:             "instrument",
				CircuitBreakerState: state,
			})
		}
	}

	// Profile providers
	for market, m := range s.profileManagers {
		for _, entry := range m.ProviderEntries() {
			state := "closed"
			if sp, ok := entry.Provider.(mw.CircuitBreakerStateProvider); ok {
				state = sp.CircuitBreakerState().String()
			}
			entries = append(entries, ProviderStatusEntry{
				Name:               entry.Name,
				Market:             market,
				Domain:             "profile",
				CircuitBreakerState: state,
			})
		}
	}

	// Financial providers
	for market, m := range s.financialManagers {
		for _, entry := range m.ProviderEntries() {
			state := "closed"
			if sp, ok := entry.Provider.(mw.CircuitBreakerStateProvider); ok {
				state = sp.CircuitBreakerState().String()
			}
			entries = append(entries, ProviderStatusEntry{
				Name:               entry.Name,
				Market:             market,
				Domain:             "financial",
				CircuitBreakerState: state,
			})
		}
	}

	// Announcement providers
	for market, m := range s.announcementManagers {
		for _, entry := range m.ProviderEntries() {
			state := "closed"
			if sp, ok := entry.Provider.(mw.CircuitBreakerStateProvider); ok {
				state = sp.CircuitBreakerState().String()
			}
			entries = append(entries, ProviderStatusEntry{
				Name:               entry.Name,
				Market:             market,
				Domain:             "announcement",
				CircuitBreakerState: state,
			})
		}
	}

	return entries
}

func (s *Service) Close() {
	// Close all managers
	for _, m := range s.klineManagers {
		m.Close()
	}
	for _, m := range s.spotManagers {
		m.Close()
	}
	for _, m := range s.instrumentManagers {
		m.Close()
	}
	for _, m := range s.profileManagers {
		m.Close()
	}
	for _, m := range s.financialManagers {
		m.Close()
	}
	for _, m := range s.announcementManagers {
		m.Close()
	}
	if s.httpClient != nil {
		s.httpClient.Close()
	}
}
