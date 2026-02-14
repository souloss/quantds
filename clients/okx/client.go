package okx

import (
        "time"

        "github.com/failsafe-go/failsafe-go/timeout"
        "github.com/souloss/quantds/request"
)

const (
        DefaultBaseURL = "https://www.okx.com"
        AwsBaseURL     = "https://aws.okx.com"
)

type Client struct {
        http    request.Client
        BaseURL string
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

// WithBaseURL sets the base URL
func WithBaseURL(url string) Option {
        return func(c *Client) {
                c.BaseURL = url
        }
}

// NewClient creates a new OKX client
// If no options are provided, it uses the default configuration
func NewClient(opts ...Option) *Client {
        c := &Client{
                http:    request.NewClient(request.DefaultConfig()),
                BaseURL: DefaultBaseURL,
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
