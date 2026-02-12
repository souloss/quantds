package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/souloss/quantds/request"
)

// TickerParams represents parameters for 24hr ticker request
type TickerParams struct {
	Symbol string // Trading pair symbol (optional, empty returns all symbols)
}

// TickerResult represents the 24hr ticker result
type TickerResult struct {
	Tickers []TickerData // List of tickers
	Count   int          // Number of tickers
}

// TickerData represents a single 24hr ticker
type TickerData struct {
	Symbol             string    // Trading pair symbol
	PriceChange        float64   // Price change
	PriceChangePercent float64   // Price change percent
	WeightedAvgPrice   float64   // Weighted average price
	PrevClosePrice     float64   // Previous close price
	LastPrice          float64   // Last price
	LastQty            float64   // Last quantity
	BidPrice           float64   // Best bid price
	BidQty             float64   // Best bid quantity
	AskPrice           float64   // Best ask price
	AskQty             float64   // Best ask quantity
	OpenPrice          float64   // Open price
	HighPrice          float64   // High price
	LowPrice           float64   // Low price
	Volume             float64   // Trading volume (base asset)
	QuoteVolume        float64   // Trading volume (quote asset)
	OpenTime           time.Time // Open time
	CloseTime          time.Time // Close time
	Trades             int       // Number of trades
}

// GetTicker24hr retrieves 24hr ticker data
func (c *Client) GetTicker24hr(ctx context.Context, params *TickerParams) (*TickerResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s", BaseURL, Ticker24hrAPI)
	if params.Symbol != "" {
		url += fmt.Sprintf("?symbol=%s", params.Symbol)
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

	result, err := parseTickerResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

// binanceTickerResponse represents a single ticker from Binance
type binanceTickerResponse struct {
	Symbol             string  `json:"symbol"`
	PriceChange        float64 `json:"priceChange,string"`
	PriceChangePercent float64 `json:"priceChangePercent,string"`
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,string"`
	PrevClosePrice     float64 `json:"prevClosePrice,string"`
	LastPrice          float64 `json:"lastPrice,string"`
	LastQty            float64 `json:"lastQty,string"`
	BidPrice           float64 `json:"bidPrice,string"`
	BidQty             float64 `json:"bidQty,string"`
	AskPrice           float64 `json:"askPrice,string"`
	AskQty             float64 `json:"askQty,string"`
	OpenPrice          float64 `json:"openPrice,string"`
	HighPrice          float64 `json:"highPrice,string"`
	LowPrice           float64 `json:"lowPrice,string"`
	Volume             float64 `json:"volume,string"`
	QuoteVolume        float64 `json:"quoteVolume,string"`
	OpenTime           int64   `json:"openTime"`
	CloseTime          int64   `json:"closeTime"`
	Trades             int     `json:"count"`
}

func parseTickerResponse(body []byte) (*TickerResult, error) {
	// Try single ticker first
	var single binanceTickerResponse
	if err := json.Unmarshal(body, &single); err == nil && single.Symbol != "" {
		return &TickerResult{
			Tickers: []TickerData{tickerFromBinance(single)},
			Count:   1,
		}, nil
	}

	// Try array of tickers
	var arr []binanceTickerResponse
	if err := json.Unmarshal(body, &arr); err != nil {
		return nil, err
	}

	tickers := make([]TickerData, 0, len(arr))
	for _, t := range arr {
		tickers = append(tickers, tickerFromBinance(t))
	}

	return &TickerResult{
		Tickers: tickers,
		Count:   len(tickers),
	}, nil
}

func tickerFromBinance(t binanceTickerResponse) TickerData {
	return TickerData{
		Symbol:             t.Symbol,
		PriceChange:        t.PriceChange,
		PriceChangePercent: t.PriceChangePercent,
		WeightedAvgPrice:   t.WeightedAvgPrice,
		PrevClosePrice:     t.PrevClosePrice,
		LastPrice:          t.LastPrice,
		LastQty:            t.LastQty,
		BidPrice:           t.BidPrice,
		BidQty:             t.BidQty,
		AskPrice:           t.AskPrice,
		AskQty:             t.AskQty,
		OpenPrice:          t.OpenPrice,
		HighPrice:          t.HighPrice,
		LowPrice:           t.LowPrice,
		Volume:             t.Volume,
		QuoteVolume:        t.QuoteVolume,
		OpenTime:           time.UnixMilli(t.OpenTime),
		CloseTime:          time.UnixMilli(t.CloseTime),
		Trades:             t.Trades,
	}
}

// PriceParams represents parameters for price request
type PriceParams struct {
	Symbol string // Trading pair symbol (optional)
}

// PriceResult represents the price result
type PriceResult struct {
	Prices []PriceData
	Count  int
}

// PriceData represents a single price
type PriceData struct {
	Symbol string  // Trading pair symbol
	Price  float64 // Current price
}

// GetPrice retrieves current price for one or all symbols
func (c *Client) GetPrice(ctx context.Context, params *PriceParams) (*PriceResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s", BaseURL, TickerPriceAPI)
	if params.Symbol != "" {
		url += fmt.Sprintf("?symbol=%s", params.Symbol)
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

	result, err := parsePriceResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type binancePriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func parsePriceResponse(body []byte) (*PriceResult, error) {
	// Try single price
	var single binancePriceResponse
	if err := json.Unmarshal(body, &single); err == nil && single.Symbol != "" {
		price := parseFloatStr(single.Price)
		return &PriceResult{
			Prices: []PriceData{{Symbol: single.Symbol, Price: price}},
			Count:  1,
		}, nil
	}

	// Try array
	var arr []binancePriceResponse
	if err := json.Unmarshal(body, &arr); err != nil {
		return nil, err
	}

	prices := make([]PriceData, 0, len(arr))
	for _, p := range arr {
		prices = append(prices, PriceData{
			Symbol: p.Symbol,
			Price:  parseFloatStr(p.Price),
		})
	}

	return &PriceResult{
		Prices: prices,
		Count:  len(prices),
	}, nil
}

func parseFloatStr(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// SpotParams is an alias for TickerParams
type SpotParams = TickerParams

// SpotResult is an alias for TickerResult
type SpotResult = TickerResult

// SpotData is an alias for TickerData
type SpotData = TickerData

// GetSpot is an alias for GetTicker24hr
func (c *Client) GetSpot(ctx context.Context, params *SpotParams) (*SpotResult, *request.Record, error) {
	return c.GetTicker24hr(ctx, params)
}
