package formatter

import (
	"fmt"
	"strings"

	"github.com/souloss/quantds/domain"
)

// EODHDFormatter converts symbols for EODHD API format.
// EODHD uses CODE.EXCHANGE format: "AAPL.NYSE", "600519.SH"
type EODHDFormatter struct{}

var eodhdExchangeMap = map[domain.Exchange]string{
	// US
	domain.ExchangeNYSE:   "NYSE",
	domain.ExchangeNASDAQ: "NASDAQ",
	// CN
	domain.ExchangeSH: "SH",
	domain.ExchangeSZ: "SZ",
	domain.ExchangeBJ: "BJ",
	// HK
	domain.ExchangeHKEX: "HK",
}

var eodhdExchangeFromStr = map[string]domain.Exchange{
	"NYSE":   domain.ExchangeNYSE,
	"NASDAQ": domain.ExchangeNASDAQ,
	"SH":     domain.ExchangeSH,
	"SZ":     domain.ExchangeSZ,
	"BJ":     domain.ExchangeBJ,
	"HK":     domain.ExchangeHKEX,
}

func (f *EODHDFormatter) Format(sym domain.Symbol) string {
	suffix, ok := eodhdExchangeMap[sym.Exchange]
	if ok {
		return sym.Code + "." + suffix
	}
	return sym.Code
}

func (f *EODHDFormatter) Parse(raw string) (domain.Symbol, bool) {
	parts := strings.SplitN(raw, ".", 2)
	if len(parts) == 2 {
		code := parts[0]
		exchangeStr := parts[1]
		exchange, ok := eodhdExchangeFromStr[exchangeStr]
		if !ok {
			return domain.Symbol{}, false
		}
		var market domain.Market
		switch {
		case exchange == domain.ExchangeSH || exchange == domain.ExchangeSZ || exchange == domain.ExchangeBJ:
			market = domain.MarketCN
		case exchange == domain.ExchangeNYSE || exchange == domain.ExchangeNASDAQ:
			market = domain.MarketUS
		case exchange == domain.ExchangeHKEX:
			market = domain.MarketHK
		default:
			return domain.Symbol{}, false
		}
		return domain.Symbol{
			Code:      code,
			Market:    market,
			Exchange:  exchange,
			AssetType: domain.AssetTypeStock,
			Standard:  fmt.Sprintf("%s.%s.%s", code, market, exchange),
		}, true
	}
	return domain.Symbol{}, false
}
