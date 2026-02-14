package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/souloss/quantds/clients/okx"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// SpotAdapter adapts OKX ticker data to domain spot
type SpotAdapter struct {
	client *okx.Client
}

// NewSpotAdapter creates a new spot adapter
func NewSpotAdapter(client *okx.Client) *SpotAdapter {
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

// Fetch retrieves real-time quotes from OKX
func (a *SpotAdapter) Fetch(ctx context.Context, _ request.Client, req spot.Request) (spot.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	quotes := make([]spot.Quote, 0, len(req.Symbols))

	for _, symbol := range req.Symbols {
		instID := toOKXInstID(symbol)

		ticker, record, err := a.client.GetTicker(ctx, &okx.TickerRequest{
			InstID: instID,
		})
		trace.AddRequest(record)

		if err != nil {
			continue
		}

		last, _ := strconv.ParseFloat(ticker.Last, 64)
		open, _ := strconv.ParseFloat(ticker.Open24h, 64)
		high, _ := strconv.ParseFloat(ticker.High24h, 64)
		low, _ := strconv.ParseFloat(ticker.Low24h, 64)
		vol, _ := strconv.ParseFloat(ticker.Vol24h, 64)
		volCcy, _ := strconv.ParseFloat(ticker.VolCcy24h, 64)
		bidPx, _ := strconv.ParseFloat(ticker.BidPx, 64)
		bidSz, _ := strconv.ParseFloat(ticker.BidSz, 64)
		askPx, _ := strconv.ParseFloat(ticker.AskPx, 64)
		askSz, _ := strconv.ParseFloat(ticker.AskSz, 64)
		ts, _ := strconv.ParseInt(ticker.Ts, 10, 64)

		change := last - open
		changeRate := 0.0
		if open > 0 {
			changeRate = change / open * 100
		}

		quotes = append(quotes, spot.Quote{
			Symbol:     fromOKXInstID(ticker.InstID),
			Name:       ticker.InstID,
			Latest:     last,
			Open:       open,
			High:       high,
			Low:        low,
			PreClose:   open, // OKX does not provide previous close directly
			Change:     change,
			ChangeRate: changeRate,
			Volume:     vol,
			Turnover:   volCcy,
			BidPrice:   bidPx,
			BidVolume:  bidSz,
			AskPrice:   askPx,
			AskVolume:  askSz,
			Timestamp:  time.UnixMilli(ts),
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
