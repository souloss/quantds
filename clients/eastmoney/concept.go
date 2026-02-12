package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/souloss/quantds/request"
)

const (
	ConceptListAPI = "/api/qt/clist/get"
)

// ConceptListParams represents parameters for concept list request
type ConceptListParams struct {
	PageSize int
	PageNo   int
}

// ConceptItem represents a concept/sector
type ConceptItem struct {
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	ChangePercent  float64 `json:"change_percent"`
	MainNetInflow  float64 `json:"main_net_inflow"`
	TotalMarketCap float64 `json:"total_market_cap"`
}

// ConceptStocksParams represents parameters for stocks in a concept
type ConceptStocksParams struct {
	ConceptCode string
	PageSize    int
	PageNo      int
}

// GetConceptList retrieves list of concepts (themes)
func (c *Client) GetConceptList(ctx context.Context, params *ConceptListParams) ([]ConceptItem, *request.Record, error) {
	if params.PageSize <= 0 {
		params.PageSize = 100
	}
	if params.PageNo <= 0 {
		params.PageNo = 1
	}

	query := url.Values{}
	query.Set("pn", strconv.Itoa(params.PageNo))
	query.Set("pz", strconv.Itoa(params.PageSize))
	query.Set("po", "1")
	query.Set("np", "1")
	query.Set("fltt", "2")
	query.Set("invt", "2")
	query.Set("fid", "f3")
	query.Set("fs", "m:90+t:2+f:!50") // Concept board
	query.Set("fields", "f12,f14,f3,f62,f20") // Code, Name, Change%, NetInflow, MarketCap

	url := fmt.Sprintf("%s%s?%s", PushURL, ConceptListAPI, query.Encode())

	req := request.Request{
		Method: "GET",
		URL:    url,
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

	items, err := parseConceptListResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}
	return items, record, nil
}

func parseConceptListResponse(body []byte) ([]ConceptItem, error) {
	var raw struct {
		Data *struct {
			Diff []map[string]interface{} `json:"diff"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	if raw.Data == nil {
		return nil, nil // Empty result
	}

	items := make([]ConceptItem, 0, len(raw.Data.Diff))
	for _, d := range raw.Data.Diff {
		items = append(items, ConceptItem{
			Code:           getString(d, "f12"),
			Name:           getString(d, "f14"),
			ChangePercent:  getFloat(d, "f3"),
			MainNetInflow:  getFloat(d, "f62"),
			TotalMarketCap: getFloat(d, "f20"),
		})
	}

	return items, nil
}

// GetConceptStocks retrieves stocks in a concept
func (c *Client) GetConceptStocks(ctx context.Context, params *ConceptStocksParams) ([]QuoteData, *request.Record, error) {
	if params.ConceptCode == "" {
		return nil, nil, fmt.Errorf("concept code required")
	}
	if params.PageSize <= 0 {
		params.PageSize = 100
	}
	if params.PageNo <= 0 {
		params.PageNo = 1
	}

	query := url.Values{}
	query.Set("pn", strconv.Itoa(params.PageNo))
	query.Set("pz", strconv.Itoa(params.PageSize))
	query.Set("po", "1")
	query.Set("np", "1")
	query.Set("fltt", "2")
	query.Set("invt", "2")
	query.Set("fid", "f3")
	query.Set("fs", fmt.Sprintf("b:%s", params.ConceptCode))
	query.Set("fields", QuoteFields) // Use QuoteFields from quote.go

	url := fmt.Sprintf("%s%s?%s", PushURL, ConceptListAPI, query.Encode())

	req := request.Request{
		Method: "GET",
		URL:    url,
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

	items, err := parseConceptStocksResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}
	return items, record, nil
}

func parseConceptStocksResponse(body []byte) ([]QuoteData, error) {
	var raw struct {
		Data *struct {
			Diff []map[string]interface{} `json:"diff"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	if raw.Data == nil {
		return nil, nil
	}

	items := make([]QuoteData, 0, len(raw.Data.Diff))
	for _, d := range raw.Data.Diff {
		items = append(items, QuoteData{
			Code:         getString(d, "f12"),
			Name:         getString(d, "f14"),
			Latest:       getFloat(d, "f2"),
			Open:         getFloat(d, "f17"),
			High:         getFloat(d, "f15"),
			Low:          getFloat(d, "f16"),
			PreClose:     getFloat(d, "f18"),
			Change:       getFloat(d, "f4"),
			ChangeRate:   getFloat(d, "f3"),
			Volume:       getFloat(d, "f5"),
			Turnover:     getFloat(d, "f6"),
			Amplitude:    getFloat(d, "f7"),
			TurnoverRate: getFloat(d, "f8"),
			PE:           getFloat(d, "f9"),
			VolumeRatio:  getFloat(d, "f10"),
			BidPrice:     getFloat(d, "f31"),
			BidVolume:    getFloat(d, "f32"),
			AskPrice:     getFloat(d, "f33"),
			AskVolume:    getFloat(d, "f34"),
			MarketID:     getInt(d, "f13"),
		})
	}
	return items, nil
}
