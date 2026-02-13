package binance

import (
	"errors"
	"testing"

	"github.com/souloss/quantds/request"
)

func checkAPIError(t *testing.T, err error) {
	if err == nil {
		return
	}
	var reqErr *request.RequestError
	if errors.As(err, &reqErr) {
		// 451: Unavailable For Legal Reasons (Geo-blocked)
		// 403: Forbidden (Often WAF block)
		if reqErr.StatusCode == 451 || reqErr.StatusCode == 403 {
			t.Skipf("Skipping test due to API restriction (status %d): %v", reqErr.StatusCode, err)
		}
	}
	t.Fatalf("API request failed: %v", err)
}
