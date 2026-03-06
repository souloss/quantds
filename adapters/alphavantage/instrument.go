package alphavantage

import (
	"context"

	"github.com/souloss/quantds/clients/alphavantage"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

type InstrumentAdapter struct {
	client *alphavantage.Client
}

func NewInstrumentAdapter(client *alphavantage.Client) *InstrumentAdapter {
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

	keywords := req.Market
	if keywords == "" {
		keywords = "US"
	}

	result, record, err := a.client.SearchSymbol(ctx, &alphavantage.SearchParams{Keywords: keywords})
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(result.Matches))
	for _, m := range result.Matches {
		exchange := domain.ExchangeNASDAQ
		assetType := instrument.AssetTypeStock
		if m.Type == "ETF" {
			assetType = instrument.AssetTypeETF
		}

		instruments = append(instruments, instrument.Instrument{
			Symbol:    domain.FormatSymbol(m.Symbol, exchange),
			Code:      m.Symbol,
			Name:      m.Name,
			Exchange:  exchange,
			Currency:  m.Currency,
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
