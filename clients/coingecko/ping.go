package coingecko

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

// API Endpoint Constants
const (
	EndpointPing = "/ping"
)

// PingResponse represents the response structure
type PingResponse struct {
	GeckoSays string `json:"gecko_says"`
}

// Ping checks API server status
func (c *Client) Ping(ctx context.Context) (*PingResponse, error) {
	req := request.Request{
		Method: "GET",
		URL:    BaseURL + EndpointPing,
	}

	resp, _, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result PingResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
