package twelvedata

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/souloss/quantds/request"
)

type TimeSeriesParams struct {
	Symbol     string
	Interval   string
	OutputSize int
	StartDate  string
	EndDate    string
}

type TimeSeriesResult struct {
	Symbol   string
	Interval string
	Data     []TimeSeriesData
	Count    int
}

type TimeSeriesData struct {
	Datetime string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   float64
}

func (c *Client) GetTimeSeries(ctx context.Context, params *TimeSeriesParams) (*TimeSeriesResult, *request.Record, error) {
	interval := params.Interval
	if interval == "" {
		interval = Interval1day
	}
	outputSize := params.OutputSize
	if outputSize <= 0 {
		outputSize = 30
	}

	url := fmt.Sprintf("%s%s?symbol=%s&interval=%s&outputsize=%d&apikey=%s",
		BaseURL, TimeSeriesAPI, params.Symbol, interval, outputSize, c.apiKey)

	if params.StartDate != "" {
		url += "&start_date=" + params.StartDate
	}
	if params.EndDate != "" {
		url += "&end_date=" + params.EndDate
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

	result, err := parseTimeSeriesResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type twelvedataTimeSeriesResponse struct {
	Meta struct {
		Symbol   string `json:"symbol"`
		Interval string `json:"interval"`
	} `json:"meta"`
	Values []struct {
		Datetime string `json:"datetime"`
		Open     string `json:"open"`
		High     string `json:"high"`
		Low      string `json:"low"`
		Close    string `json:"close"`
		Volume   string `json:"volume"`
	} `json:"values"`
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func parseTimeSeriesResponse(body []byte) (*TimeSeriesResult, error) {
	var resp twelvedataTimeSeriesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Status == "error" {
		return nil, fmt.Errorf("API error (code %d)", resp.Code)
	}

	data := make([]TimeSeriesData, 0, len(resp.Values))
	for _, v := range resp.Values {
		open, _ := strconv.ParseFloat(v.Open, 64)
		high, _ := strconv.ParseFloat(v.High, 64)
		low, _ := strconv.ParseFloat(v.Low, 64)
		close_, _ := strconv.ParseFloat(v.Close, 64)
		volume, _ := strconv.ParseFloat(v.Volume, 64)
		data = append(data, TimeSeriesData{
			Datetime: v.Datetime, Open: open, High: high, Low: low, Close: close_, Volume: volume,
		})
	}

	return &TimeSeriesResult{
		Symbol:   resp.Meta.Symbol,
		Interval: resp.Meta.Interval,
		Data:     data,
		Count:    len(data),
	}, nil
}
