package finnhub

import (
	"os"
	"time"

	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/souloss/quantds/request"
)

const (
	BaseURL         = "https://finnhub.io/api/v1"
	QuoteAPI        = "/quote"
	StockCandleAPI  = "/stock/candle"
	StockSymbolAPI  = "/stock/symbol"
	ForexRatesAPI   = "/forex/rates"
	ForexCandleAPI  = "/forex/candle"
	ForexSymbolAPI  = "/forex/symbol"
	CryptoCandleAPI = "/crypto/candle"
	CryptoSymbolAPI = "/crypto/symbol"
)

var DefaultHeaders = map[string]string{
	"User-Agent": "quantds/1.0",
	"Accept":     "application/json",
}

const (
	Res1  = "1"
	Res5  = "5"
	Res15 = "15"
	Res30 = "30"
	Res60 = "60"
	ResD  = "D"
	ResW  = "W"
	ResM  = "M"
)

type Client struct {
	http   request.Client
	apiKey string
}

type Option func(*Client)

func WithHTTPClient(httpClient request.Client) Option {
	return func(c *Client) {
		c.http = httpClient
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.http = request.NewClient(request.DefaultConfig(
			request.WithTimeout(timeout.New[request.Response](d)),
		))
	}
}

func WithConfig(cfg *request.Config) Option {
	return func(c *Client) {
		c.http = request.NewClient(cfg)
	}
}

func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		http:   request.NewClient(request.DefaultConfig()),
		apiKey: os.Getenv("FINNHUB_API_KEY"),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Close() {
	c.http.Close()
}

func ToResolution(tf string) string {
	switch tf {
	case "1m":
		return Res1
	case "5m":
		return Res5
	case "15m":
		return Res15
	case "30m":
		return Res30
	case "60m", "1h":
		return Res60
	case "1d", "":
		return ResD
	case "1w":
		return ResW
	case "1M":
		return ResM
	default:
		return ResD
	}
}
