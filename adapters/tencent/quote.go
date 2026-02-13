package tencent

import (
	"context"

	"github.com/souloss/quantds/clients/tencent"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// QuoteAdapter adapts Tencent real-time quote data
// Note: Quote is an alias for Spot in this domain
type QuoteAdapter struct {
	client *tencent.Client
}

// NewQuoteAdapter creates a new quote adapter
func NewQuoteAdapter(client *tencent.Client) *QuoteAdapter {
	return &QuoteAdapter{client: client}
}

// Name returns the adapter name
func (a *QuoteAdapter) Name() string {
	return Name
}

// SupportedMarkets returns supported markets
func (a *QuoteAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

// CanHandle checks if the adapter can handle the symbol
func (a *QuoteAdapter) CanHandle(symbol string) bool {
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

// Fetch retrieves real-time quote data
func (a *QuoteAdapter) Fetch(ctx context.Context, _ request.Client, req spot.Request) (spot.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	if len(req.Symbols) == 0 {
		return spot.Response{}, trace, nil
	}

	params := &tencent.QuoteParams{
		Symbols: req.Symbols,
	}

	result, record, err := a.client.GetQuotes(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return spot.Response{}, trace, err
	}

	quotes := make([]spot.Quote, 0, len(result.Data))
	for _, q := range result.Data {
		amplitude := 0.0
		if q.PreClose > 0 {
			amplitude = (q.High - q.Low) / q.PreClose * 100
		}

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
			Amplitude:  amplitude,
			Timestamp:  parseTencentTime(q.Time),
		})
	}

	trace.Finish()
	return spot.Response{
		Quotes: quotes,
		Total:  len(quotes),
		Source: Name,
	}, trace, nil
}

var _ manager.Provider[spot.Request, spot.Response] = (*QuoteAdapter)(nil)
