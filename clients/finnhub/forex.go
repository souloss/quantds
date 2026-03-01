package finnhub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

type ForexRatesParams struct {
	Base string
}

type ForexRatesResult struct {
	Base  string
	Rates map[string]float64
}

func (c *Client) GetForexRates(ctx context.Context, params *ForexRatesParams) (*ForexRatesResult, *request.Record, error) {
	base := params.Base
	if base == "" {
		base = "USD"
	}

	url := fmt.Sprintf("%s%s?base=%s&token=%s", BaseURL, ForexRatesAPI, base, c.apiKey)

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

	result, err := parseForexRatesResponse(resp.Body, base)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type finnhubForexRatesResponse struct {
	Base  string             `json:"base"`
	Quote map[string]float64 `json:"quote"`
}

func parseForexRatesResponse(body []byte, base string) (*ForexRatesResult, error) {
	var resp finnhubForexRatesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &ForexRatesResult{
		Base:  resp.Base,
		Rates: resp.Quote,
	}, nil
}

func (c *Client) GetForexCandles(ctx context.Context, params *CandleParams) (*CandleResult, *request.Record, error) {
	return c.getCandles(ctx, ForexCandleAPI, params)
}

type ForexSymbolData struct {
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
}

type ForexSymbolResult struct {
	Symbols []ForexSymbolData
	Count   int
}

func (c *Client) GetForexSymbols(ctx context.Context, exchange string) (*ForexSymbolResult, *request.Record, error) {
	if exchange == "" {
		exchange = "oanda"
	}

	url := fmt.Sprintf("%s%s?exchange=%s&token=%s", BaseURL, ForexSymbolAPI, exchange, c.apiKey)

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

	var symbols []ForexSymbolData
	if err := json.Unmarshal(resp.Body, &symbols); err != nil {
		return nil, record, err
	}

	return &ForexSymbolResult{
		Symbols: symbols,
		Count:   len(symbols),
	}, record, nil
}
