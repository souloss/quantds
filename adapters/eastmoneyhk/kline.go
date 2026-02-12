package eastmoneyhk

import (
	"context"

	"github.com/souloss/quantds/clients/eastmoneyhk"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// Adapter name
const Name = "eastmoneyhk"

// Supported markets for EastMoney HK adapter
var supportedMarkets = []domain.Market{domain.MarketHK}

// KlineAdapter adapts EastMoney HK K-line data
type KlineAdapter struct {
	client *eastmoneyhk.Client
}

// NewKlineAdapter creates a new K-line adapter
func NewKlineAdapter(client *eastmoneyhk.Client) *KlineAdapter {
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

	// Parse HK symbol
	code, ok := eastmoneyhk.ParseHKSymbol(req.Symbol)
	if !ok {
		return kline.Response{}, trace, nil
	}

	params := &eastmoneyhk.KlineParams{
		Symbol:    code + ".HK",
		StartDate: req.StartTime.Format("20060102"),
		EndDate:   req.EndTime.Format("20060102"),
		Period:    eastmoneyhk.ToPeriod(string(req.Timeframe)),
		Adjust:    eastmoneyhk.ToAdjust(string(req.Adjust)),
	}

	result, record, err := a.client.GetKline(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Data))
	for _, d := range result.Data {
		bars = append(bars, kline.Bar{
			Timestamp:    d.Timestamp,
			Open:         d.Open,
			High:         d.High,
			Low:          d.Low,
			Close:        d.Close,
			Volume:       d.Volume,
			Turnover:     d.Turnover,
			Change:       d.Change,
			ChangeRate:   d.ChangeRate,
			TurnoverRate: d.TurnoverRate,
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
