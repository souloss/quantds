package tushare

import (
	"context"

	"github.com/souloss/quantds/clients/tushare"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "tushare"

// supportedMarkets 定义 Tushare 适配器支持的市场
var supportedMarkets = []domain.Market{domain.MarketCN}

type KlineAdapter struct {
	client *tushare.Client
}

func NewKlineAdapter(client *tushare.Client) *KlineAdapter {
	return &KlineAdapter{client: client}
}

func (a *KlineAdapter) Name() string {
	return Name
}

func (a *KlineAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

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

	params := &tushare.KlineParams{
		Symbol:    req.Symbol,
		StartDate: req.StartTime.Format("20060102"),
		EndDate:   req.EndTime.Format("20060102"),
		Period:    tushare.ToPeriod(string(req.Timeframe)),
	}

	result, record, err := a.client.GetKline(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Data))
	for _, bar := range result.Data {
		bars = append(bars, kline.Bar{
			Timestamp:  tushare.ParseDate(bar.Date),
			Open:       bar.Open,
			High:       bar.High,
			Low:        bar.Low,
			Close:      bar.Close,
			Volume:     bar.Volume,
			Turnover:   bar.Amount,
			Change:     bar.Change,
			ChangeRate: bar.ChangeRate,
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
