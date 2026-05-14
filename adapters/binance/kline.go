package binance

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/binance"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
)

// Adapter name
const Name = "binance"

// Supported markets for Binance adapter
var supportedMarkets = []domain.Market{domain.MarketCrypto}

// KlineAdapter adapts Binance K-line data
type KlineAdapter struct {
	client *binance.Client
}

// NewKlineAdapter creates a new K-line adapter
func NewKlineAdapter(client *binance.Client) *KlineAdapter {
	return &KlineAdapter{client: client}
}

// Name returns the adapter name
func (a *KlineAdapter) Name() string {
	return Name
}

// SupportedMarkets returns supported markets
func (a *KlineAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

// CanHandle checks if the adapter can handle the symbol
func (a *KlineAdapter) CanHandle(symbol string) bool {
	var sym domain.Symbol
	if err := sym.Parse(symbol); err != nil {
		// Also try direct crypto format
		return binance.IsCryptoSymbol(symbol)
	}
	for _, m := range supportedMarkets {
		if sym.Market == m {
			return true
		}
	}
	return false
}

// Fetch retrieves K-line data
func (a *KlineAdapter) Fetch(ctx context.Context, req kline.Request) (kline.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	// Convert symbol to Binance format
	symbol, err := binance.ToBinanceSymbol(req.Symbol)
	if err != nil {
		return kline.Response{}, trace, err
	}

	// Calculate limit based on date range
	limit := 500
	if !req.StartTime.IsZero() && !req.EndTime.IsZero() {
		days := int(req.EndTime.Sub(req.StartTime).Hours() / 24)
		if days > 0 && days < 1000 {
			limit = days + 1
		}
	}

	params := &binance.KlineParams{
		Symbol:    symbol,
		Interval:  binance.ToInterval(string(req.Timeframe)),
		Limit:     limit,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	result, record, err := a.client.GetKline(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Data))
	for _, d := range result.Data {
		bars = append(bars, kline.Bar{
			Timestamp: binance.ParseOpenTime(d.OpenTime),
			Open:      d.Open,
			High:      d.High,
			Low:       d.Low,
			Close:     d.Close,
			Volume:    d.Volume,
			Turnover:  d.QuoteVol,
		})
	}

	trace.Finish()
	return kline.Response{
		Symbol:      req.Symbol,
		Bars:        bars,
		Source:      Name,
		DataVersion: 1,
	}, trace, nil
}

// inferIntervalFromTimeframe infers K-line interval from timeframe
func inferIntervalFromTimeframe(tf kline.Timeframe) time.Duration {
	switch tf {
	case kline.Timeframe1m:
		return time.Minute
	case kline.Timeframe5m:
		return 5 * time.Minute
	case kline.Timeframe15m:
		return 15 * time.Minute
	case kline.Timeframe30m:
		return 30 * time.Minute
	case kline.Timeframe60m, kline.Timeframe1H:
		return time.Hour
	case kline.Timeframe4H:
		return 4 * time.Hour
	case kline.Timeframe1d, "":
		return 24 * time.Hour
	case kline.Timeframe1w:
		return 7 * 24 * time.Hour
	case kline.Timeframe1M:
		return 30 * 24 * time.Hour
	default:
		return 24 * time.Hour
	}
}

var _ manager.Provider[kline.Request, kline.Response] = (*KlineAdapter)(nil)

// FuturesName is the adapter name for Binance Futures
const FuturesName = "binance_futures"

// Supported markets for Binance Futures adapter
var futuresSupportedMarkets = []domain.Market{domain.MarketFutures}

// FuturesKlineAdapter adapts Binance Futures K-line data
type FuturesKlineAdapter struct {
	client *binance.Client
}

// NewFuturesKlineAdapter creates a new K-line adapter for Binance Futures
func NewFuturesKlineAdapter(client *binance.Client) *FuturesKlineAdapter {
	return &FuturesKlineAdapter{client: client}
}

// Name returns the adapter name
func (a *FuturesKlineAdapter) Name() string {
	return FuturesName
}

// SupportedMarkets returns supported markets
func (a *FuturesKlineAdapter) SupportedMarkets() []domain.Market {
	return futuresSupportedMarkets
}

// CanHandle checks if the adapter can handle the symbol
func (a *FuturesKlineAdapter) CanHandle(symbol string) bool {
	var sym domain.Symbol
	if err := sym.Parse(symbol); err != nil {
		return binance.IsCryptoSymbol(symbol)
	}
	for _, m := range futuresSupportedMarkets {
		if sym.Market == m {
			return true
		}
	}
	return false
}

// Fetch retrieves K-line data from Binance Futures
func (a *FuturesKlineAdapter) Fetch(ctx context.Context, req kline.Request) (kline.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(FuturesName)

	// Convert symbol to Binance format
	symbol, err := binance.ToBinanceSymbol(req.Symbol)
	if err != nil {
		return kline.Response{}, trace, err
	}

	// Calculate limit based on date range
	limit := 500
	if !req.StartTime.IsZero() && !req.EndTime.IsZero() {
		days := int(req.EndTime.Sub(req.StartTime).Hours() / 24)
		if days > 0 && days < 1000 {
			limit = days + 1
		}
	}

	params := &binance.KlineParams{
		Symbol:    symbol,
		Interval:  binance.ToInterval(string(req.Timeframe)),
		Limit:     limit,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	result, record, err := a.client.GetKline(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Data))
	for _, d := range result.Data {
		bars = append(bars, kline.Bar{
			Timestamp: binance.ParseOpenTime(d.OpenTime),
			Open:      d.Open,
			High:      d.High,
			Low:       d.Low,
			Close:     d.Close,
			Volume:    d.Volume,
			Turnover:  d.QuoteVol,
		})
	}

	trace.Finish()
	return kline.Response{
		Symbol:      req.Symbol,
		Bars:        bars,
		Source:      FuturesName,
		DataVersion: 1,
	}, trace, nil
}

var _ manager.Provider[kline.Request, kline.Response] = (*FuturesKlineAdapter)(nil)
