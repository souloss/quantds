package binance

import (
	"fmt"
	"strings"
)

// Common quote assets (currencies)
const (
	QuoteUSDT = "USDT"
	QuoteBUSD = "BUSD"
	QuoteUSDC = "USDC"
	QuoteBTC  = "BTC"
	QuoteETH  = "ETH"
	QuoteBNB  = "BNB"
)

// Common base assets
const (
	BaseBTC  = "BTC"
	BaseETH  = "ETH"
	BaseBNB  = "BNB"
	BaseSOL  = "SOL"
	BaseXRP  = "XRP"
	BaseADA  = "ADA"
	BaseDOGE = "DOGE"
	BaseDOT  = "DOT"
)

// CommonTradingPairs contains the most popular trading pairs
var CommonTradingPairs = []string{
	"BTCUSDT", "ETHUSDT", "BNBUSDT", "SOLUSDT", "XRPUSDT",
	"ADAUSDT", "DOGEUSDT", "DOTUSDT", "MATICUSDT", "LTCUSDT",
	"BTCBUSD", "ETHBUSD", "BNBBUSD",
}

// ParseBinanceSymbol parses a Binance trading pair symbol
// Accepts formats: BTCUSDT, BTC-USDT, BTC/USDT
func ParseBinanceSymbol(symbol string) (base string, quote string, ok bool) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	// Handle BTC-USDT or BTC/USDT format
	if strings.Contains(symbol, "-") {
		parts := strings.Split(symbol, "-")
		if len(parts) == 2 {
			return parts[0], parts[1], true
		}
		return "", "", false
	}
	if strings.Contains(symbol, "/") {
		parts := strings.Split(symbol, "/")
		if len(parts) == 2 {
			return parts[0], parts[1], true
		}
		return "", "", false
	}

	// Handle standard Binance format (BTCUSDT)
	// Try to find known quote assets
	knownQuotes := []string{QuoteUSDT, QuoteBUSD, QuoteUSDC, QuoteBTC, QuoteETH, QuoteBNB}
	for _, q := range knownQuotes {
		if strings.HasSuffix(symbol, q) {
			base = symbol[:len(symbol)-len(q)]
			quote = q
			if len(base) > 0 {
				return base, quote, true
			}
		}
	}

	return "", "", false
}

// ToBinanceSymbol converts to standard Binance format (BTCUSDT)
func ToBinanceSymbol(symbol string) (string, error) {
	base, quote, ok := ParseBinanceSymbol(symbol)
	if !ok {
		// Already in correct format or unparseable
		if isBinanceFormat(symbol) {
			return symbol, nil
		}
		return "", fmt.Errorf("invalid Binance symbol: %s", symbol)
	}
	return base + quote, nil
}

// FromBinanceSymbol converts from Binance format to standard format
func FromBinanceSymbol(symbol string) string {
	base, quote, ok := ParseBinanceSymbol(symbol)
	if !ok {
		return symbol
	}
	return fmt.Sprintf("%s.%s.%s", base+quote, "CRYPTO", "BINANCE")
}

// isBinanceFormat checks if symbol is in Binance format (e.g., BTCUSDT)
func isBinanceFormat(symbol string) bool {
	_, _, ok := ParseBinanceSymbol(symbol)
	return ok
}

// IsCryptoSymbol checks if a symbol is a cryptocurrency
func IsCryptoSymbol(symbol string) bool {
	_, _, ok := ParseBinanceSymbol(symbol)
	if ok {
		return true
	}
	// Also check for CRYPTO.MARKET.EXCHANGE format
	if strings.Contains(symbol, ".CRYPTO.") {
		return true
	}
	return false
}

// GetQuoteAsset extracts the quote asset from a symbol
func GetQuoteAsset(symbol string) string {
	_, quote, ok := ParseBinanceSymbol(symbol)
	if !ok {
		return ""
	}
	return quote
}

// GetBaseAsset extracts the base asset from a symbol
func GetBaseAsset(symbol string) string {
	base, _, ok := ParseBinanceSymbol(symbol)
	if !ok {
		return ""
	}
	return base
}
