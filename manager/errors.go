package manager

import "errors"

var (
	ErrNoProvider        = errors.New("no provider available")
	ErrAllProviderFailed = errors.New("all providers failed")
)
