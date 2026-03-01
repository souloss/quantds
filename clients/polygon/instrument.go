package polygon

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

type TickerParams struct {
	Type     string // CS (common stock), ETF, FUND, etc.
	Market   string // stocks, crypto, fx, otc
	Exchange string
	Search   string
	Limit    int
}

type TickerResult struct {
	Tickers []TickerData
	Count   int
}

type TickerData struct {
	Ticker          string `json:"ticker"`
	Name            string `json:"name"`
	Market          string `json:"market"`
	Locale          string `json:"locale"`
	PrimaryExchange string `json:"primary_exchange"`
	Type            string `json:"type"`
	Active          bool   `json:"active"`
	CurrencyName    string `json:"currency_name"`
}

func (c *Client) GetTickers(ctx context.Context, params *TickerParams) (*TickerResult, *request.Record, error) {
	limit := params.Limit
	if limit <= 0 {
		limit = 100
	}

	url := fmt.Sprintf("%s%s?limit=%d&apiKey=%s", BaseURL, TickersAPI, limit, c.apiKey)

	if params.Type != "" {
		url += fmt.Sprintf("&type=%s", params.Type)
	}
	if params.Market != "" {
		url += fmt.Sprintf("&market=%s", params.Market)
	}
	if params.Exchange != "" {
		url += fmt.Sprintf("&exchange=%s", params.Exchange)
	}
	if params.Search != "" {
		url += fmt.Sprintf("&search=%s", params.Search)
	}

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

	result, err := parseTickerResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type polygonTickerResponse struct {
	Results []TickerData `json:"results"`
	Status  string       `json:"status"`
	Count   int          `json:"count"`
}

func parseTickerResponse(body []byte) (*TickerResult, error) {
	var resp polygonTickerResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &TickerResult{
		Tickers: resp.Results,
		Count:   len(resp.Results),
	}, nil
}
