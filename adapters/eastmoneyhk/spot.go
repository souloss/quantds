package eastmoneyhk

import (
	"context"
	"fmt"

	"github.com/souloss/quantds/clients/eastmoneyhk"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// SpotAdapter adapts EastMoney HK real-time quote data
type SpotAdapter struct {
	client *eastmoneyhk.Client
}

// NewSpotAdapter creates a new spot adapter
func NewSpotAdapter(client *eastmoneyhk.Client) *SpotAdapter {
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

	// Get quotes for specific symbols or all HK stocks
	var result *eastmoneyhk.QuoteResult
	var record *request.Record
	var err error

	if len(req.Symbols) > 0 {
		result, record, err = a.client.GetQuotesBySymbols(ctx, req.Symbols)
	} else {
		result, record, err = a.client.GetQuote(ctx, &eastmoneyhk.QuoteParams{
			PageSize: 500,
		})
	}
	trace.AddRequest(record)

	if err != nil {
		return spot.Response{}, trace, err
	}

	quotes := make([]spot.Quote, 0, len(result.Quotes))
	for _, q := range result.Quotes {
		quotes = append(quotes, spot.Quote{
			Symbol:     formatHKSymbol(q.Code),
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
		})
	}

	trace.Finish()
	return spot.Response{
		Quotes: quotes,
		Total:  len(quotes),
		Source: Name,
	}, trace, nil
}

// formatHKSymbol formats HK stock code to standard format
func formatHKSymbol(code string) string {
	return fmt.Sprintf("%s.HK.HKEX", code)
}

var _ manager.Provider[spot.Request, spot.Response] = (*SpotAdapter)(nil)
