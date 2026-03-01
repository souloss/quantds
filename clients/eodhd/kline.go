package eodhd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/souloss/quantds/request"
)

type EODParams struct {
	Symbol string // e.g., "AAPL.US", "VOD.LSE"
	From   string // YYYY-MM-DD
	To     string // YYYY-MM-DD
	Period string // d (daily), w (weekly), m (monthly)
}

type EODResult struct {
	Symbol string
	Data   []EODData
	Count  int
}

type EODData struct {
	Date          string
	Open          float64
	High          float64
	Low           float64
	Close         float64
	AdjustedClose float64
	Volume        float64
}

func (c *Client) GetEOD(ctx context.Context, params *EODParams) (*EODResult, *request.Record, error) {
	period := params.Period
	if period == "" {
		period = "d"
	}

	url := fmt.Sprintf("%s%s/%s?fmt=json&period=%s&api_token=%s",
		BaseURL, EODAPI, params.Symbol, period, c.apiKey)

	if params.From != "" {
		url += "&from=" + params.From
	}
	if params.To != "" {
		url += "&to=" + params.To
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

	result, err := parseEODResponse(resp.Body, params.Symbol)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type eodItem struct {
	Date          string      `json:"date"`
	Open          interface{} `json:"open"`
	High          interface{} `json:"high"`
	Low           interface{} `json:"low"`
	Close         interface{} `json:"close"`
	AdjustedClose interface{} `json:"adjusted_close"`
	Volume        interface{} `json:"volume"`
}

func parseEODResponse(body []byte, symbol string) (*EODResult, error) {
	var items []eodItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}

	data := make([]EODData, 0, len(items))
	for _, item := range items {
		data = append(data, EODData{
			Date:          item.Date,
			Open:          toFloat(item.Open),
			High:          toFloat(item.High),
			Low:           toFloat(item.Low),
			Close:         toFloat(item.Close),
			AdjustedClose: toFloat(item.AdjustedClose),
			Volume:        toFloat(item.Volume),
		})
	}

	return &EODResult{
		Symbol: symbol,
		Data:   data,
		Count:  len(data),
	}, nil
}

func toFloat(v interface{}) float64 {
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
