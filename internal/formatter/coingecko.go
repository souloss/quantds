package formatter

import (
	"strings"

	"github.com/souloss/quantds/domain"
)

// CoinGeckoFormatter converts symbols for CoinGecko API format.
// CoinGecko uses lowercase coin IDs: "bitcoin", "ethereum"
//
// Unlike other formatters that use deterministic rules, CoinGeFormatter
// requires dynamic lookup because coin IDs are not derivable from
// ticker symbols. The formatter provides:
// - Format: converts Symbol.Code to lowercase as a best-effort guess
// - Parse: converts lowercase coin ID back (limited, needs API for full support)
//
// For production use, the adapter should use CoinGecko's /search API
// to resolve symbols to coin IDs dynamically.
type CoinGeckoFormatter struct{}

// Common ticker-to-coinID mappings for the most popular coins.
// This is NOT a comprehensive mapping table — it's a small cache for
// the most frequently used coins. For anything not in this list,
// the adapter should use CoinGecko's /search API.
var coinGeckoCommonIDs = map[string]string{
	"BTC":   "bitcoin",
	"ETH":   "ethereum",
	"BNB":   "binancecoin",
	"SOL":   "solana",
	"XRP":   "ripple",
	"ADA":   "cardano",
	"DOGE":  "dogecoin",
	"DOT":   "polkadot",
	"MATIC": "matic-network",
	"LTC":   "litecoin",
	"AVAX":  "avalanche-2",
	"LINK":  "chainlink",
	"UNI":   "uniswap",
	"ATOM":  "cosmos",
	"NEAR":  "near",
}

// Reverse mapping: coin ID → ticker
var coinGeckoReverseIDs map[string]string

func init() {
	coinGeckoReverseIDs = make(map[string]string, len(coinGeckoCommonIDs))
	for k, v := range coinGeckoCommonIDs {
		coinGeckoReverseIDs[v] = k
	}
}

func (f *CoinGeckoFormatter) Format(sym domain.Symbol) string {
	switch sym.Market {
	case domain.MarketCrypto:
		// Extract base from code (e.g., BTCUSDT → BTC)
		base, _, ok := domain.MatchCryptoQuoteSuffix(sym.Code)
		if !ok {
			base = sym.Code
		}
		// Look up coin ID
		if id, ok := coinGeckoCommonIDs[strings.ToUpper(base)]; ok {
			return id
		}
		// Fallback: lowercase the base
		return strings.ToLower(base)
	default:
		return strings.ToLower(sym.Code)
	}
}

func (f *CoinGeckoFormatter) Parse(raw string) (domain.Symbol, bool) {
	lower := strings.ToLower(raw)
	if ticker, ok := coinGeckoReverseIDs[lower]; ok {
		return domain.Symbol{
			Code:      ticker,
			Market:    domain.MarketCrypto,
			Exchange:  domain.ExchangeBinance,
			AssetType: domain.AssetTypeCrypto,
			Standard:  ticker + ".CRYPTO.BINANCE",
		}, true
	}
	return domain.Symbol{}, false
}
