package eodhd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/souloss/quantds/request"
)

type RealTimeParams struct {
	Symbol string // e.g., "AAPL.US"
}

type RealTimeResult struct {
	Symbol        string
	Open          float64
	High          float64
	Low           float64
	Close         float64
	Volume        float64
	PreviousClose float64
	Change        float64
	ChangePercent float64
	Timestamp     int64
}

func (c *Client) GetRealTimeQuote(ctx context.Context, params *RealTimeParams) (*RealTimeResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s/%s?fmt=json&api_token=%s",
		BaseURL, RealTimeAPI, params.Symbol, c.apiKey)

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

	result, err := parseRealTimeResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type realTimeResponse struct {
	Code          string      `json:"code"`
	Timestamp     int64       `json:"timestamp"`
	Open          interface{} `json:"open"`
	High          interface{} `json:"high"`
	Low           interface{} `json:"low"`
	Close         interface{} `json:"close"`
	Volume        interface{} `json:"volume"`
	PreviousClose interface{} `json:"previousClose"`
	Change        interface{} `json:"change"`
	ChangeP       interface{} `json:"change_p"`
}

func parseRealTimeResponse(body []byte) (*RealTimeResult, error) {
	var resp realTimeResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &RealTimeResult{
		Symbol:        resp.Code,
		Open:          toFloat64(resp.Open),
		High:          toFloat64(resp.High),
		Low:           toFloat64(resp.Low),
		Close:         toFloat64(resp.Close),
		Volume:        toFloat64(resp.Volume),
		PreviousClose: toFloat64(resp.PreviousClose),
		Change:        toFloat64(resp.Change),
		ChangePercent: toFloat64(resp.ChangeP),
		Timestamp:     resp.Timestamp,
	}, nil
}

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}
