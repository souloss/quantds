package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/souloss/quantds/manager"
)

func Logging[Req, Resp any](logger *slog.Logger) Middleware[Req, Resp] {
	return func(next manager.Provider[Req, Resp]) manager.Provider[Req, Resp] {
		return ProviderFunc[Req, Resp](next.Name(), func(ctx context.Context, req Req) (Resp, *manager.RequestTrace, error) {
			start := time.Now()
			resp, trace, err := next.Fetch(ctx, req)
			duration := time.Since(start)

			attrs := []any{
				"provider", next.Name(),
				"duration", duration,
			}
			if err != nil {
				attrs = append(attrs, "error", err)
				logger.Error("provider fetch failed", attrs...)
			} else {
				logger.Info("provider fetch succeeded", attrs...)
			}

			return resp, trace, err
		})
	}
}
