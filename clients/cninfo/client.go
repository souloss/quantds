package cninfo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/souloss/quantds/request"
)

const BaseURL = "http://www.cninfo.com.cn"

type Client struct {
	http    request.Client
	baseURL string
}

type Option func(*Client)

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

func NewClient(httpClient request.Client, opts ...Option) *Client {
	if httpClient == nil {
		httpClient = request.NewClient(request.DefaultConfig())
	}
	c := &Client{
		http:    httpClient,
		baseURL: BaseURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Close() {
	c.http.Close()
}

func (c *Client) do(ctx context.Context, method, path string, form url.Values, result any) (*request.Record, error) {
	u := c.baseURL + path
	var body []byte
	var contentType string

	if form != nil {
		body = []byte(form.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	req := request.Request{
		Method: method,
		URL:    u,
		Headers: map[string]string{
			"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Accept":           "application/json, text/javascript, */*",
			"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8",
			"Referer":          c.baseURL + "/new/disclosure/list/notice",
			"X-Requested-With": "XMLHttpRequest",
		},
	}
	if len(body) > 0 {
		req.Headers["Content-Type"] = contentType
		req.Body = body
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return record, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return record, fmt.Errorf("http status %d", resp.StatusCode)
	}

	if result != nil {
		if err := json.Unmarshal(resp.Body, result); err != nil {
			return record, fmt.Errorf("decode response: %w", err)
		}
	}

	return record, nil
}
