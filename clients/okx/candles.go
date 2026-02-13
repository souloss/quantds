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
	EndpointCandles = "/api/v5/market/candles"

	// Query Parameters
	ParamBar   = "bar"
	ParamLimit = "limit"
)

// CandlestickRequest represents request parameters for candlesticks
type CandlestickRequest struct {
	InstID string
	Bar    string
	Limit  int
}

// CandlestickResponse represents a single K-line bar
// OKX returns: [ts, o, h, l, c, vol, volCcy, volCcyQuote, confirm]
type CandlestickResponse []string

// GetCandlesticks gets historical k-line data
// bar: 1m, 3m, 5m, 15m, 30m, 1H, 2H, 4H, 6H, 12H, 1D, 1W, 1M, 3M
func (c *Client) GetCandlesticks(ctx context.Context, params *CandlestickRequest) ([]CandlestickResponse, *request.Record, error) {
	if params.InstID == "" {
		return nil, nil, fmt.Errorf("instId is required")
	}

	u, _ := url.Parse(c.BaseURL + EndpointCandles)
	q := u.Query()
	q.Add(ParamInstID, params.InstID)
	if params.Bar != "" {
		q.Add(ParamBar, params.Bar)
	}
	if params.Limit > 0 {
		q.Add(ParamLimit, fmt.Sprintf("%d", params.Limit))
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

	var result Response
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, err
	}

	if result.Code != "0" {
		return nil, record, fmt.Errorf("api error: %s (code: %s)", result.Msg, result.Code)
	}

	var candles []CandlestickResponse
	if err := json.Unmarshal(result.Data, &candles); err != nil {
		return nil, record, err
	}

	return candles, record, nil
}
