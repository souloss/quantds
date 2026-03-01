package finnhub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

type SymbolParams struct {
	Exchange string
}

type SymbolResult struct {
	Symbols []SymbolData
	Count   int
}

type SymbolData struct {
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Currency    string `json:"currency"`
	Type        string `json:"type"`
	Exchange    string `json:"mic"`
}

func (c *Client) GetStockSymbols(ctx context.Context, params *SymbolParams) (*SymbolResult, *request.Record, error) {
	exchange := params.Exchange
	if exchange == "" {
		exchange = "US"
	}

	url := fmt.Sprintf("%s%s?exchange=%s&token=%s", BaseURL, StockSymbolAPI, exchange, c.apiKey)

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

	result, err := parseSymbolResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

func parseSymbolResponse(body []byte) (*SymbolResult, error) {
	var symbols []SymbolData
	if err := json.Unmarshal(body, &symbols); err != nil {
		return nil, err
	}

	return &SymbolResult{
		Symbols: symbols,
		Count:   len(symbols),
	}, nil
}
