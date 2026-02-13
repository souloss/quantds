package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/souloss/quantds/request"
)

// API Endpoint Constants
const (
	EndpointSimplePrice = "/simple/price"
	
	// Query Parameters
	ParamIDs                 = "ids"
	ParamVsCurrencies        = "vs_currencies"
	ParamIncludeMarketCap    = "include_market_cap"
	ParamInclude24hrVol      = "include_24hr_vol"
	ParamInclude24hrChange   = "include_24hr_change"
	ParamIncludeLastUpdatedAt = "include_last_updated_at"
	
	// Default Values
	ValueTrue = "true"
)

// SimplePriceRequest represents parameters for simple price request
type SimplePriceRequest struct {
	IDs                 []string // Coin IDs (e.g. "bitcoin", "ethereum")
	VsCurrencies        []string // Target currencies (e.g. "usd", "cny")
	IncludeMarketCap    bool
	Include24hrVol      bool
	Include24hrChange   bool
	IncludeLastUpdatedAt bool
}

// SimplePriceResponse represents the response structure
// Map key is coin ID, value is map of currency/metric to value
type SimplePriceResponse map[string]map[string]float64

// GetSimplePrice gets the current price of any cryptocurrencies in any other supported currencies that you need.
func (c *Client) GetSimplePrice(ctx context.Context, params *SimplePriceRequest) (SimplePriceResponse, *request.Record, error) {
	if len(params.IDs) == 0 || len(params.VsCurrencies) == 0 {
		return nil, nil, fmt.Errorf("ids and vs_currencies are required")
	}

	u, _ := url.Parse(BaseURL + EndpointSimplePrice)
	q := u.Query()
	q.Add(ParamIDs, strings.Join(params.IDs, ","))
	q.Add(ParamVsCurrencies, strings.Join(params.VsCurrencies, ","))
	
	if params.IncludeMarketCap {
		q.Add(ParamIncludeMarketCap, ValueTrue)
	}
	if params.Include24hrVol {
		q.Add(ParamInclude24hrVol, ValueTrue)
	}
	if params.Include24hrChange {
		q.Add(ParamInclude24hrChange, ValueTrue)
	}
	if params.IncludeLastUpdatedAt {
		q.Add(ParamIncludeLastUpdatedAt, ValueTrue)
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

	var result SimplePriceResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, err
	}

	return result, record, nil
}
