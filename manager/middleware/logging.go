package middleware

import (
	"context"
	"log"

	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

func Logging[Req, Resp any](logger *log.Logger) Middleware[Req, Resp] {
	return func(next manager.Provider[Req, Resp]) manager.Provider[Req, Resp] {
		return ProviderFunc[Req, Resp](next.Name(), func(ctx context.Context, client request.Client, req Req) (Resp, *manager.RequestTrace, error) {
			resp, trace, err := next.Fetch(ctx, client, req)
			if logger != nil {
				if err != nil {
					logger.Printf("[Fetch] provider=%s error=%v", next.Name(), err)
				} else {
					logger.Printf("[Fetch] provider=%s requests=%d duration=%v", next.Name(), trace.TotalRequests(), trace.TotalDuration())
				}
			}
			return resp, trace, err
		})
	}
}
