package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/souloss/quantds/request"
)

// InstrumentParams represents parameters for instrument list request
type InstrumentParams struct {
	Symbol     string // Trading pair symbol (e.g., "BTCUSDT")
	Asset      string // Base asset (e.g., "BTC")
	Quote      string // Quote asset (e.g., "USDT")
	Status     string // Trading status (TRADING)
	PageSize   int
	PageNumber int
}

// InstrumentResult represents the instrument list result
type InstrumentResult struct {
	Instruments []InstrumentData
	Total       int
}

// InstrumentData represents a single instrument/trading pair
type InstrumentData struct {
	Symbol             string // Trading pair symbol (e.g., "BTCUSDT")
	BaseAsset          string // Base asset (e.g., "BTC")
	QuoteAsset         string // Quote asset (e.g., "USDT")
	BaseAssetPrecision int    // Base asset precision
	QuotePrecision     int    // Quote asset precision
	Status             string // Trading status
	Market             string // Market category (SPOT, MARGIN, etc.)
}

// GetExchangeInfo retrieves exchange information including all trading pairs
func (c *Client) GetExchangeInfo(ctx context.Context) (*InstrumentResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s", BaseURL, ExchangeInfoAPI)

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

	result, err := parseExchangeInfoResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

// GetInstruments retrieves trading pairs with optional filters
func (c *Client) GetInstruments(ctx context.Context, params *InstrumentParams) (*InstrumentResult, *request.Record, error) {
	result, record, err := c.GetExchangeInfo(ctx)
	if err != nil {
		return nil, record, err
	}

	// Apply filters
	instruments := result.Instruments
	if params != nil {
		filtered := make([]InstrumentData, 0)
		for _, inst := range instruments {
			// Filter by symbol
			if params.Symbol != "" && !strings.Contains(inst.Symbol, strings.ToUpper(params.Symbol)) {
				continue
			}
			// Filter by base asset
			if params.Asset != "" && inst.BaseAsset != strings.ToUpper(params.Asset) {
				continue
			}
			// Filter by quote asset
			if params.Quote != "" && inst.QuoteAsset != strings.ToUpper(params.Quote) {
				continue
			}
			// Filter by status
			if params.Status != "" && inst.Status != params.Status {
				continue
			}
			filtered = append(filtered, inst)
		}
		instruments = filtered
	}

	return &InstrumentResult{
		Instruments: instruments,
		Total:       len(instruments),
	}, record, nil
}

// GetAllSpotInstruments retrieves all spot trading pairs
func (c *Client) GetAllSpotInstruments(ctx context.Context) (*InstrumentResult, *request.Record, error) {
	return c.GetInstruments(ctx, &InstrumentParams{
		Status: "TRADING",
	})
}

// GetInstrumentsByQuote retrieves all instruments with a specific quote asset
func (c *Client) GetInstrumentsByQuote(ctx context.Context, quote string) (*InstrumentResult, *request.Record, error) {
	return c.GetInstruments(ctx, &InstrumentParams{
		Quote:  strings.ToUpper(quote),
		Status: "TRADING",
	})
}

type exchangeInfoResponse struct {
	Timezone   string            `json:"timezone"`
	ServerTime int64             `json:"serverTime"`
	Symbols    []json.RawMessage `json:"symbols"`
}

type symbolInfo struct {
	Symbol              string `json:"symbol"`
	BaseAsset           string `json:"baseAsset"`
	BaseAssetPrecision  int    `json:"baseAssetPrecision"`
	QuoteAsset          string `json:"quoteAsset"`
	QuotePrecision      int    `json:"quotePrecision"`
	QuoteAssetPrecision int    `json:"quoteAssetPrecision"`
	Status              string `json:"status"`
	Market              string `json:"market"` // Added in newer API versions
}

func parseExchangeInfoResponse(body []byte) (*InstrumentResult, error) {
	var resp exchangeInfoResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	instruments := make([]InstrumentData, 0, len(resp.Symbols))
	for _, sym := range resp.Symbols {
		var info symbolInfo
		if err := json.Unmarshal(sym, &info); err != nil {
			continue
		}

		// Only include trading symbols
		if info.Status != "TRADING" {
			continue
		}

		instruments = append(instruments, InstrumentData{
			Symbol:             info.Symbol,
			BaseAsset:          info.BaseAsset,
			QuoteAsset:         info.QuoteAsset,
			BaseAssetPrecision: info.BaseAssetPrecision,
			QuotePrecision:     info.QuotePrecision,
			Status:             info.Status,
			Market:             "SPOT",
		})
	}

	return &InstrumentResult{
		Instruments: instruments,
		Total:       len(instruments),
	}, nil
}
