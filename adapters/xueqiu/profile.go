package xueqiu

import (
	"context"

	"github.com/souloss/quantds/clients/xueqiu"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/profile"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// ProfileAdapter adapts Xueqiu quote detail data to domain profile
type ProfileAdapter struct {
	client *xueqiu.Client
}

// NewProfileAdapter creates a new profile adapter for Xueqiu
func NewProfileAdapter(client *xueqiu.Client) *ProfileAdapter {
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

// Fetch retrieves profile data from Xueqiu
func (a *ProfileAdapter) Fetch(ctx context.Context, _ request.Client, req profile.Request) (profile.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	params := &xueqiu.QuoteDetailParams{
		Symbol: req.Symbol,
		Extend: true,
	}

	result, record, err := a.client.GetQuoteDetail(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return profile.Response{}, trace, err
	}

	trace.Finish()
	return profile.Response{
		Data: profile.Profile{
			Symbol:       result.Symbol,
			Name:         result.Name,
			ListingDate:  result.ListDate,
			Industry:     result.Industry,
			Sector:       result.Sector,
			Province:     result.Province,
			City:         result.City,
			Website:      result.Website,
			Email:        result.Email,
			Address:      result.Office,
			Chairman:     result.Chairman,
			CEO:          result.Manager,
			Secretary:    result.Secretary,
			RegCapital:   result.RegCapital,
			SetupDate:    result.SetupDate,
			Employees:    result.Employees,
			MainBusiness: result.MainBusiness,
			Introduction: result.Description,
			// 行情数据
			Open:          result.Open,
			High:          result.High,
			Low:           result.Low,
			Close:         result.Current,
			PreClose:      result.PreClose,
			Volume:        result.Volume,
			Amount:        result.Amount,
			TurnoverRate:  result.TurnoverRate,
			ChangePercent: result.Percent,
			Amplitude:     result.Amplitude,
			// 估值指标
			PE:       result.PE,
			PB:       result.PB,
			PS:       result.PS,
			DivYield: result.DividendYield,
			// 股本数据
			TotalShares: result.TotalShares,
			FloatShares: result.FloatShares,
			MarketCap:   result.TotalMarketCap,
			FloatCap:    result.FloatMarketCap,
		},
		Source: Name,
	}, trace, nil
}

var _ manager.Provider[profile.Request, profile.Response] = (*ProfileAdapter)(nil)
