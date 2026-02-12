package bse

import (
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

func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

func NewClient(httpClient request.Client, opts ...Option) *Client {
	if httpClient == nil {
		httpClient = request.NewClient(request.DefaultConfig())
	}
	c := &Client{
		http:    httpClient,
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
