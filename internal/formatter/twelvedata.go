package formatter

import (
	"github.com/souloss/quantds/domain"
)

// TwelveDataFormatter converts symbols for TwelveData API format.
// US stocks: plain code "AAPL"
// Forex: plain code "EURUSD"
type TwelveDataFormatter struct{}

func (f *TwelveDataFormatter) Format(sym domain.Symbol) string {
	return sym.Code
}

func (f *TwelveDataFormatter) Parse(raw string) (domain.Symbol, bool) {
	var s domain.Symbol
	if err := s.Parse(raw); err == nil {
		return s, true
	}
	return domain.Symbol{}, false
}
