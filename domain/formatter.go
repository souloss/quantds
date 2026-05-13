package domain

// SymbolFormatter defines the interface for converting between standard Symbol
// and data-source-specific API formats.
//
// Core insight: symbol conversion is an algorithm problem (finite rules per
// data source per market), not a data problem (infinite symbol mappings).
// Each data source has deterministic format rules that are finite and stable.
type SymbolFormatter interface {
	// Format converts a standard Symbol to the data-source's API format.
	// For example, Symbol{Code:"000001", Market:MarketCN, Exchange:ExchangeSZ}
	// → "000001.SZ" for tushare, "0.000001" for eastmoney.
	Format(sym Symbol) string

	// Parse converts a data-source's raw format back to a standard Symbol.
	// Returns (Symbol, true) on success, (zero, false) on failure.
	Parse(raw string) (Symbol, bool)
}
