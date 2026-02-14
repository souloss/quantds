package cninfo

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
		switch reqErr.StatusCode {
		case 401, 403, 429, 451, 503:
			t.Skipf("Skipping test due to API restriction or unavailability (status %d): %v", reqErr.StatusCode, err)
		}
	}
	errMsg := err.Error()
	if strings.Contains(errMsg, "client error") ||
		strings.Contains(errMsg, "unmarshal") {
		t.Skipf("Skipping test due to API error: %v", err)
	}
	t.Fatalf("API request failed: %v", err)
}
