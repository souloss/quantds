package szse

import (
	"context"

	"github.com/souloss/quantds/clients/szse"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "szse"

var supportedMarkets = []domain.Market{domain.MarketCN}

// InstrumentAdapter adapts SZSE (Shenzhen Stock Exchange) instrument list data.
type InstrumentAdapter struct {
	client *szse.Client
}

// NewInstrumentAdapter creates a new instrument adapter.
func NewInstrumentAdapter(client *szse.Client) *InstrumentAdapter {
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

	result, record, err := a.client.GetStockList(ctx, &szse.StockListParams{
		CatalogID: "1110",
		TabKey:    "tab1",
		ShowType:  "xlsx",
	})
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	// SZSE returns Excel data as [][]string
	// Typical columns: [板块, 公司代码, 公司简称, 公司全称, A股代码, A股简称, A股上市日期, ...]
	// Skip header row (first row)
	items := make([]instrument.Instrument, 0, len(result.Data))
	for i, row := range result.Data {
		if i == 0 {
			continue // skip header row
		}
		if len(row) < 7 {
			continue // skip malformed rows
		}
		code := safeCol(row, 4)     // A股代码
		name := safeCol(row, 5)     // A股简称
		listDate := safeCol(row, 6) // A股上市日期
		if code == "" {
			code = safeCol(row, 1) // fallback to 公司代码
		}
		if name == "" {
			name = safeCol(row, 2) // fallback to 公司简称
		}
		if code == "" {
			continue
		}

		items = append(items, instrument.Instrument{
			Symbol:   instrument.FormatSymbol(code, instrument.ExchangeSZ),
			Code:     code,
			Name:     name,
			Exchange: instrument.ExchangeSZ,
			ListDate: listDate,
			Status:   instrument.GuessStatus(name),
		})
	}

	trace.Finish()
	return instrument.Response{
		Data:   items,
		Total:  len(items),
		Source: Name,
	}, trace, nil
}

func safeCol(row []string, idx int) string {
	if idx < len(row) {
		return row[idx]
	}
	return ""
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
