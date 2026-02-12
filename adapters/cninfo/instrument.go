package cninfo

import (
	"context"

	"github.com/souloss/quantds/clients/cninfo"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "cninfo"

var supportedMarkets = []domain.Market{domain.MarketCN}

// InstrumentAdapter adapts CNINFO instrument list data.
type InstrumentAdapter struct {
	client *cninfo.Client
}

// NewInstrumentAdapter creates a new instrument adapter.
func NewInstrumentAdapter(client *cninfo.Client) *InstrumentAdapter {
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

	rows, record, err := a.client.GetStockList(ctx)
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	items := make([]instrument.Instrument, 0, len(rows))
	for _, row := range rows {
		code := row.Code
		exchange := inferExchangeFromCode(code)

		// Filter by requested exchange if specified
		if req.Exchange != "" && exchange != req.Exchange {
			continue
		}

		items = append(items, instrument.Instrument{
			Symbol:   instrument.FormatSymbol(code, exchange),
			Code:     code,
			Name:     row.Name,
			Exchange: exchange,
			Status:   instrument.GuessStatus(row.Name),
		})
	}

	trace.Finish()
	return instrument.Response{
		Data:   items,
		Total:  len(items),
		Source: Name,
	}, trace, nil
}

// inferExchangeFromCode infers the exchange from A-share stock code prefix.
func inferExchangeFromCode(code string) instrument.Exchange {
	if len(code) < 2 {
		return instrument.ExchangeSZ
	}
	prefix := code[:2]
	switch {
	case prefix == "60" || prefix == "68" || prefix == "90":
		return instrument.ExchangeSH
	case prefix == "00" || prefix == "30" || prefix == "20":
		return instrument.ExchangeSZ
	case prefix == "43" || prefix == "83" || prefix == "87":
		return instrument.ExchangeBJ
	default:
		return instrument.ExchangeSZ
	}
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
