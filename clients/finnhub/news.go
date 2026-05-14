package finnhub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/souloss/quantds/request"
)

// NewsParams represents parameters for company news request
type NewsParams struct {
	Symbol string // Stock symbol (e.g., "AAPL")
	From   string // Start date (YYYY-MM-DD)
	To     string // End date (YYYY-MM-DD)
}

// NewsResult represents the company news result
type NewsResult struct {
	Symbol string
	News   []NewsItem
}

// NewsItem represents a single news item
type NewsItem struct {
	ID       int64     // News ID
	Datetime time.Time // Publication time
	Headline string    // News headline
	Summary  string    // News summary
	Source   string    // News source
	URL      string    // News URL
	Image    string    // News image URL
	Category string    // News category
}

// GetCompanyNews retrieves company news
func (c *Client) GetCompanyNews(ctx context.Context, params *NewsParams) (*NewsResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s?symbol=%s&from=%s&to=%s&token=%s",
		BaseURL, CompanyNewsAPI, params.Symbol, params.From, params.To, c.apiKey)

	req := request.Request{
		Method:  "GET",
		URL:     url,
		Headers: DefaultHeaders,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseNewsResponse(resp.Body, params.Symbol)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type finnhubNewsResponse []finnhubNewsItem

type finnhubNewsItem struct {
	ID       int64  `json:"id"`
	Datetime int64  `json:"datetime"`
	Headline string `json:"headline"`
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	URL      string `json:"url"`
	Image    string `json:"image"`
	Category string `json:"category"`
}

func parseNewsResponse(body []byte, symbol string) (*NewsResult, error) {
	var resp finnhubNewsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	items := make([]NewsItem, 0, len(resp))
	for _, n := range resp {
		items = append(items, NewsItem{
			ID:       n.ID,
			Datetime: time.Unix(n.Datetime, 0),
			Headline: n.Headline,
			Summary:  n.Summary,
			Source:   n.Source,
			URL:      n.URL,
			Image:    n.Image,
			Category: n.Category,
		})
	}

	return &NewsResult{
		Symbol: symbol,
		News:   items,
	}, nil
}
