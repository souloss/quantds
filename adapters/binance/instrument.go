package binance

import (
	"context"

	"github.com/souloss/quantds/clients/binance"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// InstrumentAdapter adapts Binance instrument data
type InstrumentAdapter struct {
	client *binance.Client
}

// NewInstrumentAdapter creates a new instrument adapter
func NewInstrumentAdapter(client *binance.Client) *InstrumentAdapter {
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
	return binance.IsCryptoSymbol(symbol)
}

// Fetch retrieves instrument list
func (a *InstrumentAdapter) Fetch(ctx context.Context, _ request.Client, req instrument.Request) (instrument.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	params := &binance.InstrumentParams{
		PageSize: req.PageSize,
	}

	// Get all instruments with filters
	if req.Market != "" {
		params.Quote = req.Market // Use Market as quote asset filter (e.g., "USDT")
	}

	result, record, err := a.client.GetInstruments(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(result.Instruments))
	for _, data := range result.Instruments {
		instruments = append(instruments, instrument.Instrument{
			Symbol:    binance.FromBinanceSymbol(data.Symbol),
			Code:      data.Symbol,
			Name:      data.BaseAsset + "/" + data.QuoteAsset,
			Exchange:  domain.ExchangeBinance,
			Market:    data.Market,
			AssetType: instrument.AssetTypeStock,
			Currency:  data.QuoteAsset,
			Status:    instrument.StatusNormal,
		})
	}

	trace.Finish()
	return instrument.Response{
		Data:       instruments,
		Total:      result.Total,
		Source:     Name,
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	}, trace, nil
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
