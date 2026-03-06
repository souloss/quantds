package eodhd

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/eodhd"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

type SpotAdapter struct {
	client *eodhd.Client
}

func NewSpotAdapter(client *eodhd.Client) *SpotAdapter {
	return &SpotAdapter{client: client}
}

func (a *SpotAdapter) Name() string                      { return Name }
func (a *SpotAdapter) SupportedMarkets() []domain.Market { return supportedMarkets }

func (a *SpotAdapter) CanHandle(symbol string) bool {
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

func (a *SpotAdapter) Fetch(ctx context.Context, _ request.Client, req spot.Request) (spot.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	quotes := make([]spot.Quote, 0, len(req.Symbols))
	for _, s := range req.Symbols {
		var sym domain.Symbol
		if err := sym.Parse(s); err != nil {
			continue
		}

		result, record, err := a.client.GetRealTimeQuote(ctx, &eodhd.RealTimeParams{
			Symbol: sym.Code + ".US",
		})
		trace.AddRequest(record)
		if err != nil {
			continue
		}

		quotes = append(quotes, spot.Quote{
			Symbol:     s,
			Latest:     result.Close,
			Open:       result.Open,
			High:       result.High,
			Low:        result.Low,
			PreClose:   result.PreviousClose,
			Change:     result.Change,
			ChangeRate: result.ChangePercent,
			Volume:     result.Volume,
			Timestamp:  time.Unix(result.Timestamp, 0),
		})
	}

	trace.Finish()
	return spot.Response{
		Quotes: quotes,
		Total:  len(quotes),
		Source: Name,
	}, trace, nil
}

var _ manager.Provider[spot.Request, spot.Response] = (*SpotAdapter)(nil)
