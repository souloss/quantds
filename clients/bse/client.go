package bse

import (
        "time"

        "github.com/failsafe-go/failsafe-go/timeout"
        "github.com/souloss/quantds/request"
)

const (
        BaseURL = "https://www.bse.cn/nqxxController/nqxxCnzq.do"
        Referer = "https://www.bse.cn/nq/listedcompany.html"
)

type Client struct {
        http    request.Client
        baseURL string
        headers map[string]string
}

type Option func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient request.Client) Option {
        return func(c *Client) { c.http = httpClient }
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
        return func(c *Client) { c.http = request.NewClient(cfg) }
}

// WithBaseURL sets the base URL
func WithBaseURL(url string) Option {
        return func(c *Client) { c.baseURL = url }
}

// NewClient creates a new BSE client
// If no options are provided, it uses the default configuration
func NewClient(opts ...Option) *Client {
        c := &Client{
                http:    request.NewClient(request.DefaultConfig()),
                baseURL: BaseURL,
                headers: map[string]string{
                        "User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
                        "Referer":          Referer,
                        "Accept":           "application/json",
                        "X-Requested-With": "XMLHttpRequest",
                },
        }
        for _, opt := range opts {
                opt(c)
        }
        return c
}

func (c *Client) Close() {
        c.http.Close()
}
