package eastmoney

import (
	"context"

	"github.com/souloss/quantds/clients/eastmoney"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// InstrumentAdapter adapts Eastmoney instrument data
type InstrumentAdapter struct {
	client *eastmoney.Client
}

// NewInstrumentAdapter creates a new instrument adapter
func NewInstrumentAdapter(client *eastmoney.Client) *InstrumentAdapter {
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

	params := &eastmoney.InstrumentParams{
		PageSize: 5000,
	}

	if req.Exchange != "" {
		params.Market = string(req.Exchange)
	}

	result, record, err := a.client.GetInstruments(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(result.Data))
	for _, data := range result.Data {
		exchange := instrument.ExchangeSZ
		if data.MarketID == 1 {
			exchange = instrument.ExchangeSH
		}
		instruments = append(instruments, instrument.Instrument{
			Symbol:   instrument.FormatSymbol(data.Code, exchange),
			Code:     data.Code,
			Name:     data.Name,
			Exchange: exchange,
		})
	}

	trace.Finish()
	return instrument.Response{
		Data:   instruments,
		Source: Name,
	}, trace, nil
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
