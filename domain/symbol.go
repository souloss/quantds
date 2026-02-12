package domain

import (
	"fmt"
	"strings"
)

// ========== 资产类型 ==========
type AssetType string

const (
	AssetTypeStock     AssetType = "STOCK"     // 股票
	AssetTypeCrypto    AssetType = "CRYPTO"    // 加密货币
	AssetTypeForex     AssetType = "FOREX"     // 外汇
	AssetTypeFutures   AssetType = "FUTURES"   // 期货
	AssetTypeOption    AssetType = "OPTION"    // 期权
	AssetTypeBond      AssetType = "BOND"      // 债券
	AssetTypeFund      AssetType = "FUND"      // 基金
	AssetTypeIndex     AssetType = "INDEX"     // 指数
	AssetTypeCommodity AssetType = "COMMODITY" // 大宗商品
)

// ========== 市场分类 ==========
type Market string

const (
	// 中国A股
	MarketCN Market = "CN"
	// 港股
	MarketHK Market = "HK"
	// 美股
	MarketUS Market = "US"
	// 加密货币（无国界，按交易所分）
	MarketCrypto Market = "CRYPTO"
	// 外汇（全球市场）
	MarketForex Market = "FOREX"
	// 期货（按交易所分）
	MarketFutures Market = "FUTURES"
)

type Exchange string

// A股交易所
const (
	ExchangeSH Exchange = "SH"
	ExchangeSZ Exchange = "SZ"
	ExchangeBJ Exchange = "BJ"
)

// 港股交易所
const (
	ExchangeHKEX Exchange = "HKEX" // 香港交易所
)

// 美股交易所
const (
	ExchangeNYSE   Exchange = "NYSE"   // 纽交所
	ExchangeNASDAQ Exchange = "NASDAQ" // 纳斯达克
	ExchangeAMEX   Exchange = "AMEX"   // 美交所
	ExchangeOTC    Exchange = "OTC"    // 场外交易
)

// 加密货币交易所
const (
	ExchangeBinance  Exchange = "BINANCE"  // 币安
	ExchangeCoinbase Exchange = "COINBASE" // Coinbase
	ExchangeOKX      Exchange = "OKX"      // OKX
	ExchangeBitget   Exchange = "BITGET"   // Bitget
)

// 外汇交易所/平台
const (
	ExchangeForexSpot Exchange = "FOREX_SPOT" // 外汇现货
	ExchangeForexCFD  Exchange = "FOREX_CFD"  // 外汇CFD
	ExchangeInterbank Exchange = "INTERBANK"  // 银行间市场
)

// 期货交易所
const (
	// 中国期货
	ExchangeSHFE  Exchange = "SHFE"  // 上期所
	ExchangeDCE   Exchange = "DCE"   // 大商所
	ExchangeCZCE  Exchange = "CZCE"  // 郑商所
	ExchangeCFFEX Exchange = "CFFEX" // 中金所
	ExchangeINE   Exchange = "INE"   // 上期能源
	ExchangeGFEX  Exchange = "GFEX"  // 广期所
)

// ========== 市场元数据配置 ==========
// 市场配置
type MarketConfig struct {
	AssetType       AssetType
	Market          Market
	DefaultExchange Exchange
	CodeRules       CodeRules // 代码规则
}

// 代码规则
type CodeRules struct {
	MinLength    int      // 最小长度
	MaxLength    int      // 最大长度
	AllowLetters bool     // 是否允许字母
	AllowDigits  bool     // 是否允许数字
	Prefixes     []string // 特定前缀（如港股0开头）
	Suffix       string   // 固定后缀
}

// 市场配置表（可扩展）
var MarketConfigs = map[Market]MarketConfig{
	MarketCN: {
		AssetType:       AssetTypeStock,
		Market:          MarketCN,
		DefaultExchange: ExchangeSZ,
		CodeRules: CodeRules{
			MinLength:    6,
			MaxLength:    6,
			AllowLetters: false,
			AllowDigits:  true,
		},
	},
	MarketHK: {
		AssetType:       AssetTypeStock,
		Market:          MarketHK,
		DefaultExchange: ExchangeHKEX,
		CodeRules: CodeRules{
			MinLength:    5,
			MaxLength:    5,
			AllowLetters: true,
			AllowDigits:  true,
		},
	},
	MarketUS: {
		AssetType:       AssetTypeStock,
		Market:          MarketUS,
		DefaultExchange: ExchangeNASDAQ,
		CodeRules: CodeRules{
			MinLength:    1,
			MaxLength:    5,
			AllowLetters: true,
			AllowDigits:  true,
		},
	},
	MarketCrypto: {
		AssetType:       AssetTypeCrypto,
		Market:          MarketCrypto,
		DefaultExchange: ExchangeBinance,
		CodeRules: CodeRules{
			MinLength:    2,
			MaxLength:    10,
			AllowLetters: true,
			AllowDigits:  true,
		},
	},
	MarketForex: {
		AssetType:       AssetTypeForex,
		Market:          MarketForex,
		DefaultExchange: ExchangeForexSpot,
		CodeRules: CodeRules{
			MinLength:    6,
			MaxLength:    6,
			AllowLetters: true,
			AllowDigits:  false,
			Suffix:       "", // 如 EURUSD
		},
	},
	MarketFutures: {
		AssetType:       AssetTypeFutures,
		Market:          MarketFutures,
		DefaultExchange: ExchangeSHFE,
		CodeRules: CodeRules{
			MinLength:    3,
			MaxLength:    20,
			AllowLetters: true,
			AllowDigits:  true, // 包含到期月份
		},
	},
}

// ========== Symbol 结构体 ==========
type Symbol struct {
	Code      string    // 原始代码（如 000001, AAPL, BTCUSDT）
	Market    Market    // 市场（CN, US, CRYPTO）
	Exchange  Exchange  // 交易所（SZ, NASDAQ, BINANCE）
	AssetType AssetType // 资产类型（推导得出）
	Standard  string    // 标准化格式：CODE.MARKET.EXCHANGE
}

// Parse 解析各种格式的代码
func (s *Symbol) Parse(input string) error {
	input = strings.TrimSpace(strings.ToUpper(input))

	// 尝试解析三级格式：CODE.MARKET.EXCHANGE
	if parts := strings.Split(input, "."); len(parts) == 3 {
		s.Code = parts[0]
		s.Market = Market(parts[1])
		s.Exchange = Exchange(parts[2])
		s.AssetType = deriveAssetType(s.Market)
		s.Standard = input
		return s.Validate()
	}

	// 尝试解析二级格式：CODE.EXCHANGE（兼容旧A股）
	if parts := strings.Split(input, "."); len(parts) == 2 {
		s.Code = parts[0]
		s.Exchange = Exchange(parts[1])
		s.Market = deriveMarketFromExchange(s.Exchange)
		s.AssetType = deriveAssetType(s.Market)
		s.Standard = fmt.Sprintf("%s.%s.%s", s.Code, s.Market, s.Exchange)
		return s.Validate()
	}

	// 尝试解析前缀格式：EXCHANGE+CODE（如 SH000001）
	if len(input) >= 8 {
		prefix := input[:2]
		code := input[2:]
		switch prefix {
		case "SH", "SZ", "BJ":
			s.Code = code
			s.Exchange = Exchange(prefix)
			s.Market = MarketCN
			s.AssetType = AssetTypeStock
			s.Standard = fmt.Sprintf("%s.%s.%s", s.Code, s.Market, s.Exchange)
			return s.Validate()
		}
	}

	// 尝试智能识别（无后缀）
	return s.SmartParse(input)
}

// SmartParse 智能识别无后缀代码
func (s *Symbol) SmartParse(code string) error {
	s.Code = strings.ToUpper(code)

	// A股：6位纯数字
	if isPureDigits(code) && len(code) == 6 {
		s.Market = MarketCN
		s.Exchange = inferCNExchange(code) // 根据代码规则推断
		s.AssetType = AssetTypeStock
		s.Standard = fmt.Sprintf("%s.%s.%s", s.Code, s.Market, s.Exchange)
		return nil
	}

	// 港股：5位，0开头或字母开头
	if len(code) == 5 && (code[0] == '0' || isLetter(code[0])) {
		s.Market = MarketHK
		s.Exchange = ExchangeHKEX
		s.AssetType = AssetTypeStock
		s.Standard = fmt.Sprintf("%s.%s.%s", s.Code, s.Market, s.Exchange)
		return nil
	}

	// 美股：1-5位字母
	if len(code) >= 1 && len(code) <= 5 && isAllLetters(code) {
		s.Market = MarketUS
		s.Exchange = ExchangeNASDAQ // 默认，可后续修正
		s.AssetType = AssetTypeStock
		s.Standard = fmt.Sprintf("%s.%s.%s", s.Code, s.Market, s.Exchange)
		return nil
	}

	// 加密货币：常见格式（BTCUSDT, BTC-USD）
	cryptoCode := normalizeCryptoCode(code)
	if isCryptoFormat(cryptoCode) {
		s.Code = cryptoCode
		s.Market = MarketCrypto
		s.Exchange = ExchangeBinance // 默认
		s.AssetType = AssetTypeCrypto
		s.Standard = fmt.Sprintf("%s.%s.%s", s.Code, s.Market, s.Exchange)
		return nil
	}

	// 外汇：6位字母（EURUSD）
	if len(code) == 6 && isAllLetters(code) {
		s.Market = MarketForex
		s.Exchange = ExchangeForexSpot
		s.AssetType = AssetTypeForex
		s.Standard = fmt.Sprintf("%s.%s.%s", s.Code, s.Market, s.Exchange)
		return nil
	}

	return fmt.Errorf("无法识别代码格式: %s", code)
}

// Validate 验证代码合法性
func (s *Symbol) Validate() error {
	config, ok := MarketConfigs[s.Market]
	if !ok {
		return fmt.Errorf("未知市场: %s", s.Market)
	}

	rules := config.CodeRules

	// 长度检查
	if len(s.Code) < rules.MinLength || len(s.Code) > rules.MaxLength {
		return fmt.Errorf("代码长度非法: %s (需%d-%d位)", s.Code, rules.MinLength, rules.MaxLength)
	}

	// 字符检查
	for _, c := range s.Code {
		if !rules.AllowLetters && isLetter(byte(c)) {
			return fmt.Errorf("代码包含非法字母: %s", s.Code)
		}
		if !rules.AllowDigits && isDigit(byte(c)) {
			return fmt.Errorf("代码包含非法数字: %s", s.Code)
		}
	}

	return nil
}

// String 返回标准化字符串
func (s Symbol) String() string {
	return s.Standard
}

// Short 返回短格式（兼容旧系统）
func (s Symbol) Short() string {
	// A股保持 000001.SZ 格式
	if s.Market == MarketCN {
		return fmt.Sprintf("%s.%s", s.Code, s.Exchange)
	}
	return s.Standard
}

// ========== 辅助函数 ==========

func deriveAssetType(market Market) AssetType {
	if config, ok := MarketConfigs[market]; ok {
		return config.AssetType
	}
	return AssetTypeStock
}

func deriveMarketFromExchange(ex Exchange) Market {
	switch ex {
	case ExchangeSH, ExchangeSZ, ExchangeBJ:
		return MarketCN
	case ExchangeHKEX:
		return MarketHK
	case ExchangeNYSE, ExchangeNASDAQ, ExchangeAMEX, ExchangeOTC:
		return MarketUS
	default:
		if strings.HasPrefix(string(ex), "BINANCE") ||
			strings.HasPrefix(string(ex), "COINBASE") ||
			strings.HasPrefix(string(ex), "OKX") {
			return MarketCrypto
		}
		if strings.HasPrefix(string(ex), "SHFE") ||
			strings.HasPrefix(string(ex), "CME") {
			return MarketFutures
		}
		return MarketUS // 默认
	}
}

func inferCNExchange(code string) Exchange {
	// 上交所：60, 68, 90, 89 开头
	// 深交所：00, 30, 20 开头
	// 北交所：43, 83, 87 开头
	if len(code) >= 2 {
		prefix := code[:2]
		switch {
		case prefix == "60" || prefix == "68" || prefix == "90":
			return ExchangeSH
		case prefix == "00" || prefix == "30" || prefix == "20":
			return ExchangeSZ
		case prefix == "43" || prefix == "83" || prefix == "87":
			return ExchangeBJ
		}
	}
	return ExchangeSZ // 默认深交所
}

func isPureDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isAllLetters(s string) bool {
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return false
		}
	}
	return true
}

func isLetter(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func normalizeCryptoCode(code string) string {
	// 统一格式：去除分隔符，大写
	code = strings.ToUpper(code)
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, "/", "")
	code = strings.ReplaceAll(code, "_", "")
	return code
}

func isCryptoFormat(code string) bool {
	// 简单启发式：包含常见后缀
	commonBases := []string{"BTC", "ETH", "USDT", "USDC", "USD", "BUSD"}
	commonQuotes := []string{"USDT", "USDC", "USD", "BUSD", "BTC", "ETH"}

	for _, base := range commonBases {
		for _, quote := range commonQuotes {
			if code == base+quote || code == quote+base {
				return true
			}
		}
	}
	return false
}

// ========== 便捷解析函数（兼容旧接口）==========

func ParseSymbol(input string) (code string, exchange Exchange, ok bool) {
	var s Symbol
	if err := s.Parse(input); err != nil {
		return "", "", false
	}
	return s.Code, s.Exchange, true
}

func FormatSymbol(code string, exchange Exchange) string {
	market := deriveMarketFromExchange(exchange)
	return fmt.Sprintf("%s.%s.%s", code, market, exchange)
}

// 扩展：完整格式化
func FormatFullSymbol(code string, market Market, exchange Exchange) string {
	return fmt.Sprintf("%s.%s.%s", code, market, exchange)
}
