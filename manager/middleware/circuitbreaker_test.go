package middleware

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/manager"
)

// mockProvider is a test helper that implements manager.Provider.
type mockProvider[Req, Resp any] struct {
	name    string
	fetch   func(ctx context.Context, req Req) (Resp, *manager.RequestTrace, error)
	markets []domain.Market
}

func (m *mockProvider[Req, Resp]) Name() string                      { return m.name }
func (m *mockProvider[Req, Resp]) SupportedMarkets() []domain.Market { return m.markets }
func (m *mockProvider[Req, Resp]) CanHandle(symbol string) bool      { return true }
func (m *mockProvider[Req, Resp]) Fetch(ctx context.Context, req Req) (Resp, *manager.RequestTrace, error) {
	return m.fetch(ctx, req)
}

// errTransient is a transient error that should trigger the circuit breaker.
var errTransient = fmt.Errorf("timeout: connection refused")

// errPermanent is a permanent error that should NOT trigger the circuit breaker.
var errPermanent = errors.New("invalid symbol")

func TestCircuitBreaker_ClosedState(t *testing.T) {
	cb := CircuitBreaker[string, string](3, 2, 1*time.Second)

	callCount := 0
	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			callCount++
			return "ok", nil, nil
		},
	}

	wrapped := cb(next)

	// In closed state, all requests should pass through
	for i := 0; i < 5; i++ {
		resp, _, err := wrapped.Fetch(context.Background(), "req")
		if err != nil {
			t.Errorf("request %d: unexpected error: %v", i, err)
		}
		if resp != "ok" {
			t.Errorf("request %d: got %q, want %q", i, resp, "ok")
		}
	}
	if callCount != 5 {
		t.Errorf("expected 5 calls, got %d", callCount)
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cb := CircuitBreaker[string, string](3, 2, 1*time.Second)

	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			return "", nil, errTransient
		},
	}

	wrapped := cb(next)

	// Fail 3 times to trigger open (transient errors trigger CB)
	for i := 0; i < 3; i++ {
		_, _, err := wrapped.Fetch(context.Background(), "req")
		if err != errTransient {
			t.Errorf("failure %d: got error %v, want %v", i, err, errTransient)
		}
	}

	// Next request should be rejected with ErrCircuitOpen
	_, _, err := wrapped.Fetch(context.Background(), "req")
	if err != manager.ErrCircuitOpen {
		t.Errorf("after threshold: got error %v, want ErrCircuitOpen", err)
	}
}

func TestCircuitBreaker_PermanentErrorDoesNotOpen(t *testing.T) {
	cb := CircuitBreaker[string, string](2, 2, 1*time.Second)

	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			return "", nil, errPermanent
		},
	}

	wrapped := cb(next)

	// Permanent errors should NOT trigger the circuit breaker
	for i := 0; i < 5; i++ {
		_, _, err := wrapped.Fetch(context.Background(), "req")
		if err != errPermanent {
			t.Errorf("permanent failure %d: got error %v, want %v", i, err, errPermanent)
		}
	}

	// Circuit should still be closed — requests still pass through
	_, _, err := wrapped.Fetch(context.Background(), "req")
	if err != errPermanent {
		t.Errorf("after permanent failures: got error %v, want errPermanent (circuit should still be closed)", err)
	}
}

func TestCircuitBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cb := CircuitBreaker[string, string](2, 2, 100*time.Millisecond)

	failCount := 0
	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			failCount++
			if failCount <= 2 {
				return "", nil, errTransient
			}
			return "ok", nil, nil
		},
	}

	wrapped := cb(next)

	// Fail 2 times to open the circuit
	wrapped.Fetch(context.Background(), "req")
	wrapped.Fetch(context.Background(), "req")

	// Should be open now
	_, _, err := wrapped.Fetch(context.Background(), "req")
	if err != manager.ErrCircuitOpen {
		t.Errorf("before timeout: got error %v, want ErrCircuitOpen", err)
	}

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	// Should be half-open now, request goes through
	resp, _, err := wrapped.Fetch(context.Background(), "req")
	if err != nil {
		t.Errorf("half-open: unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Errorf("half-open: got %q, want %q", resp, "ok")
	}
}

func TestCircuitBreaker_HalfOpenFailureReopens(t *testing.T) {
	cb := CircuitBreaker[string, string](2, 2, 100*time.Millisecond)

	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			return "", nil, errTransient
		},
	}

	wrapped := cb(next)

	// Open the circuit
	wrapped.Fetch(context.Background(), "req")
	wrapped.Fetch(context.Background(), "req")

	// Wait for timeout to enter half-open
	time.Sleep(150 * time.Millisecond)

	// Fail in half-open state should reopen circuit
	wrapped.Fetch(context.Background(), "req")

	// Should be open again
	_, _, err := wrapped.Fetch(context.Background(), "req")
	if err != manager.ErrCircuitOpen {
		t.Errorf("after half-open failure: got error %v, want ErrCircuitOpen", err)
	}
}

func TestCircuitBreaker_HalfOpenSuccessClosesAfterThreshold(t *testing.T) {
	cb := CircuitBreaker[string, string](2, 2, 100*time.Millisecond)

	callCount := 0
	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			callCount++
			if callCount <= 2 {
				return "", nil, errTransient
			}
			return "ok", nil, nil
		},
	}

	wrapped := cb(next)

	// Open the circuit
	wrapped.Fetch(context.Background(), "req")
	wrapped.Fetch(context.Background(), "req")

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	// Succeed enough times in half-open to close the circuit
	for i := 0; i < 2; i++ {
		resp, _, err := wrapped.Fetch(context.Background(), "req")
		if err != nil {
			t.Errorf("half-open success %d: unexpected error: %v", i, err)
		}
		if resp != "ok" {
			t.Errorf("half-open success %d: got %q, want %q", i, resp, "ok")
		}
	}

	// Circuit should be closed now, requests pass through normally
	for i := 0; i < 5; i++ {
		resp, _, err := wrapped.Fetch(context.Background(), "req")
		if err != nil {
			t.Errorf("closed state %d: unexpected error: %v", i, err)
		}
		if resp != "ok" {
			t.Errorf("closed state %d: got %q, want %q", i, resp, "ok")
		}
	}
}

func TestCircuitBreaker_SuccessResetsFailures(t *testing.T) {
	cb := CircuitBreaker[string, string](3, 2, 1*time.Second)

	shouldFail := false
	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			if shouldFail {
				return "", nil, errTransient
			}
			return "ok", nil, nil
		},
	}

	wrapped := cb(next)

	// 2 failures (below threshold of 3)
	shouldFail = true
	wrapped.Fetch(context.Background(), "req")
	wrapped.Fetch(context.Background(), "req")

	// A success resets the failure counter
	shouldFail = false
	resp, _, err := wrapped.Fetch(context.Background(), "req")
	if err != nil {
		t.Errorf("success after failures: unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Errorf("success after failures: got %q, want %q", resp, "ok")
	}

	// Now 2 more failures — but since counter was reset, circuit stays closed
	shouldFail = true
	wrapped.Fetch(context.Background(), "req")
	wrapped.Fetch(context.Background(), "req")

	// Circuit should still be closed (only 2 consecutive failures, not 3)
	shouldFail = false
	resp, _, err = wrapped.Fetch(context.Background(), "req")
	if err != nil {
		t.Errorf("should still be closed: unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Errorf("should still be closed: got %q, want %q", resp, "ok")
	}
}

func TestCircuitBreaker_PreservesProviderInfo(t *testing.T) {
	cb := CircuitBreaker[string, string](3, 2, 1*time.Second)

	next := &mockProvider[string, string]{
		name:    "myprovider",
		fetch:   func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "ok", nil, nil },
		markets: []domain.Market{domain.MarketCN, domain.MarketUS},
	}

	wrapped := cb(next)

	if wrapped.Name() != "myprovider" {
		t.Errorf("Name() = %q, want %q", wrapped.Name(), "myprovider")
	}
	markets := wrapped.SupportedMarkets()
	if len(markets) != 2 {
		t.Errorf("SupportedMarkets() = %v, want 2 markets", markets)
	}
	if !wrapped.CanHandle("anything") {
		t.Error("CanHandle() = false, want true")
	}
}

func TestCircuitBreakerFromConfig_Disabled(t *testing.T) {
	cfg := CircuitBreakerConfig{Enabled: false}
	cb := CircuitBreakerFromConfig[string, string](cfg)

	next := &mockProvider[string, string]{
		name:  "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "ok", nil, nil },
	}

	wrapped := cb(next)

	// Identity middleware should just pass through
	resp, _, err := wrapped.Fetch(context.Background(), "req")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Errorf("got %q, want %q", resp, "ok")
	}
}

func TestCircuitBreakerFromConfig_Enabled(t *testing.T) {
	cfg := CircuitBreakerConfig{
		Enabled:          true,
		FailureThreshold: 2,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
	}
	cb := CircuitBreakerFromConfig[string, string](cfg)

	next := &mockProvider[string, string]{
		name:  "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "", nil, errTransient },
	}

	wrapped := cb(next)

	// Fail to open
	wrapped.Fetch(context.Background(), "req")
	wrapped.Fetch(context.Background(), "req")

	_, _, err := wrapped.Fetch(context.Background(), "req")
	if err != manager.ErrCircuitOpen {
		t.Errorf("got error %v, want ErrCircuitOpen", err)
	}
}

func TestIdentity(t *testing.T) {
	id := Identity[string, string]()

	next := &mockProvider[string, string]{
		name:  "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "ok", nil, nil },
	}

	wrapped := id(next)

	// Identity should pass through to the underlying provider
	resp, _, err := wrapped.Fetch(context.Background(), "req")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Errorf("got %q, want %q", resp, "ok")
	}
	if wrapped.Name() != "test" {
		t.Errorf("Name() = %q, want %q", wrapped.Name(), "test")
	}
}

func TestDefaultErrorClassifier(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantClass ErrorClass
	}{
		{"nil error", nil, ErrorPermanent},
		{"EOF", io.EOF, ErrorTransient},
		{"timeout message", fmt.Errorf("timeout waiting for response"), ErrorTransient},
		{"deadline exceeded", fmt.Errorf("deadline exceeded"), ErrorTransient},
		{"connection refused", fmt.Errorf("connection refused"), ErrorTransient},
		{"connection reset", fmt.Errorf("connection reset by peer"), ErrorTransient},
		{"429 rate limit", fmt.Errorf("429 rate limit exceeded"), ErrorTransient},
		{"rate limit message", fmt.Errorf("rate limit exceeded"), ErrorTransient},
		{"too many requests", fmt.Errorf("too many requests"), ErrorTransient},
		{"503 service unavailable", fmt.Errorf("503 service unavailable"), ErrorTransient},
		{"service unavailable message", fmt.Errorf("service unavailable"), ErrorTransient},
		{"bad gateway", fmt.Errorf("bad gateway"), ErrorTransient},
		{"generic error", errors.New("some random error"), ErrorPermanent},
		{"invalid symbol", errors.New("invalid symbol"), ErrorPermanent},
		{"auth error", errors.New("unauthorized"), ErrorPermanent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultErrorClassifier(tt.err)
			if got != tt.wantClass {
				t.Errorf("DefaultErrorClassifier(%v) = %v, want %v", tt.err, got, tt.wantClass)
			}
		})
	}
}

func TestCircuitBreakerWithClassifier_CustomClassifier(t *testing.T) {
	// Custom classifier that treats ALL errors as transient
	allTransient := func(err error) ErrorClass { return ErrorTransient }

	cb := CircuitBreakerWithClassifier[string, string](2, 1, 100*time.Millisecond, allTransient)

	errFail := errors.New("any error") // normally permanent, but classifier overrides
	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			return "", nil, errFail
		},
	}

	wrapped := cb(next)

	// With custom classifier treating all as transient, circuit should open
	wrapped.Fetch(context.Background(), "req")
	wrapped.Fetch(context.Background(), "req")

	_, _, err := wrapped.Fetch(context.Background(), "req")
	if err != manager.ErrCircuitOpen {
		t.Errorf("with custom classifier: got error %v, want ErrCircuitOpen", err)
	}
}

func TestCircuitBreakerState_String(t *testing.T) {
	tests := []struct {
		state CircuitBreakerState
		want  string
	}{
		{CircuitClosed, "closed"},
		{CircuitOpen, "open"},
		{CircuitHalfOpen, "half-open"},
		{CircuitBreakerState(99), "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("CircuitBreakerState(%d).String() = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}

func TestCircuitBreaker_HalfOpenPermanentErrorDoesNotReopen(t *testing.T) {
	cb := CircuitBreaker[string, string](2, 2, 100*time.Millisecond)

	callCount := 0
	next := &mockProvider[string, string]{
		name: "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) {
			callCount++
			if callCount <= 2 {
				return "", nil, errTransient // transient to open circuit
			}
			return "", nil, errPermanent // permanent in half-open should NOT reopen
		},
	}

	wrapped := cb(next)

	// Open the circuit with transient errors
	wrapped.Fetch(context.Background(), "req")
	wrapped.Fetch(context.Background(), "req")

	// Wait for timeout to enter half-open
	time.Sleep(150 * time.Millisecond)

	// Permanent error in half-open should NOT reopen the circuit
	_, _, err := wrapped.Fetch(context.Background(), "req")
	if err != errPermanent {
		t.Errorf("half-open permanent error: got %v, want errPermanent", err)
	}

	// Circuit should still be half-open (not reopened), next request should go through
	_, _, err = wrapped.Fetch(context.Background(), "req")
	// Since the mock still returns permanent errors, we get errPermanent but circuit stays half-open
	if err != errPermanent {
		t.Errorf("after half-open permanent: got %v, want errPermanent", err)
	}
}