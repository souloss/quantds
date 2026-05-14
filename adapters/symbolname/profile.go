package symbolname

import (
	"context"
	"fmt"
	"strings"

	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/profile"
	"github.com/souloss/quantds/manager"
)

// Name is the adapter name
const Name = "symbolname"

// ProfileAdapter resolves symbol names for Crypto, Forex, and Futures markets
// using local seed mapping tables. No external API calls are needed.
type ProfileAdapter struct {
	markets  []domain.Market
	resolver func(string) string
}

// NewCryptoProfileAdapter creates a profile adapter for Crypto market
func NewCryptoProfileAdapter() *ProfileAdapter {
	return &ProfileAdapter{
		markets:  []domain.Market{domain.MarketCrypto},
		resolver: resolveCryptoName,
	}
}

// NewForexProfileAdapter creates a profile adapter for Forex market
func NewForexProfileAdapter() *ProfileAdapter {
	return &ProfileAdapter{
		markets:  []domain.Market{domain.MarketForex},
		resolver: resolveForexName,
	}
}

// NewFuturesProfileAdapter creates a profile adapter for Futures market
func NewFuturesProfileAdapter() *ProfileAdapter {
	return &ProfileAdapter{
		markets:  []domain.Market{domain.MarketFutures},
		resolver: resolveFuturesName,
	}
}

// Name returns the adapter name
func (a *ProfileAdapter) Name() string { return Name }

// SupportedMarkets returns supported markets
func (a *ProfileAdapter) SupportedMarkets() []domain.Market { return a.markets }

// CanHandle checks if the adapter can handle the symbol
func (a *ProfileAdapter) CanHandle(symbol string) bool {
	var sym domain.Symbol
	if err := sym.Parse(symbol); err != nil {
		return false
	}
	for _, m := range a.markets {
		if sym.Market == m {
			return true
		}
	}
	return false
}

// Fetch resolves the symbol name and returns a profile response
func (a *ProfileAdapter) Fetch(ctx context.Context, req profile.Request) (profile.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	var sym domain.Symbol
	if err := sym.Parse(req.Symbol); err != nil {
		return profile.Response{}, trace, fmt.Errorf("invalid symbol %s: %w", req.Symbol, err)
	}

	name := a.resolver(sym.Code)

	// For futures, use exchange-aware resolution when available
	if sym.Market == domain.MarketFutures && sym.Exchange != "" {
		name = resolveFuturesNameWithExchange(sym.Code, string(sym.Exchange))
	}

	trace.Finish()
	return profile.Response{
		Data: profile.Profile{
			Symbol:   req.Symbol,
			Name:     name,
			FullName: buildFullName(name, sym),
		},
		Source:      Name,
		DataVersion: 1,
	}, trace, nil
}

// buildFullName constructs a more descriptive full name from the resolved name and symbol metadata
func buildFullName(name string, sym domain.Symbol) string {
	parts := []string{name}

	switch sym.Market {
	case domain.MarketCrypto:
		parts = append(parts, "(Crypto)")
	case domain.MarketForex:
		parts = append(parts, "(Forex)")
	case domain.MarketFutures:
		parts = append(parts, "(Futures)")
	}

	if sym.Exchange != "" && sym.Exchange != domain.ExchangeBinance &&
		sym.Exchange != domain.ExchangeForexSpot && sym.Exchange != domain.ExchangeSHFE {
		parts = append(parts, strings.ToLower(string(sym.Exchange)))
	}

	return strings.Join(parts, " ")
}

var _ manager.Provider[profile.Request, profile.Response] = (*ProfileAdapter)(nil)
