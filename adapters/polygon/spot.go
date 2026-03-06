package polygon

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/polygon"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

type SpotAdapter struct {
	client *polygon.Client
}

func NewSpotAdapter(client *polygon.Client) *SpotAdapter {
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

	tickers := make([]string, 0, len(req.Symbols))
	for _, s := range req.Symbols {
		var sym domain.Symbol
		if err := sym.Parse(s); err != nil {
			continue
		}
		tickers = append(tickers, sym.Code)
	}

	if len(tickers) == 0 {
		trace.Finish()
		return spot.Response{Source: Name}, trace, nil
	}

	result, record, err := a.client.GetSnapshot(ctx, &polygon.SnapshotParams{Tickers: tickers})
	trace.AddRequest(record)

	if err != nil {
		return spot.Response{}, trace, err
	}

	quotes := make([]spot.Quote, 0, len(result.Tickers))
	for i, t := range result.Tickers {
		symbol := req.Symbols[0]
		if i < len(req.Symbols) {
			symbol = req.Symbols[i]
		}
		quotes = append(quotes, spot.Quote{
			Symbol:     symbol,
			Latest:     t.Day.Close,
			Open:       t.Day.Open,
			High:       t.Day.High,
			Low:        t.Day.Low,
			PreClose:   t.PrevDay.Close,
			Change:     t.Change,
			ChangeRate: t.ChangePercent,
			Volume:     t.Day.Volume,
			Timestamp:  time.UnixMilli(t.Updated),
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
