package finnhub

import (
	"context"

	"github.com/souloss/quantds/clients/finnhub"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/profile"
	"github.com/souloss/quantds/manager"
)

// usMarkets defines markets supported by Finnhub profile/financial/news adapters
var usMarkets = []domain.Market{domain.MarketUS}

// ProfileAdapter adapts Finnhub company profile data for US stocks
type ProfileAdapter struct {
	client *finnhub.Client
}

// NewProfileAdapter creates a new profile adapter
func NewProfileAdapter(client *finnhub.Client) *ProfileAdapter {
	return &ProfileAdapter{client: client}
}

func (a *ProfileAdapter) Name() string                      { return Name }
func (a *ProfileAdapter) SupportedMarkets() []domain.Market { return usMarkets }

func (a *ProfileAdapter) CanHandle(symbol string) bool {
	var sym domain.Symbol
	if err := sym.Parse(symbol); err != nil {
		return false
	}
	for _, m := range usMarkets {
		if sym.Market == m {
			return true
		}
	}
	return false
}

func (a *ProfileAdapter) Fetch(ctx context.Context, req profile.Request) (profile.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	var sym domain.Symbol
	if err := sym.Parse(req.Symbol); err != nil {
		return profile.Response{}, trace, err
	}

	result, record, err := a.client.GetProfile2(ctx, &finnhub.ProfileParams{Symbol: sym.Code})
	trace.AddRequest(record)

	if err != nil {
		return profile.Response{}, trace, err
	}

	p := profile.Profile{
		Symbol:      req.Symbol,
		Name:        result.Name,
		Industry:    result.Industry,
		ListingDate: result.IPO,
		MarketCap:   result.MarketCap,
		Website:     result.WebURL,
		Phone:       result.Phone,
	}

	// Determine market cap category
	if p.MarketCap > 10e9 { // > $10B
		p.MarketCapCategory = "Large"
	} else if p.MarketCap > 2e9 { // > $2B
		p.MarketCapCategory = "Mid"
	} else {
		p.MarketCapCategory = "Small"
	}

	trace.Finish()
	return profile.Response{
		Data:        p,
		Source:      Name,
		DataVersion: 1,
	}, trace, nil
}

var _ manager.Provider[profile.Request, profile.Response] = (*ProfileAdapter)(nil)
