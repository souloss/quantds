package symbolname

// cryptoNames maps crypto base symbols to their human-readable names.
var cryptoNames = map[string]string{
	// Major
	"BTC":   "Bitcoin",
	"ETH":   "Ethereum",
	"BNB":   "BNB",
	"SOL":   "Solana",
	"XRP":   "XRP",
	"ADA":   "Cardano",
	"DOGE":  "Dogecoin",
	"DOT":   "Polkadot",
	"MATIC": "Polygon",
	"AVAX":  "Avalanche",
	// DeFi
	"UNI":   "Uniswap",
	"LINK":  "Chainlink",
	"AAVE":  "Aave",
	"MKR":   "Maker",
	"COMP":  "Compound",
	"CRV":   "Curve DAO",
	"SNX":   "Synthetix",
	"SUSHI": "SushiSwap",
	"1INCH": "1inch",
	"YFI":   "yearn.finance",
	// Layer 2
	"ARB":  "Arbitrum",
	"OP":   "Optimism",
	"STRK": "StarkNet",
	"IMX":  "Immutable",
	// Smart Contract Platforms
	"NEAR": "NEAR Protocol",
	"ALGO": "Algorand",
	"ATOM": "Cosmos",
	"FTM":  "Fantom",
	"ICP":  "Internet Computer",
	"HBAR": "Hedera",
	"SEI":  "Sei",
	"SUI":  "Sui",
	"APT":  "Aptos",
	"TIA":  "Celestia",
	// Meme
	"SHIB":  "Shiba Inu",
	"PEPE":  "Pepe",
	"FLOKI": "Floki",
	"WIF":   "dogwifhat",
	"BONK":  "Bonk",
	// Exchange
	"OKB": "OKB",
	"LEO": "UNUS SED LEO",
	"CRO": "Cronos",
	"GT":  "GateToken",
	// Stablecoins
	"USDT": "Tether",
	"USDC": "USD Coin",
	"DAI":  "Dai",
	"BUSD": "Binance USD",
	"TUSD": "TrueUSD",
	// Others
	"LTC":    "Litecoin",
	"BCH":    "Bitcoin Cash",
	"ETC":    "Ethereum Classic",
	"FIL":    "Filecoin",
	"VET":    "VeChain",
	"TRX":    "TRON",
	"EOS":    "EOS",
	"XTZ":    "Tezos",
	"THETA":  "Theta Network",
	"XLM":    "Stellar",
	"XMR":    "Monero",
	"RUNE":   "THORChain",
	"GRT":    "The Graph",
	"INJ":    "Injective",
	"TON":    "Toncoin",
	"KAS":    "Kaspa",
	"JUP":    "Jupiter",
	"ENA":    "Ethena",
	"PENDLE": "Pendle",
	"RENDER": "Render",
	"FET":    "Fetch.ai",
	"RNDR":   "Render",
	"WLD":    "Worldcoin",
	"STX":    "Stacks",
}

// resolveCryptoName resolves a crypto symbol to its human-readable name.
// It handles both base symbols (e.g., "BTC") and trading pairs (e.g., "BTCUSDT").
func resolveCryptoName(symbol string) string {
	// Try direct lookup first
	if name, ok := cryptoNames[symbol]; ok {
		return name
	}

	// Try stripping common quote currencies from trading pairs
	quoteCurrencies := []string{"USDT", "BUSD", "USDC", "BTC", "ETH", "BNB"}
	for _, quote := range quoteCurrencies {
		if len(symbol) > len(quote) && symbol[len(symbol)-len(quote):] == quote {
			base := symbol[:len(symbol)-len(quote)]
			if name, ok := cryptoNames[base]; ok {
				return name + "/" + quote
			}
		}
	}

	return symbol
}
