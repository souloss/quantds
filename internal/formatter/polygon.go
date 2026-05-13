package formatter

import (
	"github.com/souloss/quantds/domain"
)

// PolygonFormatter converts symbols for Polygon API format.
// Polygon uses plain code: "AAPL"
type PolygonFormatter struct{}

func (f *PolygonFormatter) Format(sym domain.Symbol) string {
	return sym.Code
}

func (f *PolygonFormatter) Parse(raw string) (domain.Symbol, bool) {
	var s domain.Symbol
	if err := s.Parse(raw); err == nil {
		return s, true
	}
	return domain.Symbol{}, false
}
