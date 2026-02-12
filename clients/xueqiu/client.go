package xueqiu

import (
	"github.com/souloss/quantds/request"
)

const (
	BaseURL = "https://stock.xueqiu.com/v5/stock"
)

type Client struct {
	http   request.Client
	token  string
	cookie string
}

type Option func(*Client)

func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

func WithCookie(cookie string) Option {
	return func(c *Client) {
		c.cookie = cookie
	}
}

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
	c.http.Close()
}
