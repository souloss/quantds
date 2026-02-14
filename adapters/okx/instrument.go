package okx

import (
	"context"

	"github.com/souloss/quantds/clients/okx"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// InstrumentAdapter adapts OKX instrument data
type InstrumentAdapter struct {
	client *okx.Client
}

// NewInstrumentAdapter creates a new instrument adapter
func NewInstrumentAdapter(client *okx.Client) *InstrumentAdapter {
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
	return true // OKX instruments are generic
}

// Fetch retrieves instrument list from OKX
func (a *InstrumentAdapter) Fetch(ctx context.Context, _ request.Client, req instrument.Request) (instrument.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	// Default to SPOT instruments
	instType := okx.InstTypeSpot
	if req.Market != "" {
		switch req.Market {
		case "SWAP":
			instType = okx.InstTypeSwap
		case "FUTURES":
			instType = okx.InstTypeFutures
		case "OPTION":
			instType = okx.InstTypeOption
		}
	}

	result, record, err := a.client.GetInstruments(ctx, &okx.InstrumentParams{
		InstType: instType,
	})
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(result.Instruments))
	for _, data := range result.Instruments {
		name := data.BaseCcy + "/" + data.QuoteCcy
		if data.BaseCcy == "" {
			name = data.InstID
		}

		instruments = append(instruments, instrument.Instrument{
			Symbol:    fromOKXInstID(data.InstID),
			Code:      data.InstID,
			Name:      name,
			Exchange:  domain.ExchangeOKX,
			Market:    data.InstType,
			AssetType: instrument.AssetTypeStock,
			Currency:  data.QuoteCcy,
			Status:    mapOKXState(data.State),
		})
	}

	// Apply pagination if requested
	total := len(instruments)
	if req.PageSize > 0 {
		start := 0
		if req.PageNumber > 1 {
			start = (req.PageNumber - 1) * req.PageSize
		}
		end := start + req.PageSize
		if start > total {
			instruments = nil
		} else {
			if end > total {
				end = total
			}
			instruments = instruments[start:end]
		}
	}

	trace.Finish()
	return instrument.Response{
		Data:       instruments,
		Total:      total,
		Source:     Name,
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	}, trace, nil
}

// mapOKXState maps OKX instrument state to domain status
func mapOKXState(state string) instrument.Status {
	switch state {
	case "live":
		return instrument.StatusNormal
	case "suspend":
		return instrument.StatusSuspended
	case "expired":
		return instrument.StatusDelisted
	default:
		return instrument.StatusNormal
	}
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
