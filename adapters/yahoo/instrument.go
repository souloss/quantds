package yahoo

import (
	"context"

	"github.com/souloss/quantds/clients/yahoo"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// InstrumentAdapter adapts Yahoo Finance instrument data
type InstrumentAdapter struct {
	client *yahoo.Client
}

// NewInstrumentAdapter creates a new instrument adapter
func NewInstrumentAdapter(client *yahoo.Client) *InstrumentAdapter {
	return &InstrumentAdapter{client: client}
}

// Name returns the adapter name
func (a *InstrumentAdapter) Name() string {
	return Name
}

// SupportedMarkets returns supported markets
func (a *InstrumentAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

// CanHandle checks if the adapter can handle the symbol
func (a *InstrumentAdapter) CanHandle(symbol string) bool {
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

// Fetch retrieves instrument list
func (a *InstrumentAdapter) Fetch(ctx context.Context, _ request.Client, req instrument.Request) (instrument.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	params := &yahoo.InstrumentParams{
		Query: req.Market,
		Limit: req.PageSize,
	}

	if params.Limit <= 0 {
		params.Limit = 100
	}

	// Use search-based approach
	result, record, err := a.client.GetInstruments(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		// Fall back to predefined list
		result, record, err = a.client.GetAllUSStocks(ctx)
		trace.AddRequest(record)
		if err != nil {
			return instrument.Response{}, trace, err
		}
	}

	instruments := make([]instrument.Instrument, 0, len(result.Instruments))
	for _, data := range result.Instruments {
		exchange := domain.ExchangeNASDAQ
		if data.Exchange == "NYSE" {
			exchange = domain.ExchangeNYSE
		} else if data.Exchange == "AMEX" {
			exchange = domain.ExchangeAMEX
		}

		assetType := instrument.AssetTypeStock
		if data.AssetType == "ETF" {
			assetType = instrument.AssetTypeETF
		}

		instruments = append(instruments, instrument.Instrument{
			Symbol:    yahoo.FromYahooSymbol(data.Symbol, data.Exchange),
			Code:      data.Symbol,
			Name:      data.Name,
			Exchange:  exchange,
			Currency:  data.Currency,
			AssetType: assetType,
			Status:    instrument.StatusNormal,
		})
	}

	trace.Finish()
	return instrument.Response{
		Data:       instruments,
		Total:      len(instruments),
		Source:     Name,
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	}, trace, nil
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
