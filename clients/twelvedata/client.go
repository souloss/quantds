package twelvedata

import (
	"os"
	"time"

	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/souloss/quantds/request"
)

const (
	BaseURL        = "https://api.twelvedata.com"
	TimeSeriesAPI  = "/time_series"
	QuoteAPI       = "/quote"
	PriceAPI       = "/price"
	StocksAPI      = "/stocks"
	ForexPairsAPI  = "/forex_pairs"
	CryptoAPI      = "/cryptocurrencies"
	ETFAPI         = "/etf"
	FundsAPI       = "/funds"
	BondsAPI       = "/bonds"
)

var DefaultHeaders = map[string]string{
	"User-Agent": "quantds/1.0",
	"Accept":     "application/json",
}

const (
	Interval1min  = "1min"
	Interval5min  = "5min"
	Interval15min = "15min"
	Interval30min = "30min"
	Interval1h    = "1h"
	Interval1day  = "1day"
	Interval1week = "1week"
	Interval1month = "1month"
)

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
		apiKey: os.Getenv("TWELVEDATA_API_KEY"),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Close() { c.http.Close() }

func ToInterval(tf string) string {
	switch tf {
	case "1m":
		return Interval1min
	case "5m":
		return Interval5min
	case "15m":
		return Interval15min
	case "30m":
		return Interval30min
	case "60m", "1h":
		return Interval1h
	case "1d", "":
		return Interval1day
	case "1w":
		return Interval1week
	case "1M":
		return Interval1month
	default:
		return Interval1day
	}
}
