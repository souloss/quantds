package formatter

import (
	"fmt"
	"strings"

	"github.com/souloss/quantds/domain"
)

// EastMoneyHKFormatter converts symbols for East Money HK API format.
// East Money HK uses secid format: "116.00700" (HK market code is 116)
type EastMoneyHKFormatter struct{}

func (f *EastMoneyHKFormatter) Format(sym domain.Symbol) string {
	switch sym.Market {
	case domain.MarketHK:
		return "116." + sym.Code
	default:
		return sym.Code
	}
}

func (f *EastMoneyHKFormatter) Parse(raw string) (domain.Symbol, bool) {
	parts := strings.SplitN(raw, ".", 2)
	if len(parts) == 2 && parts[0] == "116" {
		code := parts[1]
		if len(code) == 5 {
			return domain.Symbol{
				Code:      code,
				Market:    domain.MarketHK,
				Exchange:  domain.ExchangeHKEX,
				AssetType: domain.AssetTypeStock,
				Standard:  fmt.Sprintf("%s.HK.HKEX", code),
			}, true
		}
	}
	return domain.Symbol{}, false
}
