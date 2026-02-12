package bse

import (
	"context"

	"github.com/souloss/quantds/clients/bse"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "bse"

var supportedMarkets = []domain.Market{domain.MarketCN}

// InstrumentAdapter adapts BSE (Beijing Stock Exchange) instrument list data.
type InstrumentAdapter struct {
	client *bse.Client
}

// NewInstrumentAdapter creates a new instrument adapter.
func NewInstrumentAdapter(client *bse.Client) *InstrumentAdapter {
	return &InstrumentAdapter{client: client}
}

func (a *InstrumentAdapter) Name() string                      { return Name }
func (a *InstrumentAdapter) SupportedMarkets() []domain.Market { return supportedMarkets }

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

func (a *InstrumentAdapter) Fetch(ctx context.Context, _ request.Client, req instrument.Request) (instrument.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	rows, records, err := a.client.GetStockList(ctx)
	for _, record := range records {
		trace.AddRequest(record)
	}

	if err != nil {
		return instrument.Response{}, trace, err
	}

	items := make([]instrument.Instrument, 0, len(rows))
	for _, row := range rows {
		items = append(items, instrument.Instrument{
			Symbol:   instrument.FormatSymbol(row.StockCode, instrument.ExchangeBJ),
			Code:     row.StockCode,
			Name:     row.StockName,
			Exchange: instrument.ExchangeBJ,
			Industry: row.Industry,
			ListDate: row.ListDate,
			Status:   instrument.GuessStatus(row.StockName),
		})
	}

	trace.Finish()
	return instrument.Response{
		Data:   items,
		Total:  len(items),
		Source: Name,
	}, trace, nil
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
