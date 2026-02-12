package yahoo

import (
	"context"

	"github.com/souloss/quantds/clients/yahoo"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// Adapter name
const Name = "yahoo"

// Supported markets for Yahoo Finance adapter
var supportedMarkets = []domain.Market{domain.MarketUS}

// KlineAdapter adapts Yahoo Finance K-line data
type KlineAdapter struct {
	client *yahoo.Client
}

// NewKlineAdapter creates a new K-line adapter
func NewKlineAdapter(client *yahoo.Client) *KlineAdapter {
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
		return false
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

	// Convert symbol to Yahoo format
	symbol, err := yahoo.ToYahooSymbol(req.Symbol)
	if err != nil {
		return kline.Response{}, trace, err
	}

	params := &yahoo.KlineParams{
		Symbol:   symbol,
		Interval: yahoo.ToInterval(string(req.Timeframe)),
	}

	// Use range if no explicit dates, otherwise use period
	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		params.Range = yahoo.Range1y
	} else {
		params.StartDate = req.StartTime
		params.EndDate = req.EndTime
	}

	result, record, err := a.client.GetKline(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Data))
	for _, d := range result.Data {
		bars = append(bars, kline.Bar{
			Timestamp: yahoo.ParseTimestamp(d.Timestamp, result.Timezone),
			Open:      d.Open,
			High:      d.High,
			Low:       d.Low,
			Close:     d.Close,
			Volume:    d.Volume,
		})
	}

	trace.Finish()
	return kline.Response{
		Symbol: req.Symbol,
		Bars:   bars,
		Source: Name,
	}, trace, nil
}

var _ manager.Provider[kline.Request, kline.Response] = (*KlineAdapter)(nil)
