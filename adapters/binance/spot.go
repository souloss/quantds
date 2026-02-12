package binance

import (
	"context"

	"github.com/souloss/quantds/clients/binance"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// SpotAdapter adapts Binance real-time ticker data
type SpotAdapter struct {
	client *binance.Client
}

// NewSpotAdapter creates a new spot adapter
func NewSpotAdapter(client *binance.Client) *SpotAdapter {
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
		return binance.IsCryptoSymbol(symbol)
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

	// Convert symbols to Binance format
	symbols := make([]string, 0, len(req.Symbols))
	for _, s := range req.Symbols {
		symbol, err := binance.ToBinanceSymbol(s)
		if err != nil {
			continue
		}
		symbols = append(symbols, symbol)
	}

	quotes := make([]spot.Quote, 0)

	// Get ticker data for each symbol
	for _, symbol := range symbols {
		result, record, err := a.client.GetTicker24hr(ctx, &binance.TickerParams{
			Symbol: symbol,
		})
		trace.AddRequest(record)

		if err != nil {
			continue
		}

		for _, t := range result.Tickers {
			quotes = append(quotes, spot.Quote{
				Symbol:     binance.FromBinanceSymbol(t.Symbol),
				Name:       t.Symbol,
				Latest:     t.LastPrice,
				Open:       t.OpenPrice,
				High:       t.HighPrice,
				Low:        t.LowPrice,
				PreClose:   t.PrevClosePrice,
				Change:     t.PriceChange,
				ChangeRate: t.PriceChangePercent,
				Volume:     t.Volume,
				Turnover:   t.QuoteVolume,
				Timestamp:  t.CloseTime,
				BidPrice:   t.BidPrice,
				BidVolume:  t.BidQty,
				AskPrice:   t.AskPrice,
				AskVolume:  t.AskQty,
			})
		}
	}

	trace.Finish()
	return spot.Response{
		Quotes: quotes,
		Total:  len(quotes),
		Source: Name,
	}, trace, nil
}

var _ manager.Provider[spot.Request, spot.Response] = (*SpotAdapter)(nil)
