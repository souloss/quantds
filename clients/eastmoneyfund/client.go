package eastmoneyfund

import (
	"time"

	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/souloss/quantds/request"
)

const (
	FundListURL      = "http://fund.eastmoney.com/js/fundcode_search.js"
	FundEstimateURL  = "https://fundgz.1234567.com.cn/js"
	FundNAVURL       = "http://fund.eastmoney.com/f10/F10DataApi.aspx"
	FundDetailURL    = "http://fund.eastmoney.com/pingzhongdata"
)

var DefaultHeaders = map[string]string{
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	"Referer":    "http://fund.eastmoney.com/",
}

type Client struct {
	http request.Client
}

type Option func(*Client)

func WithHTTPClient(httpClient request.Client) Option {
	return func(c *Client) { c.http = httpClient }
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.http = request.NewClient(request.DefaultConfig(
			request.WithTimeout(timeout.New[request.Response](d)),
		))
	}
}

func WithConfig(cfg *request.Config) Option {
	return func(c *Client) { c.http = request.NewClient(cfg) }
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		http: request.NewClient(request.DefaultConfig()),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Close() { c.http.Close() }
