package szse

import (
	"github.com/souloss/quantds/request"
)

const (
	BaseURL = "https://www.szse.cn/api/report/ShowReport"
	Referer = "https://www.szse.com.cn/market/product/stock/list/index.html"
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
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":    Referer,
			"Accept":     "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
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
