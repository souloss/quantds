// Package yahoo provides a client for the Yahoo Finance data source.
//
// Yahoo Finance is one of the most popular financial data providers, offering free APIs for:
//   - K-line data (historical prices) for US stocks, ETFs, and indices
//   - Real-time quotes (spot prices)
//   - Stock details and company information
//
// API Features:
//   - No authentication required for public APIs
//   - Supports US stocks (NYSE, NASDAQ, AMEX)
//   - Real-time data with minimal delay
//   - Multiple periods: 1m, 5m, 15m, 30m, 60m, daily, weekly, monthly
//
// Limitations:
//   - Rate limiting may occur with high-frequency requests
//   - Some endpoints may require authentication for extended data
//   - APIs may change without notice
//
// Example:
//
//	client := yahoo.NewClient(nil)
//	defer client.Close()
//
//	// Get daily K-line for Apple (AAPL)
//	result, record, err := client.GetKline(ctx, &yahoo.KlineParams{
//	    Symbol:    "AAPL",
//	    Interval:  "1d",
//	    StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
//	    EndDate:   time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
//	})
package yahoo

import (
	"github.com/souloss/quantds/request"
)

// API endpoints
const (
	BaseURL    = "https://query1.finance.yahoo.com"
	ChartAPI   = "/v8/finance/chart"
	QuoteAPI   = "/v7/finance/quote"
	SearchAPI  = "/v1/finance/search"
	SummaryAPI = "/v10/finance/quoteSummary"
)

// HTTP headers for Yahoo Finance API
var (
	DefaultHeaders = map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Accept":     "application/json",
	}
)

// Interval constants for K-line data
const (
	Interval1m  = "1m"
	Interval5m  = "5m"
	Interval15m = "15m"
	Interval30m = "30m"
	Interval60m = "60m"
	Interval1d  = "1d"
	Interval1w  = "1wk"
	Interval1M  = "1mo"
)

// Range constants for K-line data
const (
	Range1d  = "1d"
	Range5d  = "5d"
	Range1mo = "1mo"
	Range3mo = "3mo"
	Range6mo = "6mo"
	Range1y  = "1y"
	Range2y  = "2y"
	Range5y  = "5y"
	Range10y = "10y"
	RangeYtd = "ytd"
	RangeMax = "max"
)

// Client is the Yahoo Finance API client
type Client struct {
	http request.Client
}

// Option is a function that configures the client
type Option func(*Client)

// NewClient creates a new Yahoo Finance client
func NewClient(httpClient request.Client, opts ...Option) *Client {
	if httpClient == nil {
		httpClient = request.NewClient(request.DefaultConfig())
	}
	c := &Client{
		http: httpClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Close closes the underlying HTTP client
func (c *Client) Close() {
	c.http.Close()
}
