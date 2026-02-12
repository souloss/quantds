package tushare

import (
	"context"

	"github.com/souloss/quantds/clients/tushare"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// InstrumentAdapter adapts Tushare stock basic data to instrument list.
type InstrumentAdapter struct {
	client *tushare.Client
}

// NewInstrumentAdapter creates a new instrument adapter.
func NewInstrumentAdapter(client *tushare.Client) *InstrumentAdapter {
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

	// Map instrument.Exchange to Tushare exchange code
	exchange := ""
	switch req.Exchange {
	case instrument.ExchangeSH:
		exchange = "SSE"
	case instrument.ExchangeSZ:
		exchange = "SZSE"
	case instrument.ExchangeBJ:
		exchange = "BSE"
	}

	params := &tushare.StockBasicParams{
		Exchange: exchange,
		Status:   "L", // L=上市 D=退市 P=暂停上市
	}

	rows, record, err := a.client.GetStockBasic(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	items := make([]instrument.Instrument, 0, len(rows))
	for _, row := range rows {
		ex := instrument.ExchangeSZ
		switch row.Exchange {
		case "SSE":
			ex = instrument.ExchangeSH
		case "SZSE":
			ex = instrument.ExchangeSZ
		case "BSE":
			ex = instrument.ExchangeBJ
		}

		// Format date from YYYYMMDD to YYYY-MM-DD
		listDate := row.ListDate
		if len(listDate) == 8 {
			listDate = listDate[:4] + "-" + listDate[4:6] + "-" + listDate[6:]
		}

		items = append(items, instrument.Instrument{
			Symbol:   instrument.FormatSymbol(row.Symbol, ex),
			Code:     row.Symbol,
			Name:     row.Name,
			Exchange: ex,
			Industry: row.Industry,
			Market:   row.Market,
			ListDate: listDate,
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

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
