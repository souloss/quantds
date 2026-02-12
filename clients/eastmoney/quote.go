package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/souloss/quantds/request"
)

const QuoteAPI = "https://push2.eastmoney.com/api/qt/clist/get"
const QuoteFields = "f12,f13,f14,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f15,f16,f17,f18,f20,f21,f22,f23,f24,f25,f26,f27,f28,f29,f30,f31,f32,f33,f34,f35,f36,f37,f38,f39,f40,f41,f42,f43,f44,f45,f46,f47,f48,f49,f50,f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63,f64,f65,f66,f67,f68,f69,f70,f71,f72,f73,f74,f75,f76,f77,f78,f79,f80,f81,f82,f83,f84,f85,f86,f87,f88,f89,f90,f91,f92,f93,f94,f95,f96,f97,f98,f99,f100"

// QuoteParams represents parameters for real-time quote request
type QuoteParams struct {
	Market   string // Market code: SH, SZ, BJ
	PageSize int    // Number of results per page
	PageNo   int    // Page number (0-based)
}

// QuoteResult represents the real-time quote result
type QuoteResult struct {
	Data  []QuoteData
	Total int
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
	VolumeRatio  float64 // Volume ratio
	BidPrice     float64 // Best bid price
	BidVolume    float64 // Best bid volume
	AskPrice     float64 // Best ask price
	AskVolume    float64 // Best ask volume
	MarketID     int     // Market ID (0=SZ, 1=SH)
}

// SpotParams is an alias for QuoteParams for backward compatibility
type SpotParams = QuoteParams

// SpotResult is an alias for QuoteResult for backward compatibility
type SpotResult = QuoteResult

// SpotQuote is an alias for QuoteData for backward compatibility
type SpotQuote = QuoteData

// GetQuotes retrieves real-time quotes for multiple stocks
func (c *Client) GetQuotes(ctx context.Context, params *QuoteParams) (*QuoteResult, *request.Record, error) {
	if params.PageSize <= 0 {
		params.PageSize = 100
	}

	// Build fs parameter based on market
	var fs string
	switch params.Market {
	case "SH":
		fs = "m:1+t:2,m:1+t:23"
	case "SZ":
		fs = "m:0+t:6,m:0+t:80"
	case "BJ":
		fs = "m:0+t:81+s:2048"
	default:
		// All A-shares
		fs = "m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23"
	}

	query := url.Values{}
	query.Set("pn", strconv.Itoa(params.PageNo))
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

	result, err := parseQuoteResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

// GetSpot is an alias for GetQuotes for backward compatibility
func (c *Client) GetSpot(ctx context.Context, params *SpotParams) (*SpotResult, *request.Record, error) {
	return c.GetQuotes(ctx, params)
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
			VolumeRatio:  getFloat(item, "f10"),
			High:         getFloat(item, "f15"),
			Low:          getFloat(item, "f16"),
			Open:         getFloat(item, "f17"),
			PreClose:     getFloat(item, "f18"),
			MarketID:     getInt(item, "f13"),
		}
		quotes = append(quotes, quote)
	}

	return &QuoteResult{Data: quotes, Total: resp.Data.Total}, nil
}
