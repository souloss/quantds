// Package formatter provides SymbolFormatter implementations for each data source.
//
// Each formatter implements the domain.SymbolFormatter interface with pure
// functions that convert between the standard Symbol format and the
// data-source-specific API format. The rules are finite and stable,
// requiring zero data maintenance.
package formatter

import (
	"github.com/souloss/quantds/domain"
)

// Registry holds all available formatters keyed by provider name.
type Registry struct {
	formatters map[string]domain.SymbolFormatter
}

// NewRegistry creates a new formatter registry with all built-in formatters.
func NewRegistry() *Registry {
	r := &Registry{
		formatters: make(map[string]domain.SymbolFormatter),
	}
	// CN data sources
	r.Register("tushare", &TushareFormatter{})
	r.Register("eastmoney", &EastMoneyFormatter{})
	r.Register("sina", &SinaFormatter{})
	r.Register("xueqiu", &XueqiuFormatter{})
	r.Register("cninfo", &CNInfoFormatter{})
	r.Register("akshare", &TushareFormatter{}) // same format as tushare
	// US data sources
	r.Register("yahoo", &YahooFormatter{})
	r.Register("alphavantage", &AlphaVantageFormatter{})
	r.Register("finnhub", &FinnhubFormatter{})
	r.Register("polygon", &PolygonFormatter{})
	r.Register("twelvedata", &TwelveDataFormatter{})
	r.Register("eodhd", &EODHDFormatter{})
	// HK data sources
	r.Register("eastmoneyhk", &EastMoneyHKFormatter{})
	// Crypto data sources
	r.Register("binance", &BinanceFormatter{})
	r.Register("okx", &OKXFormatter{})
	r.Register("coingecko", &CoinGeckoFormatter{})
	return r
}

// Register adds a formatter to the registry.
func (r *Registry) Register(name string, f domain.SymbolFormatter) {
	r.formatters[name] = f
}

// Get returns the formatter for the given provider name.
func (r *Registry) Get(name string) domain.SymbolFormatter {
	return r.formatters[name]
}

// Format is a convenience method that looks up the formatter and formats the symbol.
func (r *Registry) Format(provider string, sym domain.Symbol) string {
	f := r.formatters[provider]
	if f == nil {
		return sym.Code // fallback: return code as-is
	}
	return f.Format(sym)
}

// Parse is a convenience method that looks up the formatter and parses the raw string.
func (r *Registry) Parse(provider, raw string) (domain.Symbol, bool) {
	f := r.formatters[provider]
	if f == nil {
		return domain.Symbol{}, false
	}
	return f.Parse(raw)
}
