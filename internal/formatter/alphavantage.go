package formatter

import (
	"strings"

	"github.com/souloss/quantds/domain"
)

// AlphaVantageFormatter converts symbols for Alpha Vantage API format.
// US stocks: plain code "AAPL"
// Forex: "EUR_USD" (underscore separator)
type AlphaVantageFormatter struct{}

func (f *AlphaVantageFormatter) Format(sym domain.Symbol) string {
	switch sym.Market {
	case domain.MarketForex:
		// Forex pairs use underscore: EURUSD → EUR_USD
		if len(sym.Code) == 6 && isAllLetters(sym.Code) {
			return sym.Code[:3] + "_" + sym.Code[3:]
		}
		return sym.Code
	default:
		return sym.Code
	}
}

func (f *AlphaVantageFormatter) Parse(raw string) (domain.Symbol, bool) {
	// Handle forex format: EUR_USD
	if strings.Contains(raw, "_") {
		parts := strings.SplitN(raw, "_", 2)
		if len(parts) == 2 && len(parts[0]) == 3 && len(parts[1]) == 3 {
			code := parts[0] + parts[1]
			return domain.Symbol{
				Code:      code,
				Market:    domain.MarketForex,
				Exchange:  domain.ExchangeForexSpot,
				AssetType: domain.AssetTypeForex,
				Standard:  code + ".FOREX.FOREX_SPOT",
			}, true
		}
	}
	// Plain code → try general parse
	var s domain.Symbol
	if err := s.Parse(raw); err == nil {
		return s, true
	}
	return domain.Symbol{}, false
}
