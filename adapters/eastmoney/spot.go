package eastmoney

import (
	"context"

	"github.com/souloss/quantds/clients/eastmoney"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// SpotAdapter adapts Eastmoney real-time spot/quote data
type SpotAdapter struct {
	client *eastmoney.Client
}

// NewSpotAdapter creates a new spot adapter
func NewSpotAdapter(client *eastmoney.Client) *SpotAdapter {
	return &SpotAdapter{client: client}
}

func (a *SpotAdapter) Name() string              { return Name }
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

	params := &eastmoney.QuoteParams{
		PageSize: 500,
	}
	if len(req.Symbols) > 0 {
		params.PageSize = len(req.Symbols)
	}

	result, record, err := a.client.GetQuotes(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return spot.Response{}, trace, err
	}

	quotes := make([]spot.Quote, 0, len(result.Data))
	for _, data := range result.Data {
		exchange := "SZ"
		if data.MarketID == 1 {
			exchange = "SH"
		}
		symbol := data.Code + "." + exchange

		quotes = append(quotes, spot.Quote{
			Symbol:       symbol,
			Name:         data.Name,
			Latest:       data.Latest,
			Open:         data.Open,
			High:         data.High,
			Low:          data.Low,
			PreClose:     data.PreClose,
			Change:       data.Change,
			ChangeRate:   data.ChangeRate,
			Volume:       data.Volume,
			Turnover:     data.Turnover,
			Amplitude:    data.Amplitude,
			TurnoverRate: data.TurnoverRate,
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
