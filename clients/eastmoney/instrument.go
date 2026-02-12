package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/souloss/quantds/request"
)

const InstrumentAPI = "https://push2.eastmoney.com/api/qt/clist/get"

// InstrumentParams represents parameters for securities list request
type InstrumentParams struct {
	Market   string // Market filter: SH, SZ, BJ
	PageSize int    // Page size
	PageNo   int    // Page number (0-based)
}

// InstrumentResult represents the securities list result
type InstrumentResult struct {
	Data  []InstrumentData
	Total int
}

// InstrumentData represents a single security information
type InstrumentData struct {
	Code     string
	Name     string
	Exchange string
	MarketID int
	Industry string
	ListDate string
	Status   string
	Market   string
}

// StockListParams is an alias for InstrumentParams
type StockListParams = InstrumentParams

// StockListResult is an alias for InstrumentResult
type StockListResult = InstrumentResult

// GetInstruments retrieves list of securities
func (c *Client) GetInstruments(ctx context.Context, params *InstrumentParams) (*InstrumentResult, *request.Record, error) {
	if params.PageSize <= 0 {
		params.PageSize = 5000
	}

	// Build fs parameter
	var fs string
	switch params.Market {
	case "SH":
		fs = "m:1+t:2,m:1+t:23"
	case "SZ":
		fs = "m:0+t:6,m:0+t:80"
	case "BJ":
		fs = "m:0+t:81+s:2048"
	default:
		fs = "m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23"
	}

	query := url.Values{}
	query.Set("pn", strconv.Itoa(params.PageNo))
	query.Set("pz", strconv.Itoa(params.PageSize))
	query.Set("po", "1")
	query.Set("np", "1")
	query.Set("fltt", "2")
	query.Set("invt", "2")
	query.Set("fid", "f12")
	query.Set("fs", fs)
	query.Set("fields", "f12,f13,f14,f20,f21")

	apiURL := fmt.Sprintf("%s?%s", InstrumentAPI, query.Encode())

	req := request.Request{
		Method: "GET",
		URL:    apiURL,
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":    "https://quote.eastmoney.com/",
		},
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

// GetStockList is an alias for GetInstruments
func (c *Client) GetStockList(ctx context.Context, params *StockListParams) (*StockListResult, *request.Record, error) {
	return c.GetInstruments(ctx, params)
}

func parseInstrumentResponse(body []byte) (*InstrumentResult, error) {
	var resp quoteResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	instruments := make([]InstrumentData, 0, len(resp.Data.Diff))
	for _, item := range resp.Data.Diff {
		inst := InstrumentData{
			Code:     getString(item, "f12"),
			Name:     getString(item, "f14"),
			MarketID: getInt(item, "f13"),
		}
		if inst.MarketID == 1 {
			inst.Exchange = "SH"
		} else {
			inst.Exchange = "SZ"
		}
		instruments = append(instruments, inst)
	}

	return &InstrumentResult{Data: instruments, Total: resp.Data.Total}, nil
}
