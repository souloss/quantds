package formatter

import (
	"github.com/souloss/quantds/domain"
)

// YahooFormatter converts symbols for Yahoo Finance API format.
// Yahoo uses plain code format: "AAPL" for US stocks
type YahooFormatter struct{}

func (f *YahooFormatter) Format(sym domain.Symbol) string {
	return sym.Code
}

func (f *YahooFormatter) Parse(raw string) (domain.Symbol, bool) {
	var s domain.Symbol
	if err := s.Parse(raw); err == nil {
		return s, true
	}
	return domain.Symbol{}, false
}
