// Package eastmoney provides a client for the EastMoney (东方财富) data source.
//
// EastMoney is one of China's largest financial data providers, offering free APIs for:
//   - K-line data (historical prices)
//   - Real-time quotes (spot prices)
//   - Stock lists and details
//   - News and announcements
//   - Financial reports
//
// API Features:
//   - No authentication required for public APIs
//   - Supports A-shares (SH, SZ), B-shares, and Beijing Stock Exchange
//   - Real-time data with minimal delay (~3 seconds)
//   - Multiple periods: 1m, 5m, 15m, 30m, 60m, daily, weekly, monthly
//   - Adjustment types: none, forward (qfq), backward (hfq)
//
// Limitations:
//   - Single request returns max ~500 data points
//   - Intraday data only available for recent 1-3 months
//   - High-frequency requests may be rate-limited
//   - APIs may change without notice
//
// Example:
//
//	client := eastmoney.NewClient(nil)
//	defer client.Close()
//
//	// Get daily K-line for 平安银行 (000001.SZ)
//	result, record, err := client.GetKline(ctx, &eastmoney.KlineParams{
//	    Symbol:    "000001.SZ",
//	    StartDate: "20240101",
//	    EndDate:   "20241231",
//	    Period:    "101",  // daily
//	    Adjust:    "1",    // forward adjustment
//	})
//
//	// Get real-time quotes
//	quotes, record, err := client.GetSpot(ctx, &eastmoney.SpotParams{
//	    Market:   "SZ",
//	    PageSize: 100,
//	})
package eastmoney

import (
	"github.com/souloss/quantds/request"
)

const (
	BaseURL    = "https://push2his.eastmoney.com"
	PushURL    = "https://push2.eastmoney.com"
	Datacenter = "https://datacenter.eastmoney.com"
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
