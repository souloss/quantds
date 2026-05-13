package formatter

import (
	"github.com/souloss/quantds/domain"
)

// CNInfoFormatter converts symbols for CNInfo API format.
// CNInfo uses plain code format: "000001" (no exchange suffix)
type CNInfoFormatter struct{}

func (f *CNInfoFormatter) Format(sym domain.Symbol) string {
	return sym.Code
}

func (f *CNInfoFormatter) Parse(raw string) (domain.Symbol, bool) {
	var s domain.Symbol
	if err := s.Parse(raw); err == nil {
		return s, true
	}
	return domain.Symbol{}, false
}
