package xueqiu

import (
	"context"

	"github.com/souloss/quantds/clients/xueqiu"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// InstrumentAdapter adapts Xueqiu stock list data to domain instrument
type InstrumentAdapter struct {
	client *xueqiu.Client
}

// NewInstrumentAdapter creates a new instrument adapter for Xueqiu
func NewInstrumentAdapter(client *xueqiu.Client) *InstrumentAdapter {
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

// CanHandle checks if the adapter can handle the request
func (a *InstrumentAdapter) CanHandle(symbol string) bool {
	return true // Can handle all symbols
}

// Fetch retrieves instrument list from Xueqiu
func (a *InstrumentAdapter) Fetch(ctx context.Context, _ request.Client, req instrument.Request) (instrument.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	// Map domain exchange to Xueqiu exchange
	exchange := ""
	switch req.Exchange {
	case domain.ExchangeSH:
		exchange = "SH"
	case domain.ExchangeSZ:
		exchange = "SZ"
	case domain.ExchangeBJ:
		exchange = "BJ"
	}

	// Map domain asset type to Xueqiu board type
	boardType := "all"
	switch req.AssetType {
	case instrument.AssetTypeStock:
		boardType = "all"
	}

	params := &xueqiu.StockListParams{
		Market:    "CN",
		Exchange:  exchange,
		BoardType: boardType,
		Page:      req.PageNumber,
		Size:      req.PageSize,
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 90
	}

	result, record, err := a.client.GetStockList(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(result.Items))
	for _, item := range result.Items {
		// Determine exchange from symbol
		exchange := domain.ExchangeSZ
		if len(item.Symbol) >= 2 {
			prefix := item.Symbol[:2]
			switch prefix {
			case "SH":
				exchange = domain.ExchangeSH
			case "SZ":
				exchange = domain.ExchangeSZ
			case "BJ":
				exchange = domain.ExchangeBJ
			}
		}

		// Determine status
		status := instrument.StatusNormal
		if item.Status == 2 {
			status = instrument.StatusSuspended
		} else if item.Status == 3 {
			status = instrument.StatusDelisted
		}

		instruments = append(instruments, instrument.Instrument{
			Symbol:    item.Symbol,
			Code:      item.Code,
			Name:      item.Name,
			Exchange:  exchange,
			AssetType: instrument.AssetTypeStock,
			Status:    status,
		})
	}

	trace.Finish()
	return instrument.Response{
		Data:       instruments,
		Total:      result.Total,
		Source:     Name,
		PageNumber: result.Page,
		PageSize:   len(instruments),
	}, trace, nil
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
