package twelvedata

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/twelvedata"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "twelvedata"

var supportedMarkets = []domain.Market{domain.MarketUS, domain.MarketForex, domain.MarketCrypto}

type KlineAdapter struct {
	client *twelvedata.Client
}

func NewKlineAdapter(client *twelvedata.Client) *KlineAdapter {
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

	params := &twelvedata.TimeSeriesParams{
		Symbol:     sym.Code,
		Interval:   twelvedata.ToInterval(string(req.Timeframe)),
		OutputSize: 100,
	}

	if !req.StartTime.IsZero() {
		params.StartDate = req.StartTime.Format("2006-01-02")
	}
	if !req.EndTime.IsZero() {
		params.EndDate = req.EndTime.Format("2006-01-02")
	}

	result, record, err := a.client.GetTimeSeries(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Data))
	for _, d := range result.Data {
		ts, parseErr := time.Parse("2006-01-02", d.Datetime)
		if parseErr != nil {
			ts, parseErr = time.Parse("2006-01-02 15:04:05", d.Datetime)
			if parseErr != nil {
				continue
			}
		}
		bars = append(bars, kline.Bar{
			Timestamp: ts,
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
