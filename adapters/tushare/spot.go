package tushare

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/tushare"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// SpotAdapter adapts Tushare realtime quote data to domain spot
type SpotAdapter struct {
	client *tushare.Client
}

// NewSpotAdapter creates a new spot adapter for Tushare
func NewSpotAdapter(client *tushare.Client) *SpotAdapter {
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

// Fetch retrieves realtime spot data from Tushare
func (a *SpotAdapter) Fetch(ctx context.Context, _ request.Client, req spot.Request) (spot.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	if len(req.Symbols) == 0 {
		return spot.Response{}, trace, nil
	}

	quotes := make([]spot.Quote, 0, len(req.Symbols))

	for _, symbol := range req.Symbols {
		tsCode, err := tushare.ToTushareSymbol(symbol)
		if err != nil {
			continue
		}

		// Use rt_k API for realtime daily data
		result, record, err := a.client.GetRtK(ctx, &tushare.RtKParams{
			TSCode: tsCode,
		})
		trace.AddRequest(record)

		if err != nil || len(result) == 0 {
			continue
		}

		for _, r := range result {
			changeRate := 0.0
			if r.PreClose > 0 {
				changeRate = r.Change / r.PreClose * 100
			}

			quotes = append(quotes, spot.Quote{
				Symbol:     r.TSCode,
				Name:       "", // rt_k doesn't return name
				Latest:     r.Close,
				Open:       r.Open,
				High:       r.High,
				Low:        r.Low,
				PreClose:   r.PreClose,
				Change:     r.Change,
				ChangeRate: changeRate,
				Volume:     r.Vol * 100,     // 手 -> 股
				Turnover:   r.Amount * 1000, // 千元 -> 元
				Timestamp:  time.Now(),
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
