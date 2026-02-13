package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/souloss/quantds/request"
)

// API Endpoint Constants
const (
	EndpointSearch = "/search"
	
	// Query Parameters
	ParamQuery = "query"
)

// SearchRequest represents parameters for search request
type SearchRequest struct {
	Query string
}

// SearchResponse represents the response structure
type SearchResponse struct {
	Coins []SearchCoin `json:"coins"`
}

type SearchCoin struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	MarketCapRank int    `json:"market_cap_rank"`
	Thumb         string `json:"thumb"`
	Large         string `json:"large"`
}

// Search searches for coins, categories and markets
func (c *Client) Search(ctx context.Context, params *SearchRequest) (*SearchResponse, *request.Record, error) {
	u, _ := url.Parse(BaseURL + EndpointSearch)
	q := u.Query()
	q.Add(ParamQuery, params.Query)

	req := request.Request{
		Method: "GET",
		URL:    u.String() + "?" + q.Encode(),
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result SearchResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, err
	}

	return &result, record, nil
}
