package binance

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/binance"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
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
func (a *KlineAdapter) Fetch(ctx context.Context, _ request.Client, req kline.Request) (kline.Response, *manager.RequestTrace, error) {
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
		Symbol: req.Symbol,
		Bars:   bars,
		Source: Name,
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
	case kline.Timeframe60m:
		return time.Hour
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
