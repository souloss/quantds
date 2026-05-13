package testutil

import (
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
)

func TestClassifyAPIError_Nil(t *testing.T) {
	transient, reason := ClassifyAPIError(nil)
	if transient {
		t.Error("nil error should not be transient")
	}
	if reason != "" {
		t.Errorf("nil error should have empty reason, got %q", reason)
	}
}

func TestClassifyAPIError_Transient(t *testing.T) {
	tests := []struct {
		err           error
		wantTransient bool
	}{
		{io.EOF, true},
		{fmt.Errorf("EOF"), true},
		{fmt.Errorf("timeout waiting for response"), true},
		{fmt.Errorf("context deadline exceeded"), true},
		{fmt.Errorf("connection refused"), true},
		{fmt.Errorf("connection reset by peer"), true},
		{fmt.Errorf("429 Too Many Requests"), true},
		{fmt.Errorf("rate limit exceeded"), true},
		{fmt.Errorf("too many requests"), true},
		{fmt.Errorf("503 Service Unavailable"), true},
		{fmt.Errorf("service unavailable"), true},
		{fmt.Errorf("bad gateway"), true},
		{fmt.Errorf("no such host"), true},
		{fmt.Errorf("i/o timeout"), true},
		{fmt.Errorf("TLS handshake timeout"), true},
	}

	for _, tt := range tests {
		t.Run(tt.err.Error(), func(t *testing.T) {
			transient, reason := ClassifyAPIError(tt.err)
			if transient != tt.wantTransient {
				t.Errorf("ClassifyAPIError(%v) transient = %v, want %v", tt.err, transient, tt.wantTransient)
			}
			if tt.wantTransient && reason == "" {
				t.Errorf("ClassifyAPIError(%v) transient but reason is empty", tt.err)
			}
		})
	}
}

func TestClassifyAPIError_Permanent(t *testing.T) {
	tests := []struct {
		err error
	}{
		{errors.New("invalid symbol")},
		{errors.New("unauthorized")},
		{errors.New("forbidden")},
		{errors.New("bad request")},
		{errors.New("not found")},
		{fmt.Errorf("symbol required")},
	}

	for _, tt := range tests {
		t.Run(tt.err.Error(), func(t *testing.T) {
			transient, _ := ClassifyAPIError(tt.err)
			if transient {
				t.Errorf("ClassifyAPIError(%v) should not be transient", tt.err)
			}
		})
	}
}

type testNetError struct {
	msg string
}

func (e *testNetError) Error() string   { return e.msg }
func (e *testNetError) Timeout() bool   { return true }
func (e *testNetError) Temporary() bool { return true }

func TestClassifyAPIError_NetError(t *testing.T) {
	var err net.Error = &testNetError{msg: "network timeout"}
	transient, reason := ClassifyAPIError(err)
	if !transient {
		t.Error("net.Error should be transient")
	}
	if reason != "network error" {
		t.Errorf("reason = %q, want %q", reason, "network error")
	}
}

func TestRequireNonZeroData(t *testing.T) {
	// Test with all non-zero values - should pass
	t.Run("all non-zero", func(t *testing.T) {
		// Use a fake testing.T to capture errors
		fakeT := &testing.T{}
		RequireNonZeroData(fakeT, map[string]float64{"Open": 10.0, "Close": 11.0})
		if fakeT.Failed() {
			t.Error("should not fail with non-zero values")
		}
	})
}

func TestRequireDataConsistency(t *testing.T) {
	tests := []struct {
		name        string
		open, high, low, close float64
		shouldFail  bool
	}{
		{"valid", 10, 12, 9, 11, false},
		{"high < low", 10, 8, 12, 11, true},
		{"high < open", 15, 14, 9, 11, true},
		{"high < close", 10, 10, 9, 11, true},
		{"low > open", 10, 12, 11, 11, true},
		{"low > close", 10, 12, 11, 9, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeT := &testing.T{}
			RequireDataConsistency(fakeT, tt.open, tt.high, tt.low, tt.close)
			if fakeT.Failed() != tt.shouldFail {
				t.Errorf("RequireDataConsistency(%v, %v, %v, %v) failed=%v, want %v",
					tt.open, tt.high, tt.low, tt.close, fakeT.Failed(), tt.shouldFail)
			}
		})
	}
}

func TestCheckAPIResult(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := CheckAPIResult(nil, true)
		if !r.Available || !r.DataValid || r.Skipped {
			t.Errorf("unexpected result: %+v", r)
		}
	})

	t.Run("transient error", func(t *testing.T) {
		r := CheckAPIResult(fmt.Errorf("EOF"), false)
		if r.Available || !r.Skipped {
			t.Errorf("expected unavailable and skipped: %+v", r)
		}
	})

	t.Run("permanent error", func(t *testing.T) {
		r := CheckAPIResult(errors.New("invalid symbol"), false)
		if !r.Available || r.Skipped {
			t.Errorf("expected available and not skipped: %+v", r)
		}
	})
}

func TestFormatResult(t *testing.T) {
	tests := []struct {
		result APITestResult
		want   string
	}{
		{APITestResult{Available: true, DataValid: true}, "OK"},
		{APITestResult{Skipped: true, SkipReason: "EOF"}, "SKIPPED: EOF"},
		{APITestResult{Available: true, DataValid: false}, "WARN: data validation failed"},
		{APITestResult{Error: errors.New("bad")}, "FAIL: bad"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FormatResult(tt.result)
			if got != tt.want {
				t.Errorf("FormatResult() = %q, want %q", got, tt.want)
			}
		})
	}
}
