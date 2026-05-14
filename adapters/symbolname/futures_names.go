package symbolname

// futuresNames maps traditional futures contract codes to their human-readable names.
// For codes that conflict across exchanges (e.g., "ZN" on CME vs SHFE),
// the most commonly used interpretation is kept as default.
// Exchange-specific lookups are handled by resolveFuturesNameWithExchange.
var futuresNames = map[string]string{
	// CME Group - Equity Index
	"ES":  "E-mini S&P 500",
	"NQ":  "E-mini NASDAQ-100",
	"YM":  "E-mini Dow Jones",
	"RTY": "E-mini Russell 2000",
	"ND":  "NASDAQ-100",
	// CME Group - Interest Rate
	"ZN":  "10-Year T-Note",
	"ZB":  "30-Year T-Bond",
	"ZF":  "5-Year T-Note",
	"ZT":  "2-Year T-Note",
	"GE":  "Eurodollar",
	"FF":  "Federal Funds",
	"SR3": "SOFR 3-Month",
	// CME Group - Energy
	"CL": "Crude Oil (WTI)",
	"NG": "Natural Gas",
	"HO": "Heating Oil",
	// CME Group - Metals
	"GC": "Gold",
	"SI": "Silver",
	"HG": "Copper",
	"PL": "Platinum",
	"PA": "Palladium",
	// CME Group - Agriculture
	"ZC":  "Corn",
	"ZS":  "Soybeans",
	"ZW":  "Wheat",
	"ZL":  "Soybean Oil",
	"ZM":  "Soybean Meal",
	"KC":  "Coffee",
	"SB":  "Sugar #11",
	"CC":  "Cocoa",
	"CT":  "Cotton",
	"OJ":  "Orange Juice",
	"LBS": "Lumber",
	// ICE
	"BRN": "Brent Crude",
	// Shanghai (SHFE)
	"CU": "沪铜",
	"AL": "沪铝",
	"PB": "沪铅",
	"NI": "沪镍",
	"SN": "沪锡",
	"AU": "沪金",
	"AG": "沪银",
	"HC": "热卷",
	"SS": "不锈钢",
	"BU": "沥青",
	"RU": "橡胶",
	"NR": "20号胶",
	"FU": "燃油",
	"SC": "原油",
	"BC": "国际铜",
	"LU": "低硫燃油",
	// Dalian (DCE)
	"I":  "铁矿石",
	"J":  "焦炭",
	"JM": "焦煤",
	"A":  "豆一",
	"B":  "豆二",
	"M":  "豆粕",
	"Y":  "豆油",
	"P":  "棕榈油",
	"CS": "玉米淀粉",
	"JD": "鸡蛋",
	"L":  "聚乙烯",
	"V":  "聚氯乙烯",
	"PP": "聚丙烯",
	"EG": "乙二醇",
	"EB": "苯乙烯",
	"PG": "液化石油气",
	"FB": "纤维板",
	"BB": "胶合板",
	"RR": "粳米",
	"LH": "生猪",
	// Zhengzhou (CZCE)
	"CF": "棉花",
	"SR": "白糖",
	"TA": "PTA",
	"MA": "甲醇",
	"FG": "玻璃",
	"SA": "纯碱",
	"OI": "菜油",
	"RM": "菜粕",
	"AP": "苹果",
	"CJ": "红枣",
	"PF": "短纤",
	"UR": "尿素",
	"SH": "烧碱",
	// CFFEX
	"IF": "沪深300指数",
	"IC": "中证500指数",
	"IM": "中证1000指数",
	"IH": "上证50指数",
	"TF": "5年期国债",
	"TS": "2年期国债",
	"TL": "30年期国债",
}

// exchangeFuturesNames provides exchange-specific overrides for codes that conflict.
// The key format is "EXCHANGE:CODE" (e.g., "SHFE:ZN" for 沪锌).
var exchangeFuturesNames = map[string]string{
	// SHFE overrides (conflict with CME codes)
	"SHFE:ZN": "沪锌",
	"SHFE:RB": "螺纹钢",
	"SHFE:SP": "纸浆",
	"SHFE:C":  "螺纹钢(废)",
	// CME overrides (to disambiguate from SHFE)
	"CME:RB": "RBOB Gasoline",
	"CME:SP": "S&P 500",
	"CME:C":  "Corn",
}

// resolveFuturesName resolves a futures contract code to its human-readable name.
// It handles both traditional futures codes and crypto perpetual futures.
func resolveFuturesName(code string) string {
	// Try direct lookup first
	if name, ok := futuresNames[code]; ok {
		return name
	}

	// Try stripping common suffixes for traditional futures (e.g., "ES1!", "ES=F", "ESZ24")
	stripped := stripFuturesSuffix(code)
	if stripped != code {
		if name, ok := futuresNames[stripped]; ok {
			return name
		}
	}

	// Try crypto perpetual futures (e.g., "BTCUSDT" → "Bitcoin Perpetual")
	if name := resolveCryptoPerpetualName(code); name != "" {
		return name
	}

	return code
}

// resolveFuturesNameWithExchange resolves a futures contract code with exchange context.
func resolveFuturesNameWithExchange(code, exchange string) string {
	// Try exchange-specific lookup first
	if exchange != "" {
		key := exchange + ":" + code
		if name, ok := exchangeFuturesNames[key]; ok {
			return name
		}
	}

	return resolveFuturesName(code)
}

// stripFuturesSuffix removes trailing contract identifiers from futures codes.
// e.g., "ES1!" → "ES", "ES=F" → "ES", "ESZ24" → "ES"
func stripFuturesSuffix(code string) string {
	// Handle "=F" or "=X" suffix (Yahoo Finance format)
	if idx := len(code) - 2; idx > 0 {
		suffix := code[idx:]
		if suffix == "=F" || suffix == "=X" {
			return code[:idx]
		}
	}

	// Handle trailing "!" (TradingView continuous format)
	if len(code) > 1 && code[len(code)-1] == '!' {
		return code[:len(code)-1]
	}

	// Handle trailing digits (e.g., "ES1", "ESZ24")
	i := len(code) - 1
	for i >= 0 && (code[i] >= '0' && code[i] <= '9') {
		i--
	}
	if i < len(code)-1 && i > 0 {
		return code[:i+1]
	}

	return code
}

// resolveCryptoPerpetualName resolves crypto perpetual futures names.
// e.g., "BTCUSDT" → "Bitcoin Perpetual"
func resolveCryptoPerpetualName(symbol string) string {
	quoteCurrencies := []string{"USDT", "BUSD", "USDC", "BTC", "ETH"}
	for _, quote := range quoteCurrencies {
		if len(symbol) > len(quote) && symbol[len(symbol)-len(quote):] == quote {
			base := symbol[:len(symbol)-len(quote)]
			if name, ok := cryptoNames[base]; ok {
				return name + " Perpetual"
			}
		}
	}
	return ""
}
