package coingecko

import (
	"context"
	"fmt"
	"strings"

	"github.com/souloss/quantds/clients/coingecko"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/manager"
)

// SpotAdapter adapts CoinGecko simple/price data to spot domain
type SpotAdapter struct {
	client *coingecko.Client
}

// NewSpotAdapter creates a new spot adapter
func NewSpotAdapter(client *coingecko.Client) *SpotAdapter {
	return &SpotAdapter{client: client}
}

func (a *SpotAdapter) Name() string                      { return Name }
func (a *SpotAdapter) SupportedMarkets() []domain.Market { return supportedMarkets }

func (a *SpotAdapter) CanHandle(symbol string) bool {
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

func (a *SpotAdapter) Fetch(ctx context.Context, req spot.Request) (spot.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	if len(req.Symbols) == 0 {
		return spot.Response{}, trace, nil
	}

	// Extract coin IDs from symbols (e.g., "BTC.USDT" -> "bitcoin")
	coinIDs := make([]string, 0, len(req.Symbols))
	for _, symbol := range req.Symbols {
		var sym domain.Symbol
		if err := sym.Parse(symbol); err != nil {
			continue
		}
		coinID := symbolToCoinID(sym.Code)
		if coinID != "" {
			coinIDs = append(coinIDs, coinID)
		}
	}

	if len(coinIDs) == 0 {
		return spot.Response{}, trace, fmt.Errorf("no valid coin IDs found")
	}

	// Determine quote currency (default: usd)
	quoteCurrency := "usd"
	for _, symbol := range req.Symbols {
		var sym domain.Symbol
		if err := sym.Parse(symbol); err != nil {
			continue
		}
		if sym.Exchange != "" {
			quoteCurrency = strings.ToLower(string(sym.Exchange))
			break
		}
	}

	result, record, err := a.client.GetSimplePrice(ctx, &coingecko.SimplePriceRequest{
		IDs:               coinIDs,
		VsCurrencies:      []string{quoteCurrency},
		Include24hrChange: true,
	})
	trace.AddRequest(record)
	if err != nil {
		return spot.Response{}, trace, err
	}

	quotes := make([]spot.Quote, 0, len(req.Symbols))
	for _, symbol := range req.Symbols {
		var sym domain.Symbol
		if err := sym.Parse(symbol); err != nil {
			continue
		}
		coinID := symbolToCoinID(sym.Code)
		if coinID == "" {
			continue
		}
		priceData, ok := result[coinID]
		if !ok {
			continue
		}

		latest := priceData[quoteCurrency]
		changeRate := priceData[quoteCurrency+"_24h_change"]

		quotes = append(quotes, spot.Quote{
			Symbol:     symbol,
			Name:       coinID,
			Latest:     latest,
			ChangeRate: changeRate,
		})
	}

	trace.Finish()
	return spot.Response{
		Quotes:      quotes,
		Total:       len(quotes),
		Source:      Name,
		DataVersion: 1,
	}, trace, nil
}

// symbolToCoinID maps common crypto symbols to CoinGecko coin IDs
func symbolToCoinID(code string) string {
	m := map[string]string{
		"BTC":   "bitcoin",
		"ETH":   "ethereum",
		"BNB":   "binancecoin",
		"SOL":   "solana",
		"XRP":   "ripple",
		"ADA":   "cardano",
		"DOGE":  "dogecoin",
		"DOT":   "polkadot",
		"MATIC": "matic-network",
		"AVAX":  "avalanche-2",
		"LINK":  "chainlink",
		"UNI":   "uniswap",
		"ATOM":  "cosmos",
		"LTC":   "litecoin",
		"ETC":   "ethereum-classic",
		"XLM":   "stellar",
		"ALGO":  "algorand",
		"VET":   "vechain",
		"FIL":   "filecoin",
		"TRX":   "tron",
	}
	if id, ok := m[strings.ToUpper(code)]; ok {
		return id
	}
	return strings.ToLower(code)
}

var _ manager.Provider[spot.Request, spot.Response] = (*SpotAdapter)(nil)
