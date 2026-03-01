package finnhub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/souloss/quantds/request"
)

type QuoteParams struct {
	Symbol string
}

type QuoteResult struct {
	Symbol        string
	Open          float64
	High          float64
	Low           float64
	Current       float64
	PreviousClose float64
	Change        float64
	PercentChange float64
	Timestamp     time.Time
}

func (c *Client) GetQuote(ctx context.Context, params *QuoteParams) (*QuoteResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s?symbol=%s&token=%s", BaseURL, QuoteAPI, params.Symbol, c.apiKey)

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

	result, err := parseQuoteResponse(resp.Body, params.Symbol)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type finnhubQuoteResponse struct {
	Open          float64 `json:"o"`
	High          float64 `json:"h"`
	Low           float64 `json:"l"`
	Current       float64 `json:"c"`
	PreviousClose float64 `json:"pc"`
	Change        float64 `json:"d"`
	PercentChange float64 `json:"dp"`
	Timestamp     int64   `json:"t"`
}

func parseQuoteResponse(body []byte, symbol string) (*QuoteResult, error) {
	var resp finnhubQuoteResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &QuoteResult{
		Symbol:        symbol,
		Open:          resp.Open,
		High:          resp.High,
		Low:           resp.Low,
		Current:       resp.Current,
		PreviousClose: resp.PreviousClose,
		Change:        resp.Change,
		PercentChange: resp.PercentChange,
		Timestamp:     time.Unix(resp.Timestamp, 0),
	}, nil
}
