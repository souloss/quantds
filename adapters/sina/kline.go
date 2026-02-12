package sina

import (
	"context"

	"github.com/souloss/quantds/clients/sina"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "sina"

// supportedMarkets 定义 Sina 适配器支持的市场
var supportedMarkets = []domain.Market{domain.MarketCN}

type KlineAdapter struct {
	client *sina.Client
}

func NewKlineAdapter(client *sina.Client) *KlineAdapter {
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

	params := &sina.KlineParams{
		Symbol: req.Symbol,
		Period: sina.ToPeriod(string(req.Timeframe)),
	}

	result, record, err := a.client.GetKline(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Data))
	for _, bar := range result.Data {
		ts := sina.ParseDate(bar.Date)
		bars = append(bars, kline.Bar{
			Timestamp: ts,
			Open:      bar.Open,
			High:      bar.High,
			Low:       bar.Low,
			Close:     bar.Close,
			Volume:    bar.Volume,
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
