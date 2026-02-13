package eastmoneyhk

import (
	"context"

	"github.com/souloss/quantds/clients/eastmoneyhk"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// InstrumentAdapter adapts EastMoney HK instrument data
type InstrumentAdapter struct {
	client *eastmoneyhk.Client
}

// NewInstrumentAdapter creates a new instrument adapter
func NewInstrumentAdapter(client *eastmoneyhk.Client) *InstrumentAdapter {
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

	params := &eastmoneyhk.InstrumentParams{
		PageSize:   req.PageSize,
		PageNumber: req.PageNumber,
	}

	if params.PageSize <= 0 {
		params.PageSize = 500
	}

	result, record, err := a.client.GetInstruments(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(result.Instruments))
	for _, data := range result.Instruments {
		status := instrument.StatusNormal
		// Determine status based on price data
		if data.LatestPrice <= 0 {
			status = instrument.StatusSuspended
		}

		instruments = append(instruments, instrument.Instrument{
			Symbol:    data.Symbol,
			Code:      data.Code,
			Name:      data.Name,
			Exchange:  domain.ExchangeHKEX,
			Market:    "HK",
			ListDate:  data.ListDate,
			Status:    status,
			AssetType: instrument.AssetTypeStock,
			Currency:  "HKD",
		})
	}

	trace.Finish()
	return instrument.Response{
		Data:       instruments,
		Total:      result.Total,
		Source:     Name,
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	}, trace, nil
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
