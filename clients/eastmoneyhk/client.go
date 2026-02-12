// Package eastmoneyhk provides a client for the EastMoney (东方财富) Hong Kong stock data API.
//
// EastMoney HK provides free APIs for Hong Kong stocks:
//   - K-line data (historical prices)
//   - Real-time quotes (spot prices)
//   - Stock lists and details
//
// API Features:
//   - No authentication required for public APIs
//   - Supports HK stocks (HKEX)
//   - Real-time data with minimal delay
//   - Multiple periods: 1m, 5m, 15m, 30m, 60m, daily, weekly, monthly
//
// Limitations:
//   - Rate limiting may occur with high-frequency requests
//   - APIs may change without notice
//
// Example:
//
//	client := eastmoneyhk.NewClient(nil)
//	defer client.Close()
//
//	// Get daily K-line for Tencent (00700.HK)
//	result, record, err := client.GetKline(ctx, &eastmoneyhk.KlineParams{
//	    Symbol:    "00700.HK",
//	    StartDate: "20240101",
//	    EndDate:   "20241231",
//	    Period:    "101",  // daily
//	})
package eastmoneyhk

import (
	"github.com/souloss/quantds/request"
)

// API endpoints
const (
	BaseURL   = "https://push2his.eastmoney.com"
	PushURL   = "https://push2.eastmoney.com"
	QuoteAPI  = "https://push2.eastmoney.com/api/qt/clist/get"
	DetailAPI = "https://emweb.eastmoney.com/PC_HKF10/NewFinanceAnalysis/Index"
)

// Market ID for HK stocks in EastMoney system
const (
	MarketHK = 116 // HK stock market ID in EastMoney
)

// Period constants for K-line data
const (
	Period1m  = "1"
	Period5m  = "5"
	Period15m = "15"
	Period30m = "30"
	Period60m = "60"
	Period1d  = "101"
	Period1w  = "102"
	Period1M  = "103"
)

// HTTP headers for EastMoney HK API
var (
	DefaultHeaders = map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Referer":    "https://quote.eastmoney.com/hk/",
		"Accept":     "application/json",
	}
)

// Client is the EastMoney HK API client
type Client struct {
	http request.Client
}

// Option is a function that configures the client
type Option func(*Client)

// NewClient creates a new EastMoney HK client
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
