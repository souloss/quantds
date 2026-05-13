package middleware

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/souloss/quantds/manager"
)

// CircuitBreakerState represents the state of a circuit breaker.
type CircuitBreakerState int

const (
	CircuitClosed   CircuitBreakerState = iota // Normal operation
	CircuitOpen                                 // Rejecting all requests
	CircuitHalfOpen                             // Testing if provider recovered
)

func (s CircuitBreakerState) String() string {
	switch s {
	case CircuitClosed:
		return "closed"
	case CircuitOpen:
		return "open"
	case CircuitHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerStateProvider is an optional interface that providers can implement
// to expose their circuit breaker state for monitoring.
type CircuitBreakerStateProvider interface {
	CircuitBreakerState() CircuitBreakerState
}

// ErrorClass represents how an error should be treated by the circuit breaker.
type ErrorClass int

const (
	// ErrorTransient indicates a temporary failure (EOF, timeout, rate limit)
	// that should trigger the circuit breaker.
	ErrorTransient ErrorClass = iota
	// ErrorPermanent indicates a non-recoverable failure (auth, invalid symbol)
	// that should NOT trigger the circuit breaker.
	ErrorPermanent
)

// ErrorClassifier determines whether an error should trigger the circuit breaker.
// Return ErrorTransient for temporary failures (EOF, timeout, 429, 503)
// and ErrorPermanent for permanent failures (401, 403, invalid symbol).
type ErrorClassifier func(error) ErrorClass

// DefaultErrorClassifier classifies common errors.
// EOF, timeout, rate limit (429), and service unavailable (503) are transient.
// All other errors are considered permanent (auth, bad request, etc.).
func DefaultErrorClassifier(err error) ErrorClass {
	if err == nil {
		return ErrorPermanent
	}
	// Network-level transient errors
	if err == io.EOF {
		return ErrorTransient
	}
	if _, ok := err.(net.Error); ok {
		return ErrorTransient
	}
	// Check for common transient error patterns in the error message
	msg := err.Error()
	switch {
	case containsAny(msg, "EOF", "timeout", "deadline exceeded", "connection refused", "connection reset"):
		return ErrorTransient
	case containsAny(msg, "429", "rate limit", "too many requests"):
		return ErrorTransient
	case containsAny(msg, "503", "service unavailable", "bad gateway"):
		return ErrorTransient
	default:
		return ErrorPermanent
	}
}

// containsAny checks if s contains any of the substrings.
func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}

// circuitBreaker tracks the state for a single provider's circuit breaker.
type circuitBreaker struct {
	failureThreshold int
	successThreshold int
	timeout          time.Duration
	classifier       ErrorClassifier

	mu              sync.Mutex
	state           CircuitBreakerState
	failures        int
	successes       int
	lastFailureTime time.Time
}

// State returns the current circuit breaker state.
func (cb *circuitBreaker) State() CircuitBreakerState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	// Check if we should transition from open to half-open
	if cb.state == CircuitOpen && time.Since(cb.lastFailureTime) >= cb.timeout {
		return CircuitHalfOpen
	}
	return cb.state
}

// CircuitBreaker creates a circuit breaker middleware with the given configuration.
// When a provider fails consecutively beyond failureThreshold, the circuit opens
// and subsequent requests are immediately rejected with ErrCircuitOpen.
// After the timeout period, the circuit enters half-open state and allows
// one request through to test if the provider has recovered.
func CircuitBreaker[Req, Resp any](failureThreshold, successThreshold int, timeout time.Duration) Middleware[Req, Resp] {
	return CircuitBreakerWithClassifier[Req, Resp](failureThreshold, successThreshold, timeout, nil)
}

// CircuitBreakerWithClassifier creates a circuit breaker middleware with a custom error classifier.
// The classifier determines whether an error should trigger the circuit breaker.
// If classifier is nil, DefaultErrorClassifier is used.
func CircuitBreakerWithClassifier[Req, Resp any](failureThreshold, successThreshold int, timeout time.Duration, classifier ErrorClassifier) Middleware[Req, Resp] {
	if classifier == nil {
		classifier = DefaultErrorClassifier
	}
	return func(next manager.Provider[Req, Resp]) manager.Provider[Req, Resp] {
		cb := &circuitBreaker{
			failureThreshold: failureThreshold,
			successThreshold: successThreshold,
			timeout:          timeout,
			classifier:       classifier,
			state:            CircuitClosed,
		}
		return &providerFunc[Req, Resp]{
			name:             next.Name(),
			supportedMarkets: next.SupportedMarkets(),
			canHandle:        next.CanHandle,
			cbState:          cb.State,
			fetch: func(ctx context.Context, req Req) (Resp, *manager.RequestTrace, error) {
				if !cb.allowRequest() {
					var zero Resp
					return zero, nil, manager.ErrCircuitOpen
				}

				resp, trace, err := next.Fetch(ctx, req)
				if err != nil {
					if classifier(err) == ErrorTransient {
						cb.recordFailure()
					}
					return resp, trace, err
				}
				cb.recordSuccess()
				return resp, trace, nil
			},
		}
	}
}

// CircuitBreakerConfig holds configuration for circuit breaker middleware.
type CircuitBreakerConfig struct {
	Enabled          bool
	FailureThreshold int
	SuccessThreshold int
	Timeout          time.Duration
}

// CircuitBreakerFromConfig creates a circuit breaker middleware from config.
// If cfg.Enabled is false, returns an identity middleware (no-op).
func CircuitBreakerFromConfig[Req, Resp any](cfg CircuitBreakerConfig) Middleware[Req, Resp] {
	if !cfg.Enabled {
		return Identity[Req, Resp]()
	}
	return CircuitBreaker[Req, Resp](cfg.FailureThreshold, cfg.SuccessThreshold, cfg.Timeout)
}

// Identity returns a no-op middleware that passes through to the next provider.
func Identity[Req, Resp any]() Middleware[Req, Resp] {
	return func(next manager.Provider[Req, Resp]) manager.Provider[Req, Resp] {
		return next
	}
}

// allowRequest checks if a request should be allowed through the circuit breaker.
func (cb *circuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if time.Since(cb.lastFailureTime) >= cb.timeout {
			cb.state = CircuitHalfOpen
			cb.failures = 0
			cb.successes = 0
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return true
	}
}

// recordFailure records a failure and potentially opens the circuit.
func (cb *circuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case CircuitClosed:
		if cb.failures >= cb.failureThreshold {
			cb.state = CircuitOpen
		}
	case CircuitHalfOpen:
		cb.state = CircuitOpen
	}
}

// recordSuccess records a success and potentially closes the circuit.
func (cb *circuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitClosed:
		cb.failures = 0
	case CircuitHalfOpen:
		cb.successes++
		if cb.successes >= cb.successThreshold {
			cb.state = CircuitClosed
			cb.failures = 0
			cb.successes = 0
		}
	}
}