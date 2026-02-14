// Package tencent provides a client for the Tencent Finance (腾讯证券) data source.
//
// Tencent Finance provides financial data through web APIs, offering:
//   - Historical K-line data
//   - Support for A-shares (SH, SZ), B-shares, and Beijing Stock Exchange
//
// API Features:
//   - No authentication required
//   - Multiple periods: 1m, 5m, 15m, 30m, 60m, daily, weekly, monthly
//   - Returns up to ~320 data points per request
//
// Limitations:
//   - No adjustment (qfq/hfq) support
//   - API format is non-standard JSON (requires parsing)
//   - High-frequency requests may be rate-limited
//
// Example:
//
//	client := tencent.NewClient()
//	defer client.Close()
//
//	result, record, err := client.GetKline(ctx, &tencent.KlineParams{
//	    Symbol: "000001.SZ",
//	    Period: "day",
//	    Count:  100,
//	})
package tencent

import (
	"time"

	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/souloss/quantds/request"
)

const (
	BaseURL = "https://web.sqt.gtimg.cn"

	// API Endpoints
	QuoteAPI     = "http://qt.gtimg.cn/q="
	MoneyFlowAPI = "http://qt.gtimg.cn/q=ff_"

	// Default Headers
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	DefaultReferer   = "https://gu.qq.com/"
)

type Client struct {
	http request.Client
}

type Option func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient request.Client) Option {
	return func(c *Client) {
		c.http = httpClient
	}
}

// WithTimeout sets the request timeout
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.http = request.NewClient(request.DefaultConfig(
			request.WithTimeout(timeout.New[request.Response](d)),
		))
	}
}

// WithConfig sets a custom request configuration
func WithConfig(cfg *request.Config) Option {
	return func(c *Client) {
		c.http = request.NewClient(cfg)
	}
}

// NewClient creates a new Tencent Finance client
// If no options are provided, it uses the default configuration
func NewClient(opts ...Option) *Client {
	c := &Client{
		http: request.NewClient(request.DefaultConfig()),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Close() {
	c.http.Close()
}
