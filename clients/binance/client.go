package binance

// Package binance provides a client for the Binance cryptocurrency exchange API.
//
// Binance is one of the world's largest cryptocurrency exchanges, offering:
//   - K-line data (historical prices) for all trading pairs
//   - Real-time quotes (spot prices)
//   - Order book data
//   - 24-hour ticker data
//
// API Features:
//   - No authentication required for public APIs
//   - Supports all major cryptocurrencies
//   - Real-time data with minimal delay
//   - Multiple periods: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M
//
// Limitations:
//   - Rate limiting: 1200 requests per minute for public endpoints
//   - Weight-based rate limiting for some endpoints
//   - APIs may change without notice
//
// Example:
//
//      client := binance.NewClient()
//      defer client.Close()
//
//      // Get daily K-line for BTCUSDT
//      result, record, err := client.GetKline(ctx, &binance.KlineParams{
//          Symbol:    "BTCUSDT",
//          Interval:  "1d",
//          Limit:     100,
//      })

import (
        "time"

        "github.com/failsafe-go/failsafe-go/timeout"
        "github.com/souloss/quantds/request"
)

// API endpoints
const (
        BaseURL         = "https://api.binance.com"
        APIV3           = "/api/v3"
        KlineAPI        = "/api/v3/klines"
        Ticker24hrAPI   = "/api/v3/ticker/24hr"
        TickerPriceAPI  = "/api/v3/ticker/price"
        ExchangeInfoAPI = "/api/v3/exchangeInfo"
        DepthAPI        = "/api/v3/depth"
)

// HTTP headers for Binance API
var (
        DefaultHeaders = map[string]string{
                "User-Agent": "quantds/1.0",
                "Accept":     "application/json",
        }
)

// Interval constants for K-line data
const (
        Interval1m  = "1m"
        Interval3m  = "3m"
        Interval5m  = "5m"
        Interval15m = "15m"
        Interval30m = "30m"
        Interval1h  = "1h"
        Interval2h  = "2h"
        Interval4h  = "4h"
        Interval6h  = "6h"
        Interval8h  = "8h"
        Interval12h = "12h"
        Interval1d  = "1d"
        Interval3d  = "3d"
        Interval1w  = "1w"
        Interval1M  = "1M"
)

// Rate limit constants
const (
        RateLimitPerMinute = 1200
        MaxKlineLimit      = 1000
        MaxDepthLimit      = 5000
)

// Client is the Binance API client
type Client struct {
        http request.Client
}

// Option is a function that configures the client
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

// NewClient creates a new Binance client
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

// Close closes the underlying HTTP client
func (c *Client) Close() {
        c.http.Close()
}
