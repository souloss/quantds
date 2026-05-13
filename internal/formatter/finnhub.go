package formatter

import (
	"github.com/souloss/quantds/domain"
)

// FinnhubFormatter converts symbols for Finnhub API format.
// Finnhub uses plain code: "AAPL"
type FinnhubFormatter struct{}

func (f *FinnhubFormatter) Format(sym domain.Symbol) string {
	return sym.Code
}

func (f *FinnhubFormatter) Parse(raw string) (domain.Symbol, bool) {
	var s domain.Symbol
	if err := s.Parse(raw); err == nil {
		return s, true
	}
	return domain.Symbol{}, false
}
