package okx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/souloss/quantds/request"
)

// API Endpoint Constants
const (
	EndpointTicker = "/api/v5/market/ticker"

	// Query Parameters
	ParamInstID = "instId"
)

// TickerRequest represents request parameters for ticker
type TickerRequest struct {
	InstID string
}

// TickerResponse represents OKX ticker data
type TickerResponse struct {
	InstID    string `json:"instId"`
	Last      string `json:"last"`
	LastSz    string `json:"lastSz"`
	AskPx     string `json:"askPx"`
	AskSz     string `json:"askSz"`
	BidPx     string `json:"bidPx"`
	BidSz     string `json:"bidSz"`
	Open24h   string `json:"open24h"`
	High24h   string `json:"high24h"`
	Low24h    string `json:"low24h"`
	VolCcy24h string `json:"volCcy24h"`
	Vol24h    string `json:"vol24h"`
	Ts        string `json:"ts"`
}

// GetTicker gets the latest ticker for a specific instrument
func (c *Client) GetTicker(ctx context.Context, params *TickerRequest) (*TickerResponse, *request.Record, error) {
	if params.InstID == "" {
		return nil, nil, fmt.Errorf("instId is required")
	}

	u, _ := url.Parse(c.BaseURL + EndpointTicker)
	q := u.Query()
	q.Add(ParamInstID, params.InstID)

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

	var result Response
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, err
	}

	if result.Code != "0" {
		return nil, record, fmt.Errorf("api error: %s (code: %s)", result.Msg, result.Code)
	}

	var tickers []TickerResponse
	if err := json.Unmarshal(result.Data, &tickers); err != nil {
		return nil, record, err
	}

	if len(tickers) == 0 {
		return nil, record, fmt.Errorf("no ticker found")
	}

	return &tickers[0], record, nil
}
