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
	EndpointMarketChart = "/coins/%s/market_chart"
	
	// Query Parameters
	ParamVsCurrency = "vs_currency"
	ParamDays       = "days"
	ParamInterval   = "interval"
)

// MarketChartRequest represents parameters for market chart request
type MarketChartRequest struct {
	ID          string // Coin ID (e.g. "bitcoin")
	VsCurrency  string // Target currency (e.g. "usd")
	Days        string // Data up to number of days ago (e.g. "1", "14", "30", "max")
	Interval    string // Data interval (e.g. "daily")
}

// MarketChartResponse represents the response structure
type MarketChartResponse struct {
	Prices       [][]float64 `json:"prices"`
	MarketCaps   [][]float64 `json:"market_caps"`
	TotalVolumes [][]float64 `json:"total_volumes"`
}

// GetMarketChart gets historical market data include price, market cap, and 24h volume
func (c *Client) GetMarketChart(ctx context.Context, params *MarketChartRequest) (*MarketChartResponse, *request.Record, error) {
	if params.ID == "" || params.VsCurrency == "" || params.Days == "" {
		return nil, nil, fmt.Errorf("id, vs_currency and days are required")
	}

	u, _ := url.Parse(fmt.Sprintf(BaseURL+EndpointMarketChart, params.ID))
	q := u.Query()
	q.Add(ParamVsCurrency, params.VsCurrency)
	q.Add(ParamDays, params.Days)
	if params.Interval != "" {
		q.Add(ParamInterval, params.Interval)
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

	var result MarketChartResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, err
	}

	return &result, record, nil
}
