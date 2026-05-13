// Package testutil provides common testing utilities for quantds client and adapter tests.
//
// It establishes a three-layer testing specification:
//
//	Layer 1 (Normal Data): Verify API reachability, key fields are non-zero, data logical consistency.
//	Layer 2 (API Limits): Test pagination limits, rate limiting, time range constraints.
//	Layer 3 (Failure Detection): Test invalid symbol handling, empty parameter handling, error classification.
package testutil

import (
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
)

// APITestResult records the result of an API test call.
type APITestResult struct {
	Available  bool   // Whether the API endpoint is reachable
	Skipped    bool   // Whether the test was skipped
	SkipReason string // Reason for skipping
	DataValid  bool   // Whether the returned data passes basic validation
	Error      error  // Any error encountered
}

// ClassifyAPIError determines whether an error is transient (temporary, should skip test)
// or permanent (indicates a real problem, should fail test).
//
// Transient errors include: EOF, timeout, rate limit (429), service unavailable (503),
// connection refused, and DNS resolution failures.
//
// Permanent errors include: authentication failures, invalid parameters, and other
// non-recoverable issues.
func ClassifyAPIError(err error) (transient bool, reason string) {
	if err == nil {
		return false, ""
	}

	// Network-level transient errors
	if err == io.EOF {
		return true, "EOF from server"
	}
	if _, ok := err.(net.Error); ok {
		return true, "network error"
	}

	msg := err.Error()

	// Transient patterns
	transientPatterns := []struct {
		pattern string
		reason  string
	}{
		{"EOF", "EOF from server"},
		{"timeout", "request timeout"},
		{"deadline exceeded", "deadline exceeded"},
		{"connection refused", "connection refused"},
		{"connection reset", "connection reset"},
		{"429", "rate limited (429)"},
		{"rate limit", "rate limited"},
		{"too many requests", "rate limited (too many requests)"},
		{"503", "service unavailable (503)"},
		{"service unavailable", "service unavailable"},
		{"bad gateway", "bad gateway"},
		{"no such host", "DNS resolution failed"},
		{"i/o timeout", "I/O timeout"},
		{"TLS handshake", "TLS handshake failure"},
	}

	lowerMsg := strings.ToLower(msg)
	for _, p := range transientPatterns {
		if strings.Contains(lowerMsg, strings.ToLower(p.pattern)) {
			return true, p.reason
		}
	}

	return false, ""
}

// SkipOnTransientError checks if an error is transient and skips the test if so.
// For permanent errors or nil errors, it does nothing.
// Use this in integration tests that call real APIs where transient failures
// should not fail the build.
func SkipOnTransientError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	transient, reason := ClassifyAPIError(err)
	if transient {
		t.Skipf("Skipping test due to transient error: %s (%v)", reason, err)
	}
}

// RequireNonZeroData validates that the given fields are non-zero.
// fields is a map of field name to value. If any field is zero, the test fails.
func RequireNonZeroData(t *testing.T, fields map[string]float64) {
	t.Helper()
	for name, value := range fields {
		if value == 0 {
			t.Errorf("field %q is zero, expected non-zero", name)
		}
	}
}

// RequireDataConsistency validates basic data consistency rules.
// For K-line data: High >= Low, High >= Open, High >= Close, Low <= Open, Low <= Close.
func RequireDataConsistency(t *testing.T, open, high, low, close float64) {
	t.Helper()
	if high < low {
		t.Errorf("High (%.4f) < Low (%.4f)", high, low)
	}
	if high < open {
		t.Errorf("High (%.4f) < Open (%.4f)", high, open)
	}
	if high < close {
		t.Errorf("High (%.4f) < Close (%.4f)", high, close)
	}
	if low > open {
		t.Errorf("Low (%.4f) > Open (%.4f)", low, open)
	}
	if low > close {
		t.Errorf("Low (%.4f) > Close (%.4f)", low, close)
	}
}

// CheckAPIResult is a convenience function that combines ClassifyAPIError and
// returns an APITestResult. It's useful for structured test reporting.
func CheckAPIResult(err error, dataValid bool) APITestResult {
	result := APITestResult{
		DataValid: dataValid,
		Error:     err,
	}
	if err == nil {
		result.Available = true
		return result
	}
	transient, reason := ClassifyAPIError(err)
	if transient {
		result.Available = false
		result.Skipped = true
		result.SkipReason = reason
	} else {
		result.Available = true // API responded but with a permanent error
	}
	return result
}

// FormatResult formats an APITestResult for display.
func FormatResult(r APITestResult) string {
	if r.Skipped {
		return fmt.Sprintf("SKIPPED: %s", r.SkipReason)
	}
	if r.Error != nil {
		return fmt.Sprintf("FAIL: %v", r.Error)
	}
	if r.DataValid {
		return "OK"
	}
	return "WARN: data validation failed"
}
