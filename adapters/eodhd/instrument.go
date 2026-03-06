package eodhd

import (
	"context"

	"github.com/souloss/quantds/clients/eodhd"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

type InstrumentAdapter struct {
	client *eodhd.Client
}

func NewInstrumentAdapter(client *eodhd.Client) *InstrumentAdapter {
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

	result, record, err := a.client.GetExchangeSymbolList(ctx, &eodhd.ExchangeSymbolsParams{
		Exchange: "US",
	})
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(result.Symbols))
	for _, s := range result.Symbols {
		exchange := domain.ExchangeNASDAQ
		if s.Exchange == "NYSE" {
			exchange = domain.ExchangeNYSE
		} else if s.Exchange == "AMEX" {
			exchange = domain.ExchangeAMEX
		}

		assetType := instrument.AssetTypeStock
		if s.Type == "ETF" {
			assetType = instrument.AssetTypeETF
		} else if s.Type == "FUND" {
			assetType = instrument.AssetTypeFund
		} else if s.Type == "BOND" {
			assetType = instrument.AssetTypeBond
		}

		instruments = append(instruments, instrument.Instrument{
			Symbol:    domain.FormatSymbol(s.Code, exchange),
			Code:      s.Code,
			Name:      s.Name,
			Exchange:  exchange,
			Currency:  s.Currency,
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
