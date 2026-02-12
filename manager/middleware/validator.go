package middleware

import (
	"context"
	"errors"

	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

func Validator[Req, Resp any](validate func(Resp) error) Middleware[Req, Resp] {
	return func(next manager.Provider[Req, Resp]) manager.Provider[Req, Resp] {
		return ProviderFunc[Req, Resp](next.Name(), func(ctx context.Context, client request.Client, req Req) (Resp, *manager.RequestTrace, error) {
			resp, trace, err := next.Fetch(ctx, client, req)
			if err != nil {
				return resp, trace, err
			}
			if validate != nil {
				if verr := validate(resp); verr != nil {
					return resp, trace, errors.Join(manager.ErrAllProviderFailed, verr)
				}
			}
			return resp, trace, nil
		})
	}
}
