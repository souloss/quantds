package domain

import (
	"testing"
)

func TestParseSymbol(t *testing.T) {
	tests := []struct {
		symbol       string
		wantCode     string
		wantMarket   Market
		wantExchange Exchange
		wantOK       bool
	}{
		// A股二级格式
		{"000001.SZ", "000001", MarketCN, ExchangeSZ, true},
		{"600001.SH", "600001", MarketCN, ExchangeSH, true},
		{"430001.BJ", "430001", MarketCN, ExchangeBJ, true},
		// A股前缀格式
		{"SH600001", "600001", MarketCN, ExchangeSH, true},
		{"SZ000001", "000001", MarketCN, ExchangeSZ, true},
		{"BJ430001", "430001", MarketCN, ExchangeBJ, true},
		// A股三级格式
		{"000001.CN.SZ", "000001", MarketCN, ExchangeSZ, true},
		{"600001.CN.SH", "600001", MarketCN, ExchangeSH, true},
		// 美股三级格式
		{"AAPL.US.NASDAQ", "AAPL", MarketUS, ExchangeNASDAQ, true},
		// 加密货币三级格式
		{"BTCUSDT.CRYPTO.BINANCE", "BTCUSDT", MarketCrypto, ExchangeBinance, true},
		// 外汇三级格式
		{"EURUSD.FOREX.FOREX_SPOT", "EURUSD", MarketForex, ExchangeForexSpot, true},
		// 无效
		{"", "", "", "", false},
		{"600001.XX", "", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			var s Symbol
			err := s.Parse(tt.symbol)
			ok := err == nil
			if ok != tt.wantOK {
				t.Errorf("Parse() ok = %v, want %v, err = %v", ok, tt.wantOK, err)
			}
			if tt.wantOK {
				if s.Code != tt.wantCode {
					t.Errorf("Parse() code = %v, want %v", s.Code, tt.wantCode)
				}
				if s.Exchange != tt.wantExchange {
					t.Errorf("Parse() exchange = %v, want %v", s.Exchange, tt.wantExchange)
				}
				if s.Market != tt.wantMarket {
					t.Errorf("Parse() market = %v, want %v", s.Market, tt.wantMarket)
				}
			}
		})
	}
}

func TestSmartParse(t *testing.T) {
	tests := []struct {
		input        string
		wantCode     string
		wantMarket   Market
		wantExchange Exchange
	}{
		// A股：6位纯数字
		{"000001", "000001", MarketCN, ExchangeSZ},
		{"600519", "600519", MarketCN, ExchangeSH},
		{"300001", "300001", MarketCN, ExchangeSZ},
		{"688001", "688001", MarketCN, ExchangeSH},
		{"430001", "430001", MarketCN, ExchangeBJ},
		{"830001", "830001", MarketCN, ExchangeBJ},
		{"920001", "920001", MarketCN, ExchangeBJ},

		// A股ETF/基金：3位前缀判断
		{"510300", "510300", MarketCN, ExchangeSH}, // 上交所ETF
		{"511010", "511010", MarketCN, ExchangeSH}, // 上交所国债ETF
		{"520030", "520030", MarketCN, ExchangeSH}, // 上交所ETF
		{"560050", "560050", MarketCN, ExchangeSH}, // 上交所ETF
		{"580060", "580060", MarketCN, ExchangeSH}, // 上交所ETF
		{"159915", "159915", MarketCN, ExchangeSZ}, // 深交所ETF
		{"150001", "150001", MarketCN, ExchangeSZ}, // 深交所分级基金
		{"160105", "160105", MarketCN, ExchangeSZ}, // 深交所LOF基金
		{"180001", "180001", MarketCN, ExchangeSZ}, // 深交所ETF

		// 港股：5位
		{"00700", "00700", MarketHK, ExchangeHKEX},

		// 加密货币：含分隔符
		{"BTC-USD", "BTCUSD", MarketCrypto, ExchangeBinance},
		{"ETH/USDT", "ETHUSDT", MarketCrypto, ExchangeBinance},
		{"BTC_USDT", "BTCUSDT", MarketCrypto, ExchangeBinance},
		{"ETH-USDC", "ETHUSDC", MarketCrypto, ExchangeBinance},
		{"SOL_USDT", "SOLUSDT", MarketCrypto, ExchangeBinance},
		// 加密货币：CCXT 格式（含结算货币）
		{"BTC/USDT:USDT", "BTCUSDT", MarketCrypto, ExchangeBinance},
		{"ETH/USDT:USDT", "ETHUSDT", MarketCrypto, ExchangeBinance},

		// 加密货币：quote suffix 匹配
		{"BTCUSDT", "BTCUSDT", MarketCrypto, ExchangeBinance},
		{"ETHUSDC", "ETHUSDC", MarketCrypto, ExchangeBinance},
		{"SOLUSDT", "SOLUSDT", MarketCrypto, ExchangeBinance},
		{"DOGEUSD", "DOGEUSD", MarketCrypto, ExchangeBinance},
		{"XRPBTC", "XRPBTC", MarketCrypto, ExchangeBinance},
		{"LINKETH", "LINKETH", MarketCrypto, ExchangeBinance},
		{"AVAXBNB", "AVAXBNB", MarketCrypto, ExchangeBinance},
		{"PEPEUSDT", "PEPEUSDT", MarketCrypto, ExchangeBinance},

		{"XRPBTC", "XRPBTC", MarketCrypto, ExchangeBinance}, // 6字母但非forex → crypto

		// 外汇：6位纯字母（ISO 4217 货币对）
		{"EURUSD", "EURUSD", MarketForex, ExchangeForexSpot},
		{"GBPJPY", "GBPJPY", MarketForex, ExchangeForexSpot},
		{"USDCHF", "USDCHF", MarketForex, ExchangeForexSpot},

		// 美股：1-5位纯字母（最后分支）
		{"AAPL", "AAPL", MarketUS, ExchangeNASDAQ},
		{"GOOG", "GOOG", MarketUS, ExchangeNASDAQ},
		{"TSLA", "TSLA", MarketUS, ExchangeNASDAQ},
		{"SOL", "SOL", MarketUS, ExchangeNASDAQ}, // 3字母默认美股，需 .CRYPTO 强制指定
		{"BRK", "BRK", MarketUS, ExchangeNASDAQ},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var s Symbol
			if err := s.SmartParse(tt.input); err != nil {
				t.Fatalf("SmartParse(%q) error: %v", tt.input, err)
			}
			if s.Code != tt.wantCode {
				t.Errorf("SmartParse(%q) code = %q, want %q", tt.input, s.Code, tt.wantCode)
			}
			if s.Market != tt.wantMarket {
				t.Errorf("SmartParse(%q) market = %q, want %q", tt.input, s.Market, tt.wantMarket)
			}
			if s.Exchange != tt.wantExchange {
				t.Errorf("SmartParse(%q) exchange = %q, want %q", tt.input, s.Exchange, tt.wantExchange)
			}
		})
	}
}

func TestSmartParseAmbiguity(t *testing.T) {
	// "SOL" 3字母纯字母 → 默认美股
	var s1 Symbol
	if err := s1.SmartParse("SOL"); err != nil {
		t.Fatalf("SmartParse(SOL) error: %v", err)
	}
	if s1.Market != MarketUS {
		t.Errorf("SOL should default to US stock, got market %s", s1.Market)
	}

	// "SOLUSDT" → crypto (quote suffix 匹配)
	var s2 Symbol
	if err := s2.SmartParse("SOLUSDT"); err != nil {
		t.Fatalf("SmartParse(SOLUSDT) error: %v", err)
	}
	if s2.Market != MarketCrypto {
		t.Errorf("SOLUSDT should be crypto, got market %s", s2.Market)
	}

	// 用户强制指定 crypto：SOL.CRYPTO.BINANCE
	var s3 Symbol
	if err := s3.Parse("SOL.CRYPTO.BINANCE"); err != nil {
		t.Fatalf("Parse(SOL.CRYPTO.BINANCE) error: %v", err)
	}
	if s3.Market != MarketCrypto {
		t.Errorf("SOL.CRYPTO.BINANCE should be crypto, got market %s", s3.Market)
	}
}

func TestMatchCryptoQuoteSuffix(t *testing.T) {
	tests := []struct {
		code      string
		wantBase  string
		wantQuote string
		wantOK    bool
	}{
		{"BTCUSDT", "BTC", "USDT", true},
		{"ETHUSDC", "ETH", "USDC", true},
		{"SOLUSDT", "SOL", "USDT", true},
		{"XRPBTC", "XRP", "BTC", true},
		{"DOGEUSD", "DOGE", "USD", true},
		{"LINKETH", "LINK", "ETH", true},
		{"PEPEUSDT", "PEPE", "USDT", true},
		{"AVAXBNB", "AVAX", "BNB", true},
		{"SOL", "", "", false},         // 太短，无 quote suffix
		{"AAPL", "", "", false},        // 美股代码
		{"EURUSD", "EUR", "USD", true}, // MatchCryptoQuoteSuffix 本身会匹配，但 SmartParse 优先识别为 forex
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			base, quote, ok := MatchCryptoQuoteSuffix(tt.code)
			if ok != tt.wantOK {
				t.Errorf("MatchCryptoQuoteSuffix(%q) ok = %v, want %v", tt.code, ok, tt.wantOK)
			}
			if tt.wantOK {
				if base != tt.wantBase {
					t.Errorf("MatchCryptoQuoteSuffix(%q) base = %q, want %q", tt.code, base, tt.wantBase)
				}
				if quote != tt.wantQuote {
					t.Errorf("MatchCryptoQuoteSuffix(%q) quote = %q, want %q", tt.code, quote, tt.wantQuote)
				}
			}
		})
	}
}

func TestInferCNExchange(t *testing.T) {
	tests := []struct {
		code string
		want Exchange
	}{
		// 上交所主板
		{"600000", ExchangeSH},
		{"601398", ExchangeSH},
		// 科创板
		{"688001", ExchangeSH},
		// 上交所指数
		{"900001", ExchangeSH},
		{"890001", ExchangeSH},
		// 上交所ETF
		{"510300", ExchangeSH},
		{"511010", ExchangeSH},
		{"512100", ExchangeSH},
		{"513500", ExchangeSH},
		{"515050", ExchangeSH},
		{"516160", ExchangeSH},
		{"518880", ExchangeSH},
		{"520030", ExchangeSH},
		{"560050", ExchangeSH},
		{"580060", ExchangeSH},

		// 深交所主板
		{"000001", ExchangeSZ},
		{"000002", ExchangeSZ},
		// 创业板
		{"300001", ExchangeSZ},
		{"300750", ExchangeSZ},
		// 科创板
		{"200001", ExchangeSZ},
		// 深交所ETF/基金
		{"159915", ExchangeSZ},
		{"159001", ExchangeSZ},
		{"150001", ExchangeSZ},
		{"160105", ExchangeSZ},
		{"161725", ExchangeSZ},
		{"180001", ExchangeSZ},

		// 北交所
		{"430001", ExchangeBJ},
		{"830001", ExchangeBJ},
		{"870001", ExchangeBJ},
		{"920001", ExchangeBJ},

		// 默认
		{"999999", ExchangeSZ},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got := inferCNExchange(tt.code)
			if got != tt.want {
				t.Errorf("inferCNExchange(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestFormatSymbol(t *testing.T) {
	tests := []struct {
		code     string
		exchange Exchange
		want     string
	}{
		{"000001", ExchangeSZ, "000001.SZ"},
		{"600001", ExchangeSH, "600001.SH"},
		{"430001", ExchangeBJ, "430001.BJ"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := FormatSymbol(tt.code, tt.exchange); got != tt.want {
				t.Errorf("FormatSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAndFormatRoundTrip(t *testing.T) {
	symbols := []string{"000001.SZ", "600001.SH", "430001.BJ"}
	for _, s := range symbols {
		code, exchange, ok := ParseSymbol(s)
		if !ok {
			t.Errorf("ParseSymbol(%s) failed", s)
			continue
		}
		got := FormatSymbol(code, exchange)
		if got != s {
			t.Errorf("Round trip: %s -> (%s, %s) -> %s", s, code, exchange, got)
		}
	}
}

func TestNormalizeCryptoSymbol(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Slash format
		{"BTC/USDT", "BTCUSDT"},
		{"ETH/USDC", "ETHUSDC"},
		// Dash format
		{"BTC-USD", "BTCUSD"},
		{"ETH-USDC", "ETHUSDC"},
		// Underscore format
		{"BTC_USDT", "BTCUSDT"},
		{"SOL_USDT", "SOLUSDT"},
		// CCXT format with settlement currency
		{"BTC/USDT:USDT", "BTCUSDT"},
		{"ETH/USDT:USDT", "ETHUSDT"},
		// Already concatenated
		{"BTCUSDT", "BTCUSDT"},
		// Lowercase
		{"btcusdt", "BTCUSDT"},
		// Whitespace
		{"  BTC/USDT  ", "BTCUSDT"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := NormalizeCryptoSymbol(tt.input); got != tt.want {
				t.Errorf("NormalizeCryptoSymbol(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSymbolShort(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"000001.CN.SZ", "000001.SZ"},
		{"600519.CN.SH", "600519.SH"},
		{"AAPL.US.NASDAQ", "AAPL.US.NASDAQ"},
		{"BTCUSDT.CRYPTO.BINANCE", "BTCUSDT.CRYPTO.BINANCE"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var s Symbol
			if err := s.Parse(tt.input); err != nil {
				t.Fatalf("Parse(%q) error: %v", tt.input, err)
			}
			if got := s.Short(); got != tt.want {
				t.Errorf("Short() = %q, want %q", got, tt.want)
			}
		})
	}
}
