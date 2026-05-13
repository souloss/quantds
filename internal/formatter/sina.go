package formatter

import (
	"strings"

	"github.com/souloss/quantds/domain"
)

// SinaFormatter converts symbols for Sina API format.
// Sina uses prefix+code format: "sh600519", "sz000001"
type SinaFormatter struct{}

var sinaPrefix = map[domain.Exchange]string{
	domain.ExchangeSH: "sh",
	domain.ExchangeSZ: "sz",
	domain.ExchangeBJ: "bj",
}

var sinaExchangeFromPrefix = map[string]domain.Exchange{
	"sh": domain.ExchangeSH,
	"sz": domain.ExchangeSZ,
	"bj": domain.ExchangeBJ,
}

func (f *SinaFormatter) Format(sym domain.Symbol) string {
	switch sym.Market {
	case domain.MarketCN:
		prefix, ok := sinaPrefix[sym.Exchange]
		if !ok {
			prefix = "sz"
		}
		return prefix + sym.Code
	default:
		return sym.Code
	}
}

func (f *SinaFormatter) Parse(raw string) (domain.Symbol, bool) {
	raw = strings.ToLower(raw)
	if len(raw) >= 8 {
		prefix := raw[:2]
		code := raw[2:]
		exchange, ok := sinaExchangeFromPrefix[prefix]
		if ok && len(code) == 6 && isPureDigits(code) {
			return domain.Symbol{
				Code:      strings.ToUpper(code),
				Market:    domain.MarketCN,
				Exchange:  exchange,
				AssetType: domain.AssetTypeStock,
				Standard:  strings.ToUpper(code) + ".CN." + string(exchange),
			}, true
		}
	}
	return domain.Symbol{}, false
}
