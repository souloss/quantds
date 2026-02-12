package yahoo

import (
	"fmt"
	"strings"
)

// USExchangePrefix maps exchange to market ID for symbol formatting
var USExchangePrefix = map[string]string{
	"NYSE":   "NYSE",
	"NASDAQ": "NASDAQ",
	"AMEX":   "AMEX",
	"OTC":    "OTC",
}

// USExchangeNames maps exchange code to display name
var USExchangeNames = map[string]string{
	"NYSE":   "New York Stock Exchange",
	"NASDAQ": "NASDAQ Stock Market",
	"AMEX":   "American Stock Exchange",
	"OTC":    "Over-the-Counter",
}

// ParseUSSymbol parses a US stock symbol string
// Accepts formats: AAPL, AAPL.US, AAPL.NASDAQ, AAPL.NYSE
func ParseUSSymbol(symbol string) (code string, exchange string, ok bool) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))

	// Handle CODE.EXCHANGE format
	if strings.Contains(symbol, ".") {
		parts := strings.Split(symbol, ".")
		if len(parts) == 2 || len(parts) == 3 {
			code = parts[0]
			if len(parts) == 2 {
				// Format: CODE.US or CODE.EXCHANGE
				exchange = parts[1]
				if exchange == "US" {
					exchange = "NASDAQ" // Default US exchange
				}
			} else if len(parts) == 3 {
				// Format: CODE.US.EXCHANGE
				exchange = parts[2]
			}
			return code, exchange, true
		}
		return "", "", false
	}

	// Plain symbol - assume NASDAQ
	return symbol, "NASDAQ", true
}

// ToYahooSymbol converts a symbol to Yahoo Finance format
// Yahoo uses plain symbols like "AAPL" for US stocks
func ToYahooSymbol(symbol string) (string, error) {
	code, _, ok := ParseUSSymbol(symbol)
	if !ok {
		return "", fmt.Errorf("invalid US stock symbol: %s", symbol)
	}
	return code, nil
}

// FromYahooSymbol converts Yahoo symbol to standard format
func FromYahooSymbol(symbol string, exchange string) string {
	if exchange == "" {
		exchange = "NASDAQ"
	}
	return fmt.Sprintf("%s.US.%s", symbol, exchange)
}

// IsUSSymbol checks if a symbol is a US stock
func IsUSSymbol(symbol string) bool {
	code, _, ok := ParseUSSymbol(symbol)
	if !ok {
		return false
	}
	// US stock symbols are 1-5 letters, sometimes with special chars
	if len(code) < 1 || len(code) > 5 {
		// Some symbols can be longer (e.g., BRK.B)
		if len(code) > 6 {
			return false
		}
	}
	// Must start with a letter
	if len(code) > 0 && !isLetter(code[0]) {
		return false
	}
	return true
}

func isLetter(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}
