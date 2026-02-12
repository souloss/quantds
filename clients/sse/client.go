package sse

import (
	"github.com/souloss/quantds/request"
)

const (
	BaseURL = "https://query.sse.com.cn/sseQuery/commonQuery.do"
	Host    = "query.sse.com.cn"
	Referer = "https://www.sse.com.cn/assortment/stock/list/share/"
	SqlID   = "COMMON_SSE_CP_GPJCTPZ_GPLB_GP_L"
	Status  = "2,4,5,7,8"
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
			"Host":         Host,
			"Pragma":       "no-cache",
			"Referer":      Referer,
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Accept":       "application/json",
			"Content-Type": "application/json",
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
