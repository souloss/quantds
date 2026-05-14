package symbolname

// forexNames maps forex pair codes to their human-readable names.
// Keys use the standard 6-letter format (e.g., "EURUSD").
var forexNames = map[string]string{
	// Major pairs
	"EURUSD": "Euro / US Dollar",
	"GBPUSD": "British Pound / US Dollar",
	"USDJPY": "US Dollar / Japanese Yen",
	"USDCHF": "US Dollar / Swiss Franc",
	"AUDUSD": "Australian Dollar / US Dollar",
	"USDCAD": "US Dollar / Canadian Dollar",
	"NZDUSD": "New Zealand Dollar / US Dollar",
	// Cross pairs
	"EURGBP": "Euro / British Pound",
	"EURJPY": "Euro / Japanese Yen",
	"GBPJPY": "British Pound / Japanese Yen",
	"EURCHF": "Euro / Swiss Franc",
	"EURAUD": "Euro / Australian Dollar",
	"EURNZD": "Euro / New Zealand Dollar",
	"EURCAD": "Euro / Canadian Dollar",
	"GBPAUD": "British Pound / Australian Dollar",
	"GBPNZD": "British Pound / New Zealand Dollar",
	"GBPCHF": "British Pound / Swiss Franc",
	"GBPCAD": "British Pound / Canadian Dollar",
	"AUDJPY": "Australian Dollar / Japanese Yen",
	"AUDNZD": "Australian Dollar / New Zealand Dollar",
	"AUDCAD": "Australian Dollar / Canadian Dollar",
	"AUDCHF": "Australian Dollar / Swiss Franc",
	"NZDJPY": "New Zealand Dollar / Japanese Yen",
	"NZDCAD": "New Zealand Dollar / Canadian Dollar",
	"NZDCHF": "New Zealand Dollar / Swiss Franc",
	"CADJPY": "Canadian Dollar / Japanese Yen",
	"CADCHF": "Canadian Dollar / Swiss Franc",
	"CHFJPY": "Swiss Franc / Japanese Yen",
	// Exotic pairs
	"USDTRY": "US Dollar / Turkish Lira",
	"USDZAR": "US Dollar / South African Rand",
	"USDMXN": "US Dollar / Mexican Peso",
	"USDSGD": "US Dollar / Singapore Dollar",
	"USDHKD": "US Dollar / Hong Kong Dollar",
	"USDNOK": "US Dollar / Norwegian Krone",
	"USDSEK": "US Dollar / Swedish Krona",
	"USDDKK": "US Dollar / Danish Krone",
	"USDPLN": "US Dollar / Polish Zloty",
	"USDCZK": "US Dollar / Czech Koruna",
	"USDHUF": "US Dollar / Hungarian Forint",
	"USDCNH": "US Dollar / Chinese Yuan (Offshore)",
	// CNY pairs
	"USDCNY": "US Dollar / Chinese Yuan",
	"EURCNY": "Euro / Chinese Yuan",
	"GBPCNY": "British Pound / Chinese Yuan",
	"JPYCNY": "Japanese Yen / Chinese Yuan",
}

// currencyNames maps ISO 4217 currency codes to their names.
var currencyNames = map[string]string{
	"USD": "US Dollar",
	"EUR": "Euro",
	"GBP": "British Pound",
	"JPY": "Japanese Yen",
	"CHF": "Swiss Franc",
	"AUD": "Australian Dollar",
	"CAD": "Canadian Dollar",
	"NZD": "New Zealand Dollar",
	"CNY": "Chinese Yuan",
	"CNH": "Chinese Yuan (Offshore)",
	"HKD": "Hong Kong Dollar",
	"SGD": "Singapore Dollar",
	"TRY": "Turkish Lira",
	"ZAR": "South African Rand",
	"MXN": "Mexican Peso",
	"NOK": "Norwegian Krone",
	"SEK": "Swedish Krona",
	"DKK": "Danish Krone",
	"PLN": "Polish Zloty",
	"CZK": "Czech Koruna",
	"HUF": "Hungarian Forint",
	"KRW": "South Korean Won",
	"INR": "Indian Rupee",
	"RUB": "Russian Ruble",
	"BRL": "Brazilian Real",
	"THB": "Thai Baht",
}

// resolveForexName resolves a forex pair code to its human-readable name.
// It handles 6-letter codes (e.g., "EURUSD") by splitting into base/quote.
func resolveForexName(code string) string {
	// Try direct lookup first
	if name, ok := forexNames[code]; ok {
		return name
	}

	// Try splitting 6-letter code into base + quote (e.g., "EURUSD" -> "EUR" + "USD")
	if len(code) == 6 {
		base := code[:3]
		quote := code[3:]
		baseName, baseOk := currencyNames[base]
		quoteName, quoteOk := currencyNames[quote]
		if baseOk && quoteOk {
			return baseName + " / " + quoteName
		}
		if baseOk {
			return baseName + " / " + quote
		}
		if quoteOk {
			return base + " / " + quoteName
		}
	}

	return code
}
