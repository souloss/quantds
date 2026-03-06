package polygon

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/polygon"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "polygon"

var supportedMarkets = []domain.Market{domain.MarketUS}

type KlineAdapter struct {
	client *polygon.Client
}

func NewKlineAdapter(client *polygon.Client) *KlineAdapter {
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

	var sym domain.Symbol
	if err := sym.Parse(req.Symbol); err != nil {
		return kline.Response{}, trace, err
	}

	from := req.StartTime
	to := req.EndTime
	if from.IsZero() {
		from = time.Now().AddDate(-1, 0, 0)
	}
	if to.IsZero() {
		to = time.Now()
	}

	timespan, multiplier := polygon.ToTimespan(string(req.Timeframe))

	params := &polygon.AggregateParams{
		Symbol:     sym.Code,
		Multiplier: multiplier,
		Timespan:   timespan,
		From:       from.Format("2006-01-02"),
		To:         to.Format("2006-01-02"),
		Limit:      5000,
	}

	result, record, err := a.client.GetAggregates(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Bars))
	for _, b := range result.Bars {
		bars = append(bars, kline.Bar{
			Timestamp: time.UnixMilli(b.Timestamp),
			Open:      b.Open,
			High:      b.High,
			Low:       b.Low,
			Close:     b.Close,
			Volume:    b.Volume,
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
