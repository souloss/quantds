package coingecko

import (
	"context"
	"strings"

	"github.com/souloss/quantds/clients/coingecko"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
)

// InstrumentAdapter adapts CoinGecko coins/list data to instrument domain
type InstrumentAdapter struct {
	client *coingecko.Client
}

// NewInstrumentAdapter creates a new instrument adapter
func NewInstrumentAdapter(client *coingecko.Client) *InstrumentAdapter {
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

func (a *InstrumentAdapter) Fetch(ctx context.Context, req instrument.Request) (instrument.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	result, record, err := a.client.GetCoinsList(ctx, &coingecko.CoinsListRequest{
		IncludePlatform: false,
	})
	trace.AddRequest(record)
	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(result))
	for _, coin := range result {
		instruments = append(instruments, instrument.Instrument{
			Symbol: strings.ToUpper(coin.Symbol) + ".CRYPTO",
			Code:   coin.Symbol,
			Name:   coin.Name,
			Market: "CRYPTO",
		})
	}

	trace.Finish()
	return instrument.Response{
		Data:        instruments,
		Total:       len(instruments),
		Source:      Name,
		DataVersion: 1,
	}, trace, nil
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
