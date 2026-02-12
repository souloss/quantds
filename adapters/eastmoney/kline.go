package eastmoney

import (
	"context"

	"github.com/souloss/quantds/clients/eastmoney"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "eastmoney"

// supportedMarkets defines the markets supported by Eastmoney adapters
var supportedMarkets = []domain.Market{domain.MarketCN}

// KlineAdapter adapts Eastmoney kline data
type KlineAdapter struct {
	client *eastmoney.Client
}

// NewKlineAdapter creates a new kline adapter
func NewKlineAdapter(client *eastmoney.Client) *KlineAdapter {
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

// Fetch retrieves kline data
func (a *KlineAdapter) Fetch(ctx context.Context, _ request.Client, req kline.Request) (kline.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	params := &eastmoney.CandleParams{
		Symbol:    req.Symbol,
		StartDate: req.StartTime.Format("20060102"),
		EndDate:   req.EndTime.Format("20060102"),
		Period:    eastmoney.ToPeriod(string(req.Timeframe)),
		Adjust:    eastmoney.ToAdjust(string(req.Adjust)),
	}

	result, record, err := a.client.GetCandles(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Data))
	for _, data := range result.Data {
		bars = append(bars, kline.Bar{
			Timestamp:    eastmoney.ParseDate(data.Date),
			Open:         data.Open,
			High:         data.High,
			Low:          data.Low,
			Close:        data.Close,
			Volume:       data.Volume,
			Turnover:     data.Turnover,
			Change:       data.Change,
			ChangeRate:   data.ChangeRate,
			TurnoverRate: data.TurnoverRate,
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
