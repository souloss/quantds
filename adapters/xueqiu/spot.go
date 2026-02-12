package xueqiu

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/xueqiu"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// SpotAdapter adapts Xueqiu real-time spot/quote data
type SpotAdapter struct {
	client *xueqiu.Client
}

// NewSpotAdapter creates a new spot adapter
func NewSpotAdapter(client *xueqiu.Client) *SpotAdapter {
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

	params := &xueqiu.SpotParams{
		Symbols: req.Symbols,
	}

	result, record, err := a.client.GetSpot(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return spot.Response{}, trace, err
	}

	quotes := make([]spot.Quote, 0, len(result.Data))
	for _, q := range result.Data {
		quotes = append(quotes, spot.Quote{
			Symbol:     q.Symbol,
			Name:       q.Name,
			Latest:     q.Latest,
			Open:       q.Open,
			High:       q.High,
			Low:        q.Low,
			PreClose:   q.PreClose,
			Change:     q.Change,
			ChangeRate: q.ChangeRate,
			Volume:     q.Volume,
			Turnover:   q.Turnover,
			Timestamp:  time.UnixMilli(q.Timestamp),
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
