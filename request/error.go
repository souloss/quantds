package request

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

type ErrorType string

const (
	ErrorTypeNone        ErrorType = ""
	ErrorTypeNetwork     ErrorType = "network"
	ErrorTypeTimeout     ErrorType = "timeout"
	ErrorTypeRateLimited ErrorType = "rate_limited"
	ErrorTypeAuth        ErrorType = "auth"
	ErrorTypeServer      ErrorType = "server"
	ErrorTypeClient      ErrorType = "client"
	ErrorTypeUnknown     ErrorType = "unknown"
)

type RequestError struct {
	Type       ErrorType
	Message    string
	StatusCode int
	Cause      error
}

func (e *RequestError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *RequestError) Unwrap() error {
	return e.Cause
}

func (e *RequestError) Is(target error) bool {
	t, ok := target.(*RequestError)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

func ClassifyError(err error, statusCode int) *RequestError {
	if err == nil {
		if statusCode >= 500 {
			return &RequestError{Type: ErrorTypeServer, StatusCode: statusCode, Message: "server error"}
		}
		if statusCode == 429 {
			return &RequestError{Type: ErrorTypeRateLimited, StatusCode: statusCode, Message: "rate limited"}
		}
		if statusCode == 401 || statusCode == 403 {
			return &RequestError{Type: ErrorTypeAuth, StatusCode: statusCode, Message: "authentication error"}
		}
		if statusCode >= 400 {
			return &RequestError{Type: ErrorTypeClient, StatusCode: statusCode, Message: "client error"}
		}
		return nil
	}

	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		return reqErr
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return &RequestError{Type: ErrorTypeTimeout, Cause: err}
		}
		return &RequestError{Type: ErrorTypeNetwork, Cause: err}
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		if urlErr.Timeout() {
			return &RequestError{Type: ErrorTypeTimeout, Cause: err}
		}
		return &RequestError{Type: ErrorTypeNetwork, Cause: err}
	}

	if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
		return &RequestError{Type: ErrorTypeTimeout, Cause: err}
	}

	return &RequestError{Type: ErrorTypeUnknown, Cause: err}
}

func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	reqErr := ClassifyError(err, 0)
	switch reqErr.Type {
	case ErrorTypeNetwork, ErrorTypeTimeout, ErrorTypeServer, ErrorTypeRateLimited:
		return true
	default:
		return false
	}
}

func NewError(typ ErrorType, message string, cause error) *RequestError {
	return &RequestError{
		Type:    typ,
		Message: message,
		Cause:   cause,
	}
}

func NewErrorWithStatus(typ ErrorType, message string, statusCode int, cause error) *RequestError {
	return &RequestError{
		Type:       typ,
		Message:    message,
		StatusCode: statusCode,
		Cause:      cause,
	}
}
