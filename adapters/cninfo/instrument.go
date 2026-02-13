package cninfo

import (
	"context"
	"strings"

	"github.com/souloss/quantds/clients/cninfo"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "cninfo"

// supportedMarkets defines the markets supported by CNInfo adapters
var supportedMarkets = []domain.Market{domain.MarketCN}

// InstrumentAdapter adapts CNInfo instrument data
type InstrumentAdapter struct {
	client *cninfo.Client
}

// NewInstrumentAdapter creates a new instrument adapter
func NewInstrumentAdapter(client *cninfo.Client) *InstrumentAdapter {
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

// determineExchange determines the exchange from orgId or code
func determineExchange(code, orgID string) string {
	// Try to determine from orgId
	if strings.HasSuffix(orgID, ".sh") || strings.HasSuffix(orgID, ".SH") {
		return "SH"
	}
	if strings.HasSuffix(orgID, ".sz") || strings.HasSuffix(orgID, ".SZ") {
		return "SZ"
	}
	if strings.HasSuffix(orgID, ".bj") || strings.HasSuffix(orgID, ".BJ") {
		return "BJ"
	}
	// Determine from code pattern
	if len(code) >= 2 {
		prefix := code[:2]
		switch {
		case prefix == "60" || prefix == "68" || prefix == "90":
			return "SH"
		case prefix == "00" || prefix == "30" || prefix == "20":
			return "SZ"
		case prefix == "43" || prefix == "83" || prefix == "87":
			return "BJ"
		}
	}
	return "SZ"
}

// Fetch retrieves instrument list
func (a *InstrumentAdapter) Fetch(ctx context.Context, _ request.Client, req instrument.Request) (instrument.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	// Get stock list from CNInfo
	rows, record, err := a.client.GetStockList(ctx)
	trace.AddRequest(record)

	if err != nil {
		return instrument.Response{}, trace, err
	}

	instruments := make([]instrument.Instrument, 0, len(rows))
	for _, data := range rows {
		exchange := instrument.ExchangeSZ
		ex := determineExchange(data.Code, data.OrgID)
		if ex == "SH" {
			exchange = instrument.ExchangeSH
		} else if ex == "BJ" {
			exchange = instrument.ExchangeBJ
		}

		// Filter by exchange if specified
		if req.Exchange != "" && exchange != req.Exchange {
			continue
		}

		instruments = append(instruments, instrument.Instrument{
			Symbol:   instrument.FormatSymbol(data.Code, exchange),
			Code:     data.Code,
			Name:     data.Name,
			Exchange: exchange,
		})
	}

	trace.Finish()
	return instrument.Response{
		Data:   instruments,
		Total:  len(instruments),
		Source: Name,
	}, trace, nil
}

var _ manager.Provider[instrument.Request, instrument.Response] = (*InstrumentAdapter)(nil)
