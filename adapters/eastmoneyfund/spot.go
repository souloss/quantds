package eastmoneyfund

import (
	"context"
	"fmt"
	"strconv"

	"github.com/souloss/quantds/clients/eastmoneyfund"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
)

// SpotAdapter adapts EastMoney fund estimate data to spot domain
type SpotAdapter struct {
	client *eastmoneyfund.Client
}

// NewSpotAdapter creates a new spot adapter
func NewSpotAdapter(client *eastmoneyfund.Client) *SpotAdapter {
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

func (a *SpotAdapter) Fetch(ctx context.Context, req spot.Request) (spot.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	if len(req.Symbols) == 0 {
		return spot.Response{}, trace, nil
	}

	quotes := make([]spot.Quote, 0, len(req.Symbols))
	for _, symbol := range req.Symbols {
		var sym domain.Symbol
		if err := sym.Parse(symbol); err != nil {
			continue
		}

		result, record, err := a.client.GetFundEstimate(ctx, &eastmoneyfund.FundEstimateParams{
			Code: sym.Code,
		})
		trace.AddRequest(record)
		if err != nil {
			continue
		}

		estNAV, _ := strconv.ParseFloat(result.EstNAV, 64)
		changeRate, _ := strconv.ParseFloat(result.EstChange, 64)

		quotes = append(quotes, spot.Quote{
			Symbol:     symbol,
			Name:       result.Name,
			Latest:     estNAV,
			ChangeRate: changeRate,
		})
	}

	if len(quotes) == 0 {
		return spot.Response{}, trace, fmt.Errorf("no fund estimate data returned")
	}

	trace.Finish()
	return spot.Response{
		Quotes:      quotes,
		Total:       len(quotes),
		Source:      Name,
		DataVersion: 1,
	}, trace, nil
}

var _ manager.Provider[spot.Request, spot.Response] = (*SpotAdapter)(nil)
