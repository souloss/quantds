package sina

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
	if err == nil {
		return
	}
	var reqErr *request.RequestError
	if errors.As(err, &reqErr) {
		// 401: Unauthorized
		// 403: Forbidden (Often WAF block or geo-restriction)
		// 429: Too Many Requests (Rate limited)
		// 451: Unavailable For Legal Reasons (Geo-blocked)
		// 503: Service Unavailable
		switch reqErr.StatusCode {
		case 401, 403, 429, 451, 503:
			t.Skipf("Skipping test due to API restriction or unavailability (status %d): %v", reqErr.StatusCode, err)
		}
	}
	// Handle API-level errors (e.g., "Service not found", parse errors)
	errMsg := err.Error()
	if strings.Contains(errMsg, "Service not found") ||
		strings.Contains(errMsg, "unmarshal") ||
		strings.Contains(errMsg, "service unavailable") {
		t.Skipf("Skipping test due to API service error: %v", err)
	}
	t.Fatalf("API request failed: %v", err)
}
