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
	EndpointCoinsList = "/coins/list"

	// Query Parameters
	ParamIncludePlatform = "include_platform"
)

// CoinsListRequest represents parameters for coins list request
type CoinsListRequest struct {
	IncludePlatform bool
}

// CoinsListResponse represents the response structure
type CoinsListResponse []CoinInfo

type CoinInfo struct {
	ID        string            `json:"id"`
	Symbol    string            `json:"symbol"`
	Name      string            `json:"name"`
	Platforms map[string]string `json:"platforms,omitempty"`
}

// GetCoinsList gets the list of all supported coins
func (c *Client) GetCoinsList(ctx context.Context, params *CoinsListRequest) (CoinsListResponse, *request.Record, error) {
	u, _ := url.Parse(BaseURL + EndpointCoinsList)
	q := u.Query()

	if params.IncludePlatform {
		q.Add(ParamIncludePlatform, "true")
	}

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

	var result CoinsListResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, err
	}

	return result, record, nil
}
