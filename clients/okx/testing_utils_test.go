package okx

import (
	"errors"
	"testing"

	"github.com/souloss/quantds/request"
)

// checkAPIError checks if the error is due to API geo-restriction or unavailability
// and skips the test if so.
func checkAPIError(t *testing.T, err error) {
	if err == nil {
		return
	}
	var reqErr *request.RequestError
	if errors.As(err, &reqErr) {
		// 451: Unavailable For Legal Reasons (Geo-blocked)
		// 403: Forbidden (Often WAF block)
		// 503: Service Unavailable
		if reqErr.StatusCode == 451 || reqErr.StatusCode == 403 || reqErr.StatusCode == 503 {
			t.Skipf("Skipping test due to API restriction or unavailability (status %d): %v", reqErr.StatusCode, err)
		}
	}
	t.Fatalf("API request failed: %v", err)
}
