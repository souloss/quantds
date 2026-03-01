package twelvedata

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

type ListParams struct {
	Exchange string
	Country  string
	Type     string
}

type ListResult struct {
	Data  []ListItem
	Count int
}

type ListItem struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Exchange string `json:"exchange"`
	Country  string `json:"country"`
	Type     string `json:"type"`
}

func (c *Client) getList(ctx context.Context, apiPath string, params *ListParams) (*ListResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s?apikey=%s", BaseURL, apiPath, c.apiKey)
	if params != nil {
		if params.Exchange != "" {
			url += "&exchange=" + params.Exchange
		}
		if params.Country != "" {
			url += "&country=" + params.Country
		}
		if params.Type != "" {
			url += "&type=" + params.Type
		}
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

	result, err := parseListResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type twelvedataListResponse struct {
	Data   []ListItem `json:"data"`
	Status string     `json:"status"`
}

func parseListResponse(body []byte) (*ListResult, error) {
	var resp twelvedataListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &ListResult{
		Data:  resp.Data,
		Count: len(resp.Data),
	}, nil
}

func (c *Client) GetStocksList(ctx context.Context, params *ListParams) (*ListResult, *request.Record, error) {
	return c.getList(ctx, StocksAPI, params)
}

func (c *Client) GetForexPairsList(ctx context.Context) (*ListResult, *request.Record, error) {
	return c.getList(ctx, ForexPairsAPI, nil)
}

func (c *Client) GetCryptocurrenciesList(ctx context.Context) (*ListResult, *request.Record, error) {
	return c.getList(ctx, CryptoAPI, nil)
}

func (c *Client) GetETFList(ctx context.Context, params *ListParams) (*ListResult, *request.Record, error) {
	return c.getList(ctx, ETFAPI, params)
}

func (c *Client) GetFundsList(ctx context.Context, params *ListParams) (*ListResult, *request.Record, error) {
	return c.getList(ctx, FundsAPI, params)
}

func (c *Client) GetBondsList(ctx context.Context, params *ListParams) (*ListResult, *request.Record, error) {
	return c.getList(ctx, BondsAPI, params)
}
