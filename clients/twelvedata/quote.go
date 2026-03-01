package twelvedata

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/souloss/quantds/request"
)

type QuoteParams struct {
	Symbol string
}

type QuoteResult struct {
	Symbol        string
	Name          string
	Exchange      string
	Currency      string
	Open          float64
	High          float64
	Low           float64
	Close         float64
	PreviousClose float64
	Volume        float64
	Change        float64
	PercentChange float64
	Datetime      string
}

func (c *Client) GetQuote(ctx context.Context, params *QuoteParams) (*QuoteResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s?symbol=%s&apikey=%s", BaseURL, QuoteAPI, params.Symbol, c.apiKey)

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

	result, err := parseQuoteResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type twelvedataQuoteResponse struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	Exchange      string `json:"exchange"`
	Currency      string `json:"currency"`
	Open          string `json:"open"`
	High          string `json:"high"`
	Low           string `json:"low"`
	Close         string `json:"close"`
	PreviousClose string `json:"previous_close"`
	Volume        string `json:"volume"`
	Change        string `json:"change"`
	PercentChange string `json:"percent_change"`
	Datetime      string `json:"datetime"`
	Status        string `json:"status"`
	Code          int    `json:"code"`
}

func parseQuoteResponse(body []byte) (*QuoteResult, error) {
	var resp twelvedataQuoteResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Status == "error" {
		return nil, fmt.Errorf("API error (code %d)", resp.Code)
	}

	open, _ := strconv.ParseFloat(resp.Open, 64)
	high, _ := strconv.ParseFloat(resp.High, 64)
	low, _ := strconv.ParseFloat(resp.Low, 64)
	close_, _ := strconv.ParseFloat(resp.Close, 64)
	prevClose, _ := strconv.ParseFloat(resp.PreviousClose, 64)
	volume, _ := strconv.ParseFloat(resp.Volume, 64)
	change, _ := strconv.ParseFloat(resp.Change, 64)
	pctChange, _ := strconv.ParseFloat(resp.PercentChange, 64)

	return &QuoteResult{
		Symbol: resp.Symbol, Name: resp.Name, Exchange: resp.Exchange, Currency: resp.Currency,
		Open: open, High: high, Low: low, Close: close_, PreviousClose: prevClose,
		Volume: volume, Change: change, PercentChange: pctChange, Datetime: resp.Datetime,
	}, nil
}
