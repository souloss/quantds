package eodhd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

type ExchangeSymbolsParams struct {
	Exchange string // e.g., "US", "LSE", "HKEX", "SHG", "SHE"
}

type ExchangeSymbolsResult struct {
	Symbols []ExchangeSymbol
	Count   int
}

type ExchangeSymbol struct {
	Code     string `json:"Code"`
	Name     string `json:"Name"`
	Country  string `json:"Country"`
	Exchange string `json:"Exchange"`
	Currency string `json:"Currency"`
	Type     string `json:"Type"`
	ISIN     string `json:"Isin"`
}

func (c *Client) GetExchangeSymbolList(ctx context.Context, params *ExchangeSymbolsParams) (*ExchangeSymbolsResult, *request.Record, error) {
	exchange := params.Exchange
	if exchange == "" {
		exchange = "US"
	}

	url := fmt.Sprintf("%s%s/%s?fmt=json&api_token=%s",
		BaseURL, ExchangeSymbolsAPI, exchange, c.apiKey)

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

	var symbols []ExchangeSymbol
	if err := json.Unmarshal(resp.Body, &symbols); err != nil {
		return nil, record, err
	}

	return &ExchangeSymbolsResult{
		Symbols: symbols,
		Count:   len(symbols),
	}, record, nil
}
