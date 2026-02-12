package sina

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/sina"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

type SpotAdapter struct {
	client *sina.Client
}

func NewSpotAdapter(client *sina.Client) *SpotAdapter {
	return &SpotAdapter{client: client}
}

func (a *SpotAdapter) Name() string {
	return Name
}

func (a *SpotAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

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

	if len(req.Symbols) == 0 {
		return spot.Response{}, trace, nil
	}

	params := &sina.SpotParams{
		Symbols: req.Symbols,
	}

	result, record, err := a.client.GetSpot(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return spot.Response{}, trace, err
	}

	quotes := make([]spot.Quote, 0, len(result.Data))
	for _, q := range result.Data {
		change := q.Latest - q.PreClose
		changeRate := 0.0
		if q.PreClose > 0 {
			changeRate = change / q.PreClose * 100
		}
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
			Change:     change,
			ChangeRate: changeRate,
			Volume:     q.Volume * 100,
			Turnover:   q.Amount * 10000,
			Amplitude:  amplitude,
			Timestamp:  parseSpotTime(q.Date, q.Time),
		})
	}

	trace.Finish()
	return spot.Response{
		Quotes: quotes,
		Total:  len(quotes),
		Source: Name,
	}, trace, nil
}

func parseSpotTime(date, tm string) time.Time {
	if date == "" || tm == "" {
		return time.Now()
	}
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", date+" "+tm, timeLoc)
	return t
}

var timeLoc, _ = time.LoadLocation("Asia/Shanghai")

var _ manager.Provider[spot.Request, spot.Response] = (*SpotAdapter)(nil)
