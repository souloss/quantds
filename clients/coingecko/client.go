package coingecko

import (
	"github.com/souloss/quantds/request"
)

const (
	BaseURL = "https://api.coingecko.com/api/v3"
)

type Client struct {
	http request.Client
}

type Option func(*Client)

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

func (c *Client) Close() {
	if c.http != nil {
		c.http.Close()
	}
}
