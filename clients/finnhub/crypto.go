package finnhub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

func (c *Client) GetCryptoCandles(ctx context.Context, params *CandleParams) (*CandleResult, *request.Record, error) {
	return c.getCandles(ctx, CryptoCandleAPI, params)
}

type CryptoSymbolData struct {
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
}

type CryptoSymbolResult struct {
	Symbols []CryptoSymbolData
	Count   int
}

func (c *Client) GetCryptoSymbols(ctx context.Context, exchange string) (*CryptoSymbolResult, *request.Record, error) {
	if exchange == "" {
		exchange = "binance"
	}

	url := fmt.Sprintf("%s%s?exchange=%s&token=%s", BaseURL, CryptoSymbolAPI, exchange, c.apiKey)

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

	var symbols []CryptoSymbolData
	if err := json.Unmarshal(resp.Body, &symbols); err != nil {
		return nil, record, err
	}

	return &CryptoSymbolResult{
		Symbols: symbols,
		Count:   len(symbols),
	}, record, nil
}
