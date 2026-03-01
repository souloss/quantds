package eodhd

import (
	"os"
	"time"

	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/souloss/quantds/request"
)

const (
	BaseURL     = "https://eodhd.com/api"
	EODAPI      = "/eod"
	RealTimeAPI = "/real-time"
	ExchangeSymbolsAPI = "/exchange-symbol-list"
	BondFundamentalsAPI = "/bond-fundamentals"
)

var DefaultHeaders = map[string]string{
	"User-Agent": "quantds/1.0",
	"Accept":     "application/json",
}

type Client struct {
	http   request.Client
	apiKey string
}

type Option func(*Client)

func WithHTTPClient(httpClient request.Client) Option {
	return func(c *Client) { c.http = httpClient }
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.http = request.NewClient(request.DefaultConfig(
			request.WithTimeout(timeout.New[request.Response](d)),
		))
	}
}

func WithConfig(cfg *request.Config) Option {
	return func(c *Client) { c.http = request.NewClient(cfg) }
}

func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = key }
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		http:   request.NewClient(request.DefaultConfig()),
		apiKey: os.Getenv("EODHD_API_KEY"),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Close() { c.http.Close() }
