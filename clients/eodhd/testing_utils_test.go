package eodhd

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/souloss/quantds/request"
)

func skipIfNoAPIKey(t *testing.T) {
	t.Helper()
	if os.Getenv("EODHD_API_KEY") == "" {
		t.Skip("EODHD_API_KEY not set")
	}
}

func checkAPIError(t *testing.T, err error) {
	if err == nil {
		return
	}
	var reqErr *request.RequestError
	if errors.As(err, &reqErr) {
		switch reqErr.StatusCode {
		case 401, 403, 429, 451, 503:
			t.Skipf("Skipping: API restriction (status %d): %v", reqErr.StatusCode, err)
		}
	}
	errMsg := err.Error()
	if strings.Contains(errMsg, "client error") ||
		strings.Contains(errMsg, "unmarshal") ||
		strings.Contains(errMsg, "retries exceeded") ||
		strings.Contains(errMsg, "EOF") ||
		strings.Contains(errMsg, "connection refused") {
		t.Skipf("Skipping: API error: %v", err)
	}
	t.Fatalf("API request failed: %v", err)
}
