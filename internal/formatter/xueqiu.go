package formatter

import (
	"strings"

	"github.com/souloss/quantds/domain"
)

// XueqiuFormatter converts symbols for Xueqiu API format.
// Xueqiu uses prefix+code format: "SH600519", "SZ000001" (uppercase prefix)
type XueqiuFormatter struct{}

var xueqiuPrefix = map[domain.Exchange]string{
	domain.ExchangeSH: "SH",
	domain.ExchangeSZ: "SZ",
	domain.ExchangeBJ: "BJ",
}

var xueqiuExchangeFromPrefix = map[string]domain.Exchange{
	"SH": domain.ExchangeSH,
	"SZ": domain.ExchangeSZ,
	"BJ": domain.ExchangeBJ,
}

func (f *XueqiuFormatter) Format(sym domain.Symbol) string {
	switch sym.Market {
	case domain.MarketCN:
		prefix, ok := xueqiuPrefix[sym.Exchange]
		if !ok {
			prefix = "SZ"
		}
		return prefix + sym.Code
	case domain.MarketUS:
		return sym.Code
	default:
		return sym.Code
	}
}

func (f *XueqiuFormatter) Parse(raw string) (domain.Symbol, bool) {
	if len(raw) >= 8 {
		prefix := strings.ToUpper(raw[:2])
		code := raw[2:]
		exchange, ok := xueqiuExchangeFromPrefix[prefix]
		if ok && len(code) == 6 && isPureDigits(code) {
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
