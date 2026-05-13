package formatter

import (
	"testing"

	"github.com/souloss/quantds/domain"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	if r.Get("tushare") == nil {
		t.Error("Expected tushare formatter")
	}
	if r.Get("binance") == nil {
		t.Error("Expected binance formatter")
	}
	if r.Get("nonexistent") != nil {
		t.Error("Expected nil for unknown provider")
	}
}

func TestTushareFormatter(t *testing.T) {
	f := &TushareFormatter{}
	tests := []struct {
		sym  domain.Symbol
		want string
	}{
		{domain.Symbol{Code: "000001", Market: domain.MarketCN, Exchange: domain.ExchangeSZ}, "000001.SZ"},
		{domain.Symbol{Code: "600519", Market: domain.MarketCN, Exchange: domain.ExchangeSH}, "600519.SH"},
		{domain.Symbol{Code: "430001", Market: domain.MarketCN, Exchange: domain.ExchangeBJ}, "430001.BJ"},
	}
	for _, tt := range tests {
		got := f.Format(tt.sym)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.sym, got, tt.want)
		}
	}

	// Parse round-trip
	for _, tt := range tests {
		parsed, ok := f.Parse(tt.want)
		if !ok {
			t.Errorf("Parse(%q) failed", tt.want)
			continue
		}
		if parsed.Code != tt.sym.Code {
			t.Errorf("Parse(%q) code = %q, want %q", tt.want, parsed.Code, tt.sym.Code)
		}
		if parsed.Exchange != tt.sym.Exchange {
			t.Errorf("Parse(%q) exchange = %q, want %q", tt.want, parsed.Exchange, tt.sym.Exchange)
		}
	}
}

func TestEastMoneyFormatter(t *testing.T) {
	f := &EastMoneyFormatter{}
	tests := []struct {
		sym  domain.Symbol
		want string
	}{
		{domain.Symbol{Code: "000001", Market: domain.MarketCN, Exchange: domain.ExchangeSZ}, "0.000001"},
		{domain.Symbol{Code: "600519", Market: domain.MarketCN, Exchange: domain.ExchangeSH}, "1.600519"},
	}
	for _, tt := range tests {
		got := f.Format(tt.sym)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.sym, got, tt.want)
		}
	}
	// Parse
	parsed, ok := f.Parse("1.600519")
	if !ok || parsed.Code != "600519" || parsed.Exchange != domain.ExchangeSH {
		t.Errorf("Parse(1.600519) = %v, ok=%v", parsed, ok)
	}
}

func TestSinaFormatter(t *testing.T) {
	f := &SinaFormatter{}
	tests := []struct {
		sym  domain.Symbol
		want string
	}{
		{domain.Symbol{Code: "600519", Market: domain.MarketCN, Exchange: domain.ExchangeSH}, "sh600519"},
		{domain.Symbol{Code: "000001", Market: domain.MarketCN, Exchange: domain.ExchangeSZ}, "sz000001"},
	}
	for _, tt := range tests {
		got := f.Format(tt.sym)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.sym, got, tt.want)
		}
	}
}

func TestXueqiuFormatter(t *testing.T) {
	f := &XueqiuFormatter{}
	tests := []struct {
		sym  domain.Symbol
		want string
	}{
		{domain.Symbol{Code: "600519", Market: domain.MarketCN, Exchange: domain.ExchangeSH}, "SH600519"},
		{domain.Symbol{Code: "000001", Market: domain.MarketCN, Exchange: domain.ExchangeSZ}, "SZ000001"},
		{domain.Symbol{Code: "AAPL", Market: domain.MarketUS, Exchange: domain.ExchangeNASDAQ}, "AAPL"},
	}
	for _, tt := range tests {
		got := f.Format(tt.sym)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.sym, got, tt.want)
		}
	}
}

func TestAlphaVantageFormatter(t *testing.T) {
	f := &AlphaVantageFormatter{}
	tests := []struct {
		sym  domain.Symbol
		want string
	}{
		{domain.Symbol{Code: "AAPL", Market: domain.MarketUS, Exchange: domain.ExchangeNASDAQ}, "AAPL"},
		{domain.Symbol{Code: "EURUSD", Market: domain.MarketForex, Exchange: domain.ExchangeForexSpot}, "EUR_USD"},
	}
	for _, tt := range tests {
		got := f.Format(tt.sym)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.sym, got, tt.want)
		}
	}
	// Parse forex
	parsed, ok := f.Parse("EUR_USD")
	if !ok || parsed.Code != "EURUSD" || parsed.Market != domain.MarketForex {
		t.Errorf("Parse(EUR_USD) = %v, ok=%v", parsed, ok)
	}
}

func TestBinanceFormatter(t *testing.T) {
	f := &BinanceFormatter{}
	tests := []struct {
		sym  domain.Symbol
		want string
	}{
		{domain.Symbol{Code: "BTCUSDT", Market: domain.MarketCrypto, Exchange: domain.ExchangeBinance}, "BTCUSDT"},
		{domain.Symbol{Code: "ETHUSDC", Market: domain.MarketCrypto, Exchange: domain.ExchangeBinance}, "ETHUSDC"},
	}
	for _, tt := range tests {
		got := f.Format(tt.sym)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.sym, got, tt.want)
		}
	}
	// Parse
	parsed, ok := f.Parse("BTCUSDT")
	if !ok || parsed.Code != "BTCUSDT" || parsed.Market != domain.MarketCrypto {
		t.Errorf("Parse(BTCUSDT) = %v, ok=%v", parsed, ok)
	}
}

func TestOKXFormatter(t *testing.T) {
	f := &OKXFormatter{}
	tests := []struct {
		sym  domain.Symbol
		want string
	}{
		{domain.Symbol{Code: "BTCUSDT", Market: domain.MarketCrypto, Exchange: domain.ExchangeBinance}, "BTC-USDT"},
		{domain.Symbol{Code: "ETHUSDC", Market: domain.MarketCrypto, Exchange: domain.ExchangeBinance}, "ETH-USDC"},
	}
	for _, tt := range tests {
		got := f.Format(tt.sym)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.sym, got, tt.want)
		}
	}
	// Parse
	parsed, ok := f.Parse("BTC-USDT")
	if !ok || parsed.Code != "BTCUSDT" || parsed.Market != domain.MarketCrypto {
		t.Errorf("Parse(BTC-USDT) = %v, ok=%v", parsed, ok)
	}
}

func TestCoinGeckoFormatter(t *testing.T) {
	f := &CoinGeckoFormatter{}
	tests := []struct {
		sym  domain.Symbol
		want string
	}{
		{domain.Symbol{Code: "BTCUSDT", Market: domain.MarketCrypto, Exchange: domain.ExchangeBinance}, "bitcoin"},
		{domain.Symbol{Code: "ETHUSDT", Market: domain.MarketCrypto, Exchange: domain.ExchangeBinance}, "ethereum"},
		{domain.Symbol{Code: "SOLUSDT", Market: domain.MarketCrypto, Exchange: domain.ExchangeBinance}, "solana"},
	}
	for _, tt := range tests {
		got := f.Format(tt.sym)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.sym, got, tt.want)
		}
	}
	// Parse
	parsed, ok := f.Parse("bitcoin")
	if !ok || parsed.Code != "BTC" {
		t.Errorf("Parse(bitcoin) = %v, ok=%v", parsed, ok)
	}
}

func TestEODHDFormatter(t *testing.T) {
	f := &EODHDFormatter{}
	tests := []struct {
		sym  domain.Symbol
		want string
	}{
		{domain.Symbol{Code: "AAPL", Market: domain.MarketUS, Exchange: domain.ExchangeNASDAQ}, "AAPL.NASDAQ"},
		{domain.Symbol{Code: "600519", Market: domain.MarketCN, Exchange: domain.ExchangeSH}, "600519.SH"},
		{domain.Symbol{Code: "00700", Market: domain.MarketHK, Exchange: domain.ExchangeHKEX}, "00700.HK"},
	}
	for _, tt := range tests {
		got := f.Format(tt.sym)
		if got != tt.want {
			t.Errorf("Format(%v) = %q, want %q", tt.sym, got, tt.want)
		}
	}
}

func TestEastMoneyHKFormatter(t *testing.T) {
	f := &EastMoneyHKFormatter{}
	got := f.Format(domain.Symbol{Code: "00700", Market: domain.MarketHK, Exchange: domain.ExchangeHKEX})
	if got != "116.00700" {
		t.Errorf("Format() = %q, want %q", got, "116.00700")
	}
	parsed, ok := f.Parse("116.00700")
	if !ok || parsed.Code != "00700" || parsed.Market != domain.MarketHK {
		t.Errorf("Parse(116.00700) = %v, ok=%v", parsed, ok)
	}
}

func TestRegistryFormat(t *testing.T) {
	r := NewRegistry()
	sym := domain.Symbol{Code: "000001", Market: domain.MarketCN, Exchange: domain.ExchangeSZ}

	tests := []struct {
		provider string
		want     string
	}{
		{"tushare", "000001.SZ"},
		{"eastmoney", "0.000001"},
		{"sina", "sz000001"},
		{"xueqiu", "SZ000001"},
		{"cninfo", "000001"},
		{"akshare", "000001.SZ"},
	}

	for _, tt := range tests {
		got := r.Format(tt.provider, sym)
		if got != tt.want {
			t.Errorf("Registry.Format(%q, %v) = %q, want %q", tt.provider, sym, got, tt.want)
		}
	}
}
