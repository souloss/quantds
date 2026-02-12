package eastmoney

import (
	"context"

	"github.com/souloss/quantds/clients/eastmoney"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/profile"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// ProfileAdapter adapts Eastmoney profile data
type ProfileAdapter struct {
	client *eastmoney.Client
}

// NewProfileAdapter creates a new profile adapter
func NewProfileAdapter(client *eastmoney.Client) *ProfileAdapter {
	return &ProfileAdapter{client: client}
}

// Name returns the adapter name
func (a *ProfileAdapter) Name() string {
	return Name
}

// SupportedMarkets returns supported markets
func (a *ProfileAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

// CanHandle checks if the adapter can handle the symbol
func (a *ProfileAdapter) CanHandle(symbol string) bool {
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

// Fetch retrieves profile data
func (a *ProfileAdapter) Fetch(ctx context.Context, _ request.Client, req profile.Request) (profile.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	params := &eastmoney.ProfileParams{
		Symbol: req.Symbol,
	}

	result, record, err := a.client.GetProfile(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return profile.Response{}, trace, err
	}

	trace.Finish()
	return profile.Response{
		Data: profile.Profile{
			Symbol:        result.Code,
			Name:          result.Name,
			ListingDate:   result.ListingDate,
			Currency:      result.Currency,
			Open:          result.Open,
			High:          result.High,
			Low:           result.Low,
			Close:         result.LatestPrice,
			PreClose:      result.PreClose,
			Volume:        result.Volume,
			Amount:        result.Amount,
			TurnoverRate:  result.TurnoverRate,
			ChangePercent: result.ChangePct,
			PE:            result.PEDynamic,
			PEStatic:      result.PEStatic,
			PB:            result.PBRatio,
			PS:            result.PSTTM,
			TotalShares:   result.TotalShares,
			FloatShares:   result.FloatShares,
			MarketCap:     result.TotalMarketCap,
			FloatCap:      result.FloatMarketCap,
			EPS:           result.EPS,
			Industry:      result.Industry,
			VolumeRatio:   result.VolumeRatio,
			Amplitude:     result.Amplitude,
		},
		Source: Name,
	}, trace, nil
}

var _ manager.Provider[profile.Request, profile.Response] = (*ProfileAdapter)(nil)
