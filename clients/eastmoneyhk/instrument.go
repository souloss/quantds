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

// InstrumentFields defines the fields to retrieve for HK stock instruments
const InstrumentFields = "f12,f13,f14,f15,f16,f17,f18,f20,f21,f22,f23,f24,f25,f26,f27,f28,f29,f30,f31,f32,f33,f34,f35,f36,f37,f38,f39,f40"

// InstrumentParams represents parameters for instrument list request
type InstrumentParams struct {
	Exchange   string // Exchange filter (not used for HK)
	AssetType  string // Asset type filter
	PageSize   int    // Number of results per page
	PageNumber int    // Page number (0-based)
}

// InstrumentResult represents the instrument list result
type InstrumentResult struct {
	Instruments []InstrumentData
	Total       int
	PageSize    int
	PageNumber  int
}

// InstrumentData represents a single instrument
type InstrumentData struct {
	Code        string  // Stock code (e.g., "00700")
	Symbol      string  // Full symbol (e.g., "00700.HK")
	Name        string  // Stock name
	LatestPrice float64 // Latest price
	Change      float64 // Price change
	ChangeRate  float64 // Change rate (%)
	Volume      float64 // Trading volume
	Turnover    float64 // Trading turnover
	High        float64 // High price
	Low         float64 // Low price
	Open        float64 // Open price
	PreClose    float64 // Previous close
	MarketCap   float64 // Market capitalization
	PE          float64 // P/E ratio
	ListDate    string  // Listing date
	Status      string  // Trading status
}

// GetInstruments retrieves HK stock instruments
// Uses EastMoney HK stock list API
func (c *Client) GetInstruments(ctx context.Context, params *InstrumentParams) (*InstrumentResult, *request.Record, error) {
	if params.PageSize <= 0 {
		params.PageSize = 100
	}

	// Build fs parameter for HK stocks
	// HK stocks market code: m:116 (main board + growth board)
	fs := "m:116,t:23,m:116,t:80" // HK main board + HK growth board

	query := url.Values{}
	query.Set("pn", strconv.Itoa(params.PageNumber+1))
	query.Set("pz", strconv.Itoa(params.PageSize))
	query.Set("po", "1")   // Sort by change rate descending
	query.Set("np", "1")   // Need pagination
	query.Set("fltt", "2") // Float type
	query.Set("invt", "2") // Inverted sorting
	query.Set("fid", "f3") // Sort by change rate
	query.Set("fs", fs)
	query.Set("fields", InstrumentFields)

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

	result, err := parseInstrumentResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type instrumentResponse struct {
	Data struct {
		Total int                      `json:"total"`
		Diff  []map[string]interface{} `json:"diff"`
	} `json:"data"`
}

func parseInstrumentResponse(body []byte) (*InstrumentResult, error) {
	var resp instrumentResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	instruments := make([]InstrumentData, 0, len(resp.Data.Diff))
	for _, item := range resp.Data.Diff {
		code := getString(item, "f12")
		if code == "" {
			continue
		}

		// Get listing date from f17 (上市日期)
		listDate := getString(item, "f17")
		if listDate != "" {
			// Format: YYYYMMDD -> YYYY-MM-DD
			if len(listDate) == 8 {
				listDate = listDate[0:4] + "-" + listDate[4:6] + "-" + listDate[6:8]
			}
		}

		quote := InstrumentData{
			Code:        code,
			Symbol:      code + ".HK",
			Name:        getString(item, "f14"),
			LatestPrice: getFloat(item, "f2"),
			ChangeRate:  getFloat(item, "f3"),
			Change:      getFloat(item, "f4"),
			Volume:      getFloat(item, "f5"),
			Turnover:    getFloat(item, "f6"),
			PE:          getFloat(item, "f9"),
			High:        getFloat(item, "f15"),
			Low:         getFloat(item, "f16"),
			Open:        getFloat(item, "f17"),
			PreClose:    getFloat(item, "f18"),
			MarketCap:   getFloat(item, "f20"),
			ListDate:    listDate,
		}
		instruments = append(instruments, quote)
	}

	return &InstrumentResult{
		Instruments: instruments,
		Total:       resp.Data.Total,
		PageSize:    len(instruments),
		PageNumber:  0,
	}, nil
}

// GetInstrumentsByCode retrieves instruments by specific codes
func (c *Client) GetInstrumentsByCode(ctx context.Context, codes []string) (*InstrumentResult, *request.Record, error) {
	if len(codes) == 0 {
		return &InstrumentResult{}, nil, nil
	}

	// Build secids parameter
	var secids []string
	for _, code := range codes {
		// Parse code from symbol like "00700.HK"
		symbolCode, ok := ParseHKSymbol(code)
		if ok {
			secids = append(secids, fmt.Sprintf("116.%s", symbolCode))
		}
	}

	if len(secids) == 0 {
		return &InstrumentResult{}, nil, nil
	}

	apiURL := fmt.Sprintf("%s?secids=%s&fields=%s",
		QuoteAPI,
		strings.Join(secids, ","),
		InstrumentFields)

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

	result, err := parseInstrumentResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

// GetAllHKStocks retrieves all HK stocks
func (c *Client) GetAllHKStocks(ctx context.Context) (*InstrumentResult, *request.Record, error) {
	return c.GetInstruments(ctx, &InstrumentParams{
		PageSize:   5000,
		PageNumber: 0,
	})
}
