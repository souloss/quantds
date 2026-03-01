package finnhub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

type CandleParams struct {
	Symbol     string
	Resolution string
	From       int64
	To         int64
}

type CandleResult struct {
	Symbol  string
	Candles []CandleData
	Count   int
}

type CandleData struct {
	Timestamp int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

func (c *Client) GetStockCandles(ctx context.Context, params *CandleParams) (*CandleResult, *request.Record, error) {
	return c.getCandles(ctx, StockCandleAPI, params)
}

func (c *Client) getCandles(ctx context.Context, apiPath string, params *CandleParams) (*CandleResult, *request.Record, error) {
	resolution := params.Resolution
	if resolution == "" {
		resolution = ResD
	}

	url := fmt.Sprintf("%s%s?symbol=%s&resolution=%s&from=%d&to=%d&token=%s",
		BaseURL, apiPath, params.Symbol, resolution, params.From, params.To, c.apiKey)

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

	result, err := parseCandleResponse(resp.Body, params.Symbol)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type finnhubCandleResponse struct {
	Status     string    `json:"s"`
	Close      []float64 `json:"c"`
	High       []float64 `json:"h"`
	Low        []float64 `json:"l"`
	Open       []float64 `json:"o"`
	Volume     []float64 `json:"v"`
	Timestamps []int64   `json:"t"`
}

func parseCandleResponse(body []byte, symbol string) (*CandleResult, error) {
	var resp finnhubCandleResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Status == "no_data" || len(resp.Timestamps) == 0 {
		return &CandleResult{Symbol: symbol}, nil
	}

	candles := make([]CandleData, 0, len(resp.Timestamps))
	for i, ts := range resp.Timestamps {
		candle := CandleData{Timestamp: ts}
		if i < len(resp.Open) {
			candle.Open = resp.Open[i]
		}
		if i < len(resp.High) {
			candle.High = resp.High[i]
		}
		if i < len(resp.Low) {
			candle.Low = resp.Low[i]
		}
		if i < len(resp.Close) {
			candle.Close = resp.Close[i]
		}
		if i < len(resp.Volume) {
			candle.Volume = resp.Volume[i]
		}
		candles = append(candles, candle)
	}

	return &CandleResult{
		Symbol:  symbol,
		Candles: candles,
		Count:   len(candles),
	}, nil
}
