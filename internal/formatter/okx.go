package formatter

import (
	"strings"

	"github.com/souloss/quantds/domain"
)

// OKXFormatter converts symbols for OKX API format.
// OKX uses dash-separated format: "BTC-USDT"
type OKXFormatter struct{}

func (f *OKXFormatter) Format(sym domain.Symbol) string {
	switch sym.Market {
	case domain.MarketCrypto:
		// Split BTCUSDT → BTC-USDT using quote suffix matching
		base, quote, ok := domain.MatchCryptoQuoteSuffix(sym.Code)
		if ok {
			return base + "-" + quote
		}
		// Fallback: if code contains dash already, return as-is
		if strings.Contains(sym.Code, "-") {
			return sym.Code
		}
		return sym.Code
	default:
		return sym.Code
	}
}

func (f *OKXFormatter) Parse(raw string) (domain.Symbol, bool) {
	upper := strings.ToUpper(raw)
	if strings.Contains(upper, "-") {
		parts := strings.SplitN(upper, "-", 2)
		if len(parts) == 2 {
			code := parts[0] + parts[1]
			return domain.Symbol{
				Code:      code,
				Market:    domain.MarketCrypto,
				Exchange:  domain.ExchangeOKX,
				AssetType: domain.AssetTypeCrypto,
				Standard:  code + ".CRYPTO.OKX",
			}, true
		}
	}
	return domain.Symbol{}, false
}
