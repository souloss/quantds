package facade

import (
	"context"
	"fmt"
	"time"

	bseadapter "github.com/souloss/quantds/adapters/bse"
	cninfoadapter "github.com/souloss/quantds/adapters/cninfo"
	eastmoneyadapter "github.com/souloss/quantds/adapters/eastmoney"
	sinaadapter "github.com/souloss/quantds/adapters/sina"
	sseadapter "github.com/souloss/quantds/adapters/sse"
	szseadapter "github.com/souloss/quantds/adapters/szse"
	tencentadapter "github.com/souloss/quantds/adapters/tencent"
	tushareadapter "github.com/souloss/quantds/adapters/tushare"
	xueqiuadapter "github.com/souloss/quantds/adapters/xueqiu"
	bseclient "github.com/souloss/quantds/clients/bse"
	cninfoclient "github.com/souloss/quantds/clients/cninfo"
	eastmoneyclient "github.com/souloss/quantds/clients/eastmoney"
	sinaclient "github.com/souloss/quantds/clients/sina"
	sseclient "github.com/souloss/quantds/clients/sse"
	szseclient "github.com/souloss/quantds/clients/szse"
	tencentclient "github.com/souloss/quantds/clients/tencent"
	tushareclient "github.com/souloss/quantds/clients/tushare"
	xueqiuclient "github.com/souloss/quantds/clients/xueqiu"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/announcement"
	"github.com/souloss/quantds/domain/financial"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/domain/profile"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
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
}

// ServiceOption defines the option for Service.
type ServiceOption func(*Service)

// WithMetrics enables metrics collection.
func WithMetrics(collector manager.Collector) ServiceOption {
	return func(s *Service) {
		s.metrics = collector
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
	s.initManagers()
	return s
}

// GetStats returns the metrics statistics.
func (s *Service) GetStats() manager.Stats {
	return s.metrics.GetStats()
}

func (s *Service) initManagers() {
	// ========== K 线数据 ==========
	// A股 (CN) - 支持 eastmoney, sina, tencent, tushare, xueqiu
	s.klineManagers[domain.MarketCN] = manager.NewManager[kline.Request, kline.Response](
		manager.WithTwoLevelCache[kline.Request, kline.Response](time.Minute, CacheTTLKline),
		manager.WithMetrics[kline.Request, kline.Response](s.metrics),
		manager.WithProvider[kline.Request, kline.Response](
			eastmoneyadapter.NewKlineAdapter(eastmoneyclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[kline.Request, kline.Response](
			sinaadapter.NewKlineAdapter(sinaclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[kline.Request, kline.Response](
			tencentadapter.NewKlineAdapter(tencentclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityMedium),
		),
		manager.WithProvider[kline.Request, kline.Response](
			tushareadapter.NewKlineAdapter(tushareclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[kline.Request, kline.Response](
			xueqiuadapter.NewKlineAdapter(xueqiuclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityLowest),
		),
	)

	// ========== 实时行情 ==========
	// A股 (CN) - 支持 sina, tencent, eastmoney, xueqiu
	s.spotManagers[domain.MarketCN] = manager.NewManager[spot.Request, spot.Response](
		manager.WithTwoLevelCache[spot.Request, spot.Response](time.Minute, CacheTTLSpot),
		manager.WithMetrics[spot.Request, spot.Response](s.metrics),
		manager.WithProvider[spot.Request, spot.Response](
			sinaadapter.NewSpotAdapter(sinaclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[spot.Request, spot.Response](
			tencentadapter.NewSpotAdapter(tencentclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[spot.Request, spot.Response](
			eastmoneyadapter.NewSpotAdapter(eastmoneyclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityMedium),
		),
		manager.WithProvider[spot.Request, spot.Response](
			xueqiuadapter.NewSpotAdapter(xueqiuclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityLow),
		),
	)

	// ========== 证券列表 ==========
	// A股 (CN) - 支持 eastmoney, tushare, cninfo, sse, szse, bse
	s.instrumentManagers[domain.MarketCN] = manager.NewManager[instrument.Request, instrument.Response](
		manager.WithTwoLevelCache[instrument.Request, instrument.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[instrument.Request, instrument.Response](s.metrics),
		manager.WithProvider[instrument.Request, instrument.Response](
			eastmoneyadapter.NewInstrumentAdapter(eastmoneyclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			tushareadapter.NewInstrumentAdapter(tushareclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityHigh),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			cninfoadapter.NewInstrumentAdapter(cninfoclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityMedium),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			sseadapter.NewInstrumentAdapter(sseclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			szseadapter.NewInstrumentAdapter(szseclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityLow),
		),
		manager.WithProvider[instrument.Request, instrument.Response](
			bseadapter.NewInstrumentAdapter(bseclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityLowest),
		),
	)

	// ========== 个股档案 ==========
	// A股 (CN) - 支持 eastmoney, tushare
	s.profileManagers[domain.MarketCN] = manager.NewManager[profile.Request, profile.Response](
		manager.WithTwoLevelCache[profile.Request, profile.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[profile.Request, profile.Response](s.metrics),
		manager.WithProvider[profile.Request, profile.Response](
			eastmoneyadapter.NewProfileAdapter(eastmoneyclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[profile.Request, profile.Response](
			tushareadapter.NewProfileAdapter(tushareclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityMedium),
		),
	)

	// ========== 财务数据 ==========
	// A股 (CN) - 支持 eastmoney, tushare
	s.financialManagers[domain.MarketCN] = manager.NewManager[financial.Request, financial.Response](
		manager.WithTwoLevelCache[financial.Request, financial.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[financial.Request, financial.Response](s.metrics),
		manager.WithProvider[financial.Request, financial.Response](
			eastmoneyadapter.NewFinancialAdapter(eastmoneyclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[financial.Request, financial.Response](
			tushareadapter.NewFinancialAdapter(tushareclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityMedium),
		),
	)

	// ========== 公告新闻 ==========
	// A股 (CN) - 支持 eastmoney, cninfo
	s.announcementManagers[domain.MarketCN] = manager.NewManager[announcement.Request, announcement.Response](
		manager.WithTwoLevelCache[announcement.Request, announcement.Response](time.Minute, CacheTTLList),
		manager.WithMetrics[announcement.Request, announcement.Response](s.metrics),
		manager.WithProvider[announcement.Request, announcement.Response](
			eastmoneyadapter.NewAnnouncementAdapter(eastmoneyclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityHighest),
		),
		manager.WithProvider[announcement.Request, announcement.Response](
			cninfoadapter.NewAnnouncementAdapter(cninfoclient.NewClient(s.httpClient)),
			manager.WithPriority(PriorityHigh),
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
	return result.Data, result.Trace, nil
}

// GetSpot 获取实时行情。
func (s *Service) GetSpot(ctx context.Context, req spot.Request) (spot.Response, error) {
	resp, _, err := s.GetSpotWithTrace(ctx, req)
	return resp, err
}

// GetSpotWithTrace 获取实时行情并返回请求追踪信息。
func (s *Service) GetSpotWithTrace(ctx context.Context, req spot.Request) (spot.Response, *manager.RequestTrace, error) {
	if len(req.Symbols) == 0 {
		return spot.Response{}, nil, nil
	}
	market, err := s.getMarketFromSymbol(req.Symbols[0])
	if err != nil {
		return spot.Response{}, nil, err
	}
	m, ok := s.spotManagers[market]
	if !ok {
		return spot.Response{}, nil, fmt.Errorf("unsupported market for spot: %s", market)
	}
	result, err := m.Fetch(ctx, req)
	if err != nil {
		return spot.Response{}, nil, err
	}
	return result.Data, result.Trace, nil
}

// GetInstruments 获取证券列表。
func (s *Service) GetInstruments(ctx context.Context, req instrument.Request) (instrument.Response, error) {
	market := domain.MarketCN
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

// Stats 返回统计信息。
func (s *Service) Stats() manager.Stats {
	if m, ok := s.klineManagers[domain.MarketCN]; ok {
		return m.Stats()
	}
	return manager.Stats{}
}

// Close 释放资源。
func (s *Service) Close() {
	if s.httpClient != nil {
		s.httpClient.Close()
	}
}
