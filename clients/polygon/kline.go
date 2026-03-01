package polygon

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

type AggregateParams struct {
	Symbol     string
	Multiplier int
	Timespan   string
	From       string // YYYY-MM-DD
	To         string // YYYY-MM-DD
	Limit      int
}

type AggregateResult struct {
	Symbol string
	Bars   []AggregateBar
	Count  int
}

type AggregateBar struct {
	Timestamp    int64
	Open         float64
	High         float64
	Low          float64
	Close        float64
	Volume       float64
	VWAP         float64
	Transactions int
}

func (c *Client) GetAggregates(ctx context.Context, params *AggregateParams) (*AggregateResult, *request.Record, error) {
	multiplier := params.Multiplier
	if multiplier <= 0 {
		multiplier = 1
	}
	timespan := params.Timespan
	if timespan == "" {
		timespan = TimespanDay
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 120
	}

	url := fmt.Sprintf("%s%s/%s/range/%d/%s/%s/%s?adjusted=true&sort=asc&limit=%d&apiKey=%s",
		BaseURL, AggregatesAPI, params.Symbol, multiplier, timespan,
		params.From, params.To, limit, c.apiKey)

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

	result, err := parseAggregateResponse(resp.Body, params.Symbol)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type polygonAggResponse struct {
	Ticker       string `json:"ticker"`
	ResultsCount int    `json:"resultsCount"`
	Results      []struct {
		T  int64   `json:"t"`
		O  float64 `json:"o"`
		H  float64 `json:"h"`
		L  float64 `json:"l"`
		C  float64 `json:"c"`
		V  float64 `json:"v"`
		VW float64 `json:"vw"`
		N  int     `json:"n"`
	} `json:"results"`
	Status string `json:"status"`
}

func parseAggregateResponse(body []byte, symbol string) (*AggregateResult, error) {
	var resp polygonAggResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	bars := make([]AggregateBar, 0, len(resp.Results))
	for _, r := range resp.Results {
		bars = append(bars, AggregateBar{
			Timestamp:    r.T,
			Open:         r.O,
			High:         r.H,
			Low:          r.L,
			Close:        r.C,
			Volume:       r.V,
			VWAP:         r.VW,
			Transactions: r.N,
		})
	}

	return &AggregateResult{
		Symbol: symbol,
		Bars:   bars,
		Count:  len(bars),
	}, nil
}
