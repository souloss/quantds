package formatter

import (
	"strings"

	"github.com/souloss/quantds/domain"
)

// TushareFormatter converts symbols for Tushare API format.
// Tushare uses CODE.EXCHANGE format for A-shares: "000001.SZ", "600519.SH"
type TushareFormatter struct{}

func (f *TushareFormatter) Format(sym domain.Symbol) string {
	switch sym.Market {
	case domain.MarketCN:
		return sym.Code + "." + string(sym.Exchange)
	default:
		return sym.Code
	}
}

func (f *TushareFormatter) Parse(raw string) (domain.Symbol, bool) {
	parts := strings.SplitN(raw, ".", 2)
	if len(parts) == 2 {
		var s domain.Symbol
		if err := s.Parse(parts[0] + "." + parts[1]); err == nil {
			return s, true
		}
		// Try as CN 2-level format: "000001.SZ"
		code := parts[0]
		exchange := domain.Exchange(parts[1])
		if len(code) == 6 && isPureDigits(code) {
			return domain.Symbol{
				Code:      code,
				Market:    domain.MarketCN,
				Exchange:  exchange,
				AssetType: domain.AssetTypeStock,
				Standard:  code + ".CN." + string(exchange),
			}, true
		}
	}
	return domain.Symbol{}, false
}
