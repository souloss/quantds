package middleware

import (
	"context"

	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/manager"
)

type Middleware[Req, Resp any] func(next manager.Provider[Req, Resp]) manager.Provider[Req, Resp]

func Chain[Req, Resp any](middlewares ...Middleware[Req, Resp]) Middleware[Req, Resp] {
	return func(next manager.Provider[Req, Resp]) manager.Provider[Req, Resp] {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

type providerFunc[Req, Resp any] struct {
	name             string
	fetch            func(ctx context.Context, req Req) (Resp, *manager.RequestTrace, error)
	supportedMarkets []domain.Market
	canHandle        func(symbol string) bool
	cbState          func() CircuitBreakerState // optional: exposes circuit breaker state
}

func (p *providerFunc[Req, Resp]) Name() string {
	return p.name
}

func (p *providerFunc[Req, Resp]) SupportedMarkets() []domain.Market {
	if p.supportedMarkets == nil {
		return []domain.Market{domain.MarketCN}
	}
	return p.supportedMarkets
}

func (p *providerFunc[Req, Resp]) CanHandle(symbol string) bool {
	if p.canHandle == nil {
		return true
	}
	return p.canHandle(symbol)
}

// CircuitBreakerState implements CircuitBreakerStateProvider.
// Returns CircuitClosed if no circuit breaker is attached.
func (p *providerFunc[Req, Resp]) CircuitBreakerState() CircuitBreakerState {
	if p.cbState != nil {
		return p.cbState()
	}
	return CircuitClosed
}

func (p *providerFunc[Req, Resp]) Fetch(ctx context.Context, req Req) (Resp, *manager.RequestTrace, error) {
	return p.fetch(ctx, req)
}

func ProviderFunc[Req, Resp any](name string, fetch func(ctx context.Context, req Req) (Resp, *manager.RequestTrace, error)) manager.Provider[Req, Resp] {
	return &providerFunc[Req, Resp]{name: name, fetch: fetch}
}
