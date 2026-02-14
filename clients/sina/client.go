// Package sina provides a client for the Sina Finance (新浪财经) data source.
//
// Sina Finance is one of China's most popular financial data platforms, offering:
//   - Real-time stock quotes (spot)
//   - Historical K-line data
//   - Support for A-shares (SH, SZ), B-shares, and Beijing Stock Exchange
//
// API Features:
//   - No authentication required
//   - Simple and reliable APIs
//   - Real-time data with minimal delay
//   - Multiple periods: 5m, 15m, 30m, 60m, daily, weekly, monthly
//
// Limitations:
//   - Spot API returns GBK encoded data (requires conversion)
//   - K-line API returns max ~500 data points per request
//   - High-frequency requests may be rate-limited
//   - No adjustment (qfq/hfq) support
//
// Example:
//
//	client := sina.NewClient()
//	defer client.Close()
//
//	// Get K-line data
//	result, record, err := client.GetKline(ctx, &sina.KlineParams{
//	    Symbol: "000001.SZ",
//	    Period: "d",
//	})
//
//	// Get real-time quotes
//	quotes, record, err := client.GetSpot(ctx, &sina.SpotParams{
//	    Symbols: []string{"000001.SZ", "600001.SH"},
//	})
package sina

import (
	"time"

	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/souloss/quantds/request"
)

const (
	BaseURL = "https://quotes.sina.cn"

	// API Endpoints
	SpotAPI      = "https://hq.sinajs.cn"
	MoneyFlowAPI = "http://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/MoneyFlow.ssl_qsfx_zjll_node"

	// Default Headers
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	DefaultReferer   = "https://finance.sina.com.cn/"
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

// NewClient creates a new Sina Finance client
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
