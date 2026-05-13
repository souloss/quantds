package manager

import "errors"

var (
	ErrNoProvider        = errors.New("no provider available")
	ErrAllProviderFailed = errors.New("all providers failed")
	ErrCircuitOpen       = errors.New("circuit breaker is open")
	ErrRateLimited       = errors.New("rate limit exceeded")
)
