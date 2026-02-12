package request

import (
	"errors"
	"testing"
)

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		statusCode int
		wantType   ErrorType
	}{
		{
			name:       "nil error with 200 status",
			err:        nil,
			statusCode: 200,
			wantType:   ErrorTypeNone,
		},
		{
			name:       "nil error with 500 status",
			err:        nil,
			statusCode: 500,
			wantType:   ErrorTypeServer,
		},
		{
			name:       "nil error with 429 status",
			err:        nil,
			statusCode: 429,
			wantType:   ErrorTypeRateLimited,
		},
		{
			name:       "nil error with 401 status",
			err:        nil,
			statusCode: 401,
			wantType:   ErrorTypeAuth,
		},
		{
			name:       "nil error with 403 status",
			err:        nil,
			statusCode: 403,
			wantType:   ErrorTypeAuth,
		},
		{
			name:       "nil error with 400 status",
			err:        nil,
			statusCode: 400,
			wantType:   ErrorTypeClient,
		},
		{
			name:       "net timeout error",
			err:        &netTimeoutError{},
			statusCode: 0,
			wantType:   ErrorTypeTimeout,
		},
		{
			name:       "timeout in message",
			err:        errors.New("request timeout"),
			statusCode: 0,
			wantType:   ErrorTypeTimeout,
		},
		{
			name:       "deadline exceeded in message",
			err:        errors.New("context deadline exceeded"),
			statusCode: 0,
			wantType:   ErrorTypeTimeout,
		},
		{
			name:       "unknown error",
			err:        errors.New("some random error"),
			statusCode: 0,
			wantType:   ErrorTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyError(tt.err, tt.statusCode)
			if tt.wantType == ErrorTypeNone {
				if got != nil {
					t.Errorf("ClassifyError() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Errorf("ClassifyError() = nil, want type %s", tt.wantType)
				return
			}
			if got.Type != tt.wantType {
				t.Errorf("ClassifyError().Type = %v, want %v", got.Type, tt.wantType)
			}
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		retryable bool
	}{
		{
			name:      "nil error",
			err:       nil,
			retryable: false,
		},
		{
			name:      "timeout error",
			err:       &netTimeoutError{},
			retryable: true,
		},
		{
			name:      "network error",
			err:       &netOpError{},
			retryable: true,
		},
		{
			name:      "unknown error",
			err:       errors.New("unknown"),
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRetryableError(tt.err)
			if got != tt.retryable {
				t.Errorf("IsRetryableError() = %v, want %v", got, tt.retryable)
			}
		})
	}
}

func TestRequestError(t *testing.T) {
	cause := errors.New("underlying error")
	reqErr := &RequestError{
		Type:    ErrorTypeNetwork,
		Message: "connection failed",
		Cause:   cause,
	}

	if reqErr.Error() != "connection failed: underlying error" {
		t.Errorf("Error() = %q, want %q", reqErr.Error(), "connection failed: underlying error")
	}

	if reqErr.Unwrap() != cause {
		t.Errorf("Unwrap() = %v, want %v", reqErr.Unwrap(), cause)
	}
}

type netTimeoutError struct{}

func (e *netTimeoutError) Error() string   { return "timeout" }
func (e *netTimeoutError) Timeout() bool   { return true }
func (e *netTimeoutError) Temporary() bool { return true }

type netOpError struct{}

func (e *netOpError) Error() string   { return "network error" }
func (e *netOpError) Timeout() bool   { return false }
func (e *netOpError) Temporary() bool { return false }
func (e *netOpError) Unwrap() error   { return nil }
