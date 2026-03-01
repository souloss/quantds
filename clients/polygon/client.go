package polygon

import (
	"os"
	"time"

	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/souloss/quantds/request"
)

const (
	BaseURL       = "https://api.polygon.io"
	AggregatesAPI = "/v2/aggs/ticker"
	TickersAPI    = "/v3/reference/tickers"
	SnapshotAPI   = "/v2/snapshot/locale/us/markets/stocks/tickers"
)

var DefaultHeaders = map[string]string{
	"User-Agent": "quantds/1.0",
	"Accept":     "application/json",
}

const (
	TimespanMinute = "minute"
	TimespanHour   = "hour"
	TimespanDay    = "day"
	TimespanWeek   = "week"
	TimespanMonth  = "month"
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
		apiKey: os.Getenv("POLYGON_API_KEY"),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Close() { c.http.Close() }

func ToTimespan(tf string) (string, int) {
	switch tf {
	case "1m":
		return TimespanMinute, 1
	case "5m":
		return TimespanMinute, 5
	case "15m":
		return TimespanMinute, 15
	case "30m":
		return TimespanMinute, 30
	case "60m", "1h":
		return TimespanHour, 1
	case "1w":
		return TimespanWeek, 1
	case "1M":
		return TimespanMonth, 1
	default:
		return TimespanDay, 1
	}
}
