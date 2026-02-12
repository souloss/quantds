package eastmoneyhk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/souloss/quantds/request"
)

// QuoteFields defines the fields to retrieve for HK stock quotes
const QuoteFields = "f12,f13,f14,f2,f3,f4,f5,f6,f7,f8,f9,f10,f15,f16,f17,f18,f20,f21,f22,f23,f24,f25,f26,f27,f28,f29,f30,f31,f32,f33,f34,f35,f36,f37,f38,f39,f40"

// QuoteParams represents parameters for real-time quote request
type QuoteParams struct {
	Symbols  []string // List of stock symbols (optional, empty returns all HK stocks)
	PageSize int      // Number of results per page
	PageNo   int      // Page number (0-based)
}

// QuoteResult represents the real-time quote result
type QuoteResult struct {
	Quotes []QuoteData // List of quotes
	Total  int         // Total count
}

// QuoteData represents a single real-time quote
type QuoteData struct {
	Code         string  // Stock code
	Name         string  // Stock name
	Latest       float64 // Latest price
	Open         float64 // Opening price
	High         float64 // Highest price
	Low          float64 // Lowest price
	PreClose     float64 // Previous close
	Change       float64 // Price change
	ChangeRate   float64 // Change rate (%)
	Volume       float64 // Trading volume
	Turnover     float64 // Trading turnover
	Amplitude    float64 // Price amplitude (%)
	TurnoverRate float64 // Turnover rate (%)
	PE           float64 // P/E ratio
	High52Week   float64 // 52-week high
	Low52Week    float64 // 52-week low
	MarketCap    float64 // Market capitalization (in HKD)
}

// GetQuote retrieves real-time quotes for HK stocks
func (c *Client) GetQuote(ctx context.Context, params *QuoteParams) (*QuoteResult, *request.Record, error) {
	if params.PageSize <= 0 {
		params.PageSize = 100
	}

	// Build fs parameter for HK stocks
	// HK stocks market code: m:116
	fs := "m:116,t:23,m:116,t:80" // HK main board + HK growth board

	query := url.Values{}
	query.Set("pn", strconv.Itoa(params.PageNo+1))
	query.Set("pz", strconv.Itoa(params.PageSize))
	query.Set("po", "1")
	query.Set("np", "1")
	query.Set("fltt", "2")
	query.Set("invt", "2")
	query.Set("fid", "f3")
	query.Set("fs", fs)
	query.Set("fields", QuoteFields)

	apiURL := fmt.Sprintf("%s?%s", QuoteAPI, query.Encode())

	req := request.Request{
		Method:  "GET",
		URL:     apiURL,
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

type quoteResponse struct {
	Data struct {
		Total int                      `json:"total"`
		Diff  []map[string]interface{} `json:"diff"`
	} `json:"data"`
}

func parseQuoteResponse(body []byte) (*QuoteResult, error) {
	var resp quoteResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	quotes := make([]QuoteData, 0, len(resp.Data.Diff))
	for _, item := range resp.Data.Diff {
		quote := QuoteData{
			Code:         getString(item, "f12"),
			Name:         getString(item, "f14"),
			Latest:       getFloat(item, "f2"),
			ChangeRate:   getFloat(item, "f3"),
			Change:       getFloat(item, "f4"),
			Volume:       getFloat(item, "f5"),
			Turnover:     getFloat(item, "f6"),
			Amplitude:    getFloat(item, "f7"),
			TurnoverRate: getFloat(item, "f8"),
			PE:           getFloat(item, "f9"),
			High:         getFloat(item, "f15"),
			Low:          getFloat(item, "f16"),
			Open:         getFloat(item, "f17"),
			PreClose:     getFloat(item, "f18"),
			High52Week:   getFloat(item, "f44"),
			Low52Week:    getFloat(item, "f45"),
			MarketCap:    getFloat(item, "f20"),
		}
		quotes = append(quotes, quote)
	}

	return &QuoteResult{Quotes: quotes, Total: resp.Data.Total}, nil
}

// Helper functions
func getString(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case string:
			return val
		case float64:
			return fmt.Sprintf("%.0f", val)
		}
	}
	return ""
}

func getFloat(data map[string]interface{}, key string) float64 {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case string:
			var f float64
			fmt.Sscanf(val, "%f", &f)
			return f
		}
	}
	return 0
}

// GetQuotesBySymbols retrieves quotes for specific symbols
func (c *Client) GetQuotesBySymbols(ctx context.Context, symbols []string) (*QuoteResult, *request.Record, error) {
	if len(symbols) == 0 {
		return &QuoteResult{}, nil, nil
	}

	// Build secids parameter
	var secids []string
	for _, s := range symbols {
		code, ok := ParseHKSymbol(s)
		if ok {
			secids = append(secids, fmt.Sprintf("116.%s", code))
		}
	}

	if len(secids) == 0 {
		return &QuoteResult{}, nil, nil
	}

	// Use detail API for specific symbols
	apiURL := fmt.Sprintf("%s?secids=%s&fields=%s",
		"https://push2.eastmoney.com/api/qt/ulist.np/get",
		strings.Join(secids, ","),
		QuoteFields)

	req := request.Request{
		Method:  "GET",
		URL:     apiURL,
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

// SpotParams is an alias for QuoteParams
type SpotParams = QuoteParams

// SpotResult is an alias for QuoteResult
type SpotResult = QuoteResult

// SpotQuote is an alias for QuoteData
type SpotQuote = QuoteData

// GetSpot is an alias for GetQuote
func (c *Client) GetSpot(ctx context.Context, params *SpotParams) (*SpotResult, *request.Record, error) {
	return c.GetQuote(ctx, params)
}
