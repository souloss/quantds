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
//	client := tencent.NewClient(nil)
//	defer client.Close()
//
//	result, record, err := client.GetKline(ctx, &tencent.KlineParams{
//	    Symbol: "000001.SZ",
//	    Period: "day",
//	    Count:  100,
//	})
package tencent

import (
	"github.com/souloss/quantds/request"
)

const (
	BaseURL = "https://web.sqt.gtimg.cn"
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
