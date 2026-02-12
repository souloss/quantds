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
//	client := sina.NewClient(nil)
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
	"github.com/souloss/quantds/request"
)

const (
	BaseURL = "https://quotes.sina.cn"
)

type Client struct {
	http request.Client
}

type Option func(*Client)

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

func (c *Client) Close() {
	c.http.Close()
}
