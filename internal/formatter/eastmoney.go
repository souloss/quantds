package formatter

import (
	"fmt"
	"strings"

	"github.com/souloss/quantds/domain"
)

// EastMoneyFormatter converts symbols for East Money API format.
// East Money uses secid format: "1.600519" (SH), "0.000001" (SZ)
// Market codes: 0=SZ, 1=SH
type EastMoneyFormatter struct{}

// eastMoneyMarketCode maps exchanges to East Money market codes.
var eastMoneyMarketCode = map[domain.Exchange]string{
	domain.ExchangeSH: "1",
	domain.ExchangeSZ: "0",
	domain.ExchangeBJ: "0",
}

// eastMoneyExchangeFromCode maps East Money market codes to exchanges.
var eastMoneyExchangeFromCode = map[string]domain.Exchange{
	"1": domain.ExchangeSH,
	"0": domain.ExchangeSZ,
}

func (f *EastMoneyFormatter) Format(sym domain.Symbol) string {
	switch sym.Market {
	case domain.MarketCN:
		marketCode, ok := eastMoneyMarketCode[sym.Exchange]
		if !ok {
			marketCode = "0" // default SZ
		}
		return marketCode + "." + sym.Code
	default:
		return sym.Code
	}
}

func (f *EastMoneyFormatter) Parse(raw string) (domain.Symbol, bool) {
	parts := strings.SplitN(raw, ".", 2)
	if len(parts) == 2 {
		marketCode := parts[0]
		code := parts[1]
		if len(code) == 6 && isPureDigits(code) {
			exchange, ok := eastMoneyExchangeFromCode[marketCode]
			if !ok {
				return domain.Symbol{}, false
			}
			return domain.Symbol{
				Code:      code,
				Market:    domain.MarketCN,
				Exchange:  exchange,
				AssetType: domain.AssetTypeStock,
				Standard:  fmt.Sprintf("%s.CN.%s", code, exchange),
			}, true
		}
	}
	return domain.Symbol{}, false
}
