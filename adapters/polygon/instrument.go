package polygon

import (
	"context"

	"github.com/souloss/quantds/clients/polygon"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

type InstrumentAdapter struct {
	client *polygon.Client
}

func NewInstrumentAdapter(client *polygon.Client) *InstrumentAdapter {
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

	limit := req.PageSize
	if limit <= 0 {
		limit = 100
	}

	result, record, err := a.client.GetTickers(ctx, &polygon.TickerParams{
		Market: "stocks",
		Limit:  limit,
	})
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(result.Tickers))
	for _, t := range result.Tickers {
		exchange := domain.ExchangeNASDAQ
		if t.PrimaryExchange == "XNYS" {
			exchange = domain.ExchangeNYSE
		} else if t.PrimaryExchange == "XASE" {
			exchange = domain.ExchangeAMEX
		}

		assetType := instrument.AssetTypeStock
		if t.Type == "ETF" || t.Type == "ETP" {
			assetType = instrument.AssetTypeETF
		}

		instruments = append(instruments, instrument.Instrument{
			Symbol:    domain.FormatSymbol(t.Ticker, exchange),
			Code:      t.Ticker,
			Name:      t.Name,
			Exchange:  exchange,
			Currency:  t.CurrencyName,
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
