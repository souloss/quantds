// Package cninfo provides a client for the CNInfo (巨潮资讯) data source.
//
// CNInfo is the official information disclosure platform for China's securities market,
// operated by Shenzhen Securities Information Co., Ltd. It provides APIs for:
//   - Stock/Security lists and search
//   - Company announcements and filings
//   - Regulatory disclosures
//   - IPO prospectuses
//
// API Features:
//   - No authentication required for public APIs
//   - Supports A-shares (SH, SZ, BJ), B-shares, and funds
//   - Official regulatory data source
//   - Historical data available
//
// Limitations:
//   - APIs may be rate-limited
//   - Some endpoints may require specific headers
//   - API structure may change over time
//
// Example:
//
//	client := cninfo.NewClient(nil)
//	defer client.Close()
//
//	// Search for securities
//	result, record, err := client.GetOrgID(ctx, &cninfo.OrgIDParams{
//	    KeyWord: "000001",
//	})
//
//	// Get company announcements
//	announcements, total, record, err := client.QueryNews(ctx, &cninfo.NewsQueryParams{
//	    Stock: "000001,gssz0000001",
//	    PageNum: 1,
//	    PageSize: 30,
//	})
package cninfo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/souloss/quantds/request"
)

const (
	BaseURL         = "http://www.cninfo.com.cn"
	SearchAPI       = "http://www.cninfo.com.cn/new/information/topSearch/query"
	AnnouncementAPI = "http://www.cninfo.com.cn/new/hisAnnouncement/query"
)

type Client struct {
	http    request.Client
	baseURL string
}

type Option func(*Client)

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

// do makes an HTTP request to the CNInfo API
func (c *Client) do(ctx context.Context, method, path string, params url.Values, result interface{}) (*request.Record, error) {
	var body []byte
	if params != nil {
		body = []byte(params.Encode())
	}

	urlStr := c.baseURL + path
	// Handle full URLs (if path already includes protocol)
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		urlStr = path
	}

	req := request.Request{
		Method: method,
		URL:    urlStr,
		Headers: map[string]string{
			"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":          "http://www.cninfo.com.cn/",
			"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
			"X-Requested-With": "XMLHttpRequest",
		},
		Body: body,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return record, err
	}

	if resp.StatusCode != 200 {
		return record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Check if the response is JSONP (wrapped in callback)
	bodyContent := string(resp.Body)
	if strings.HasPrefix(bodyContent, "jQuery") || strings.HasPrefix(bodyContent, "callback") {
		start := strings.Index(bodyContent, "(")
		end := strings.LastIndex(bodyContent, ")")
		if start > 0 && end > start {
			bodyContent = bodyContent[start+1 : end]
			resp.Body = []byte(bodyContent)
		}
	}

	// Some responses have metadata wrapper, try to unwrap
	var genericResp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(resp.Body, &genericResp); err == nil {
		if genericResp.Code == 0 && len(genericResp.Data) > 0 {
			// Use the data field
			if err := json.Unmarshal(genericResp.Data, result); err != nil {
				return record, err
			}
			return record, nil
		}
	}

	if err := json.Unmarshal(resp.Body, result); err != nil {
		return record, err
	}

	return record, nil
}
