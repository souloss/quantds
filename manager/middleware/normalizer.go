package middleware

import (
	"context"

	"github.com/souloss/quantds/manager"
)

// Normalizer creates a middleware that normalizes the response data after fetching.
// The normalize function is applied to successful responses only.
// This is useful for filling in computed fields (e.g., ChangeRate from Close/PreClose)
// that some data sources may not provide directly.
func Normalizer[Req, Resp any](normalize func(Resp) Resp) Middleware[Req, Resp] {
	return func(next manager.Provider[Req, Resp]) manager.Provider[Req, Resp] {
		return &providerFunc[Req, Resp]{
			name:             next.Name(),
			supportedMarkets: next.SupportedMarkets(),
			canHandle:        next.CanHandle,
			fetch: func(ctx context.Context, req Req) (Resp, *manager.RequestTrace, error) {
				resp, trace, err := next.Fetch(ctx, req)
				if err != nil {
					return resp, trace, err
				}
				return normalize(resp), trace, nil
			},
		}
	}
}
