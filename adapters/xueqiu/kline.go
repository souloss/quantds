package xueqiu

import (
	"context"

	"github.com/souloss/quantds/clients/xueqiu"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "xueqiu"

var supportedMarkets = []domain.Market{domain.MarketCN}

// KlineAdapter adapts Xueqiu kline data
type KlineAdapter struct {
	client *xueqiu.Client
}

// NewKlineAdapter creates a new kline adapter
func NewKlineAdapter(client *xueqiu.Client) *KlineAdapter {
	return &KlineAdapter{client: client}
}

func (a *KlineAdapter) Name() string                      { return Name }
func (a *KlineAdapter) SupportedMarkets() []domain.Market { return supportedMarkets }

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

func (a *KlineAdapter) Fetch(ctx context.Context, _ request.Client, req kline.Request) (kline.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	params := &xueqiu.KlineParams{
		Symbol: req.Symbol,
		Period: xueqiu.ToPeriod(string(req.Timeframe)),
		Count:  500,
	}

	result, record, err := a.client.GetKline(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Data))
	for _, bar := range result.Data {
		ts := xueqiu.ParseTimestamp(bar.Timestamp)
		bars = append(bars, kline.Bar{
			Timestamp:  ts,
			Open:       bar.Open,
			High:       bar.High,
			Low:        bar.Low,
			Close:      bar.Close,
			Volume:     bar.Volume,
			Turnover:   bar.Turnover,
			Change:     bar.Change,
			ChangeRate: bar.ChangePercent,
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
