package formatter

import (
	"strings"

	"github.com/souloss/quantds/domain"
)

// BinanceFormatter converts symbols for Binance API format.
// Binance uses concatenated format: "BTCUSDT"
// The base and quote are determined by known quote suffix matching.
type BinanceFormatter struct{}

func (f *BinanceFormatter) Format(sym domain.Symbol) string {
	switch sym.Market {
	case domain.MarketCrypto:
		// Symbol code is already in BTCUSDT format from SmartParse
		return sym.Code
	default:
		return sym.Code
	}
}

func (f *BinanceFormatter) Parse(raw string) (domain.Symbol, bool) {
	// Handle BTC-USDT or BTC/USDT format
	upper := strings.ToUpper(raw)
	if strings.Contains(upper, "-") || strings.Contains(upper, "/") {
		var sep string
		if strings.Contains(upper, "-") {
			sep = "-"
		} else {
			sep = "/"
		}
		parts := strings.SplitN(upper, sep, 2)
		if len(parts) == 2 {
			code := parts[0] + parts[1]
			return domain.Symbol{
				Code:      code,
				Market:    domain.MarketCrypto,
				Exchange:  domain.ExchangeBinance,
				AssetType: domain.AssetTypeCrypto,
				Standard:  code + ".CRYPTO.BINANCE",
			}, true
		}
	}

	// Handle standard Binance format (BTCUSDT)
	// Use quote suffix matching
	base, quote, ok := domain.MatchCryptoQuoteSuffix(upper)
	if ok {
		return domain.Symbol{
			Code:      base + quote,
			Market:    domain.MarketCrypto,
			Exchange:  domain.ExchangeBinance,
			AssetType: domain.AssetTypeCrypto,
			Standard:  base + quote + ".CRYPTO.BINANCE",
		}, true
	}
	return domain.Symbol{}, false
}
