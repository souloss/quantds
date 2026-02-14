package xueqiu

import (
        "time"

        "github.com/failsafe-go/failsafe-go/timeout"
        "github.com/souloss/quantds/request"
)

const (
        BaseURL = "https://stock.xueqiu.com"

        // Default Headers
        DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
        DefaultReferer   = "https://xueqiu.com/"
)

type Client struct {
        http   request.Client
        token  string
        cookie string
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

// WithToken sets the API token
func WithToken(token string) Option {
        return func(c *Client) {
                c.token = token
        }
}

// WithCookie sets the cookie for authentication
func WithCookie(cookie string) Option {
        return func(c *Client) {
                c.cookie = cookie
        }
}

// NewClient creates a new Xueqiu client
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
        c.http.Close()
}

// buildHeaders returns common request headers with optional auth
func (c *Client) buildHeaders() map[string]string {
        headers := map[string]string{
                "User-Agent": DefaultUserAgent,
                "Referer":    DefaultReferer,
        }
        if c.cookie != "" {
                headers["Cookie"] = c.cookie
        }
        if c.token != "" {
                headers["X-Token"] = c.token
        }
        return headers
}
