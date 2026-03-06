package twelvedata

import (
	"context"

	"github.com/souloss/quantds/clients/twelvedata"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

type InstrumentAdapter struct {
	client *twelvedata.Client
}

func NewInstrumentAdapter(client *twelvedata.Client) *InstrumentAdapter {
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

	result, record, err := a.client.GetStocksList(ctx, &twelvedata.ListParams{
		Country: "United States",
	})
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(result.Data))
	for _, item := range result.Data {
		exchange := domain.ExchangeNASDAQ
		if item.Exchange == "NYSE" {
			exchange = domain.ExchangeNYSE
		} else if item.Exchange == "AMEX" {
			exchange = domain.ExchangeAMEX
		}

		assetType := instrument.AssetTypeStock
		if item.Type == "ETF" {
			assetType = instrument.AssetTypeETF
		} else if item.Type == "Fund" || item.Type == "Mutual Fund" {
			assetType = instrument.AssetTypeFund
		} else if item.Type == "Bond" {
			assetType = instrument.AssetTypeBond
		}

		instruments = append(instruments, instrument.Instrument{
			Symbol:    domain.FormatSymbol(item.Symbol, exchange),
			Code:      item.Symbol,
			Name:      item.Name,
			Exchange:  exchange,
			Currency:  item.Currency,
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
