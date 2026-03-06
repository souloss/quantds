package eodhd

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/eodhd"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "eodhd"

var supportedMarkets = []domain.Market{domain.MarketUS}

type KlineAdapter struct {
	client *eodhd.Client
}

func NewKlineAdapter(client *eodhd.Client) *KlineAdapter {
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

	params := &eodhd.EODParams{
		Symbol: sym.Code + ".US",
		Period: "d",
	}

	if !req.StartTime.IsZero() {
		params.From = req.StartTime.Format("2006-01-02")
	}
	if !req.EndTime.IsZero() {
		params.To = req.EndTime.Format("2006-01-02")
	}

	result, record, err := a.client.GetEOD(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Data))
	for _, d := range result.Data {
		ts, parseErr := time.Parse("2006-01-02", d.Date)
		if parseErr != nil {
			continue
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
