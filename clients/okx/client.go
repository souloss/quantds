package okx

import (
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

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.BaseURL = url
	}
}

func NewClient(httpClient request.Client, opts ...Option) *Client {
	if httpClient == nil {
		httpClient = request.NewClient(request.DefaultConfig())
	}
	c := &Client{
		http:    httpClient,
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
