package coingecko

import (
        "time"

        "github.com/failsafe-go/failsafe-go/timeout"
        "github.com/souloss/quantds/request"
)

const (
        BaseURL = "https://api.coingecko.com/api/v3"
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

// NewClient creates a new CoinGecko client
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
        if c.http != nil {
                c.http.Close()
        }
}
