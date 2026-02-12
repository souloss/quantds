package yahoo

import (
	"context"

	"github.com/souloss/quantds/clients/yahoo"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// SpotAdapter adapts Yahoo Finance real-time quote data
type SpotAdapter struct {
	client *yahoo.Client
}

// NewSpotAdapter creates a new spot adapter
func NewSpotAdapter(client *yahoo.Client) *SpotAdapter {
	return &SpotAdapter{client: client}
}

// Name returns the adapter name
func (a *SpotAdapter) Name() string {
	return Name
}

// SupportedMarkets returns supported markets
func (a *SpotAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

// CanHandle checks if the adapter can handle the symbol
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

// Fetch retrieves real-time quotes
func (a *SpotAdapter) Fetch(ctx context.Context, _ request.Client, req spot.Request) (spot.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	// Convert symbols to Yahoo format
	symbols := make([]string, 0, len(req.Symbols))
	for _, s := range req.Symbols {
		symbol, err := yahoo.ToYahooSymbol(s)
		if err != nil {
			continue
		}
		symbols = append(symbols, symbol)
	}

	if len(symbols) == 0 {
		trace.Finish()
		return spot.Response{Source: Name}, trace, nil
	}

	params := &yahoo.QuoteParams{
		Symbols: symbols,
	}

	result, record, err := a.client.GetQuote(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return spot.Response{}, trace, err
	}

	quotes := make([]spot.Quote, 0, len(result.Quotes))
	for _, q := range result.Quotes {
		quotes = append(quotes, spot.Quote{
			Symbol:     yahoo.FromYahooSymbol(q.Symbol, q.Exchange),
			Name:       q.Name,
			Latest:     q.Latest,
			Open:       q.Open,
			High:       q.High,
			Low:        q.Low,
			PreClose:   q.PreClose,
			Change:     q.Change,
			ChangeRate: q.ChangeRate,
			Volume:     q.Volume,
			Timestamp:  q.Timestamp,
			BidPrice:   q.BidPrice,
			BidVolume:  q.BidSize,
			AskPrice:   q.AskPrice,
			AskVolume:  q.AskSize,
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
