package coingecko

import (
	"errors"
	"strings"
	"testing"

	"github.com/souloss/quantds/request"
)

// checkAPIError checks if the error is due to API geo-restriction or unavailability
// and skips the test if so. This follows the README requirement for handling
// geo-restrictions gracefully in tests.
func checkAPIError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	var reqErr *request.RequestError
	if errors.As(err, &reqErr) {
		// 451: Unavailable For Legal Reasons (Geo-blocked)
		// 403: Forbidden (Often WAF block or geo-restriction)
		// 429: Too Many Requests (Rate limited)
		// 503: Service Unavailable
		if reqErr.StatusCode == 451 || reqErr.StatusCode == 403 || reqErr.StatusCode == 429 || reqErr.StatusCode == 503 {
			t.Skipf("Skipping test due to API restriction or unavailability (status %d): %v", reqErr.StatusCode, err)
		}
	}
	// Fallback: also match common error messages for rate limiting and auth
	msg := err.Error()
	if strings.Contains(msg, "rate limit") || strings.Contains(msg, "429") ||
		strings.Contains(msg, "too many requests") {
		t.Skipf("Skipping test due to rate limiting: %v", err)
	}
	t.Fatalf("API request failed: %v", err)
}
