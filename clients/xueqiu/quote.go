package xueqiu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/souloss/quantds/request"
)

const (
	QuoteDetailAPI = "/v5/stock/quote.json"
)

// QuoteDetailResult represents the response for quote detail
type QuoteDetailResult struct {
	Data             QuoteDetailData `json:"data"`
	ErrorCode        int             `json:"error_code"`
	ErrorDescription string          `json:"error_description"`
}

// QuoteDetailData contains the quote and market info
type QuoteDetailData struct {
	Quote  Quote  `json:"quote"`
	Market Market `json:"market"`
}

// Quote represents the detailed quote information
type Quote struct {
	Symbol             string  `json:"symbol"`
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	Exchange           string  `json:"exchange"`
	Currency           string  `json:"currency"`
	Current            float64 `json:"current"`
	Percent            float64 `json:"percent"`
	Chg                float64 `json:"chg"`
	High               float64 `json:"high"`
	Low                float64 `json:"low"`
	Open               float64 `json:"open"`
	LastClose          float64 `json:"last_close"`
	MarketCapital      float64 `json:"market_capital"`
	FloatMarketCapital float64 `json:"float_market_capital"`
	TotalShares        float64 `json:"total_shares"`
	FloatShares        float64 `json:"float_shares"`
	PeTtm              float64 `json:"pe_ttm"`
	PeLyr              float64 `json:"pe_lyr"`
	Pb                 float64 `json:"pb"`
	Eps                float64 `json:"eps"`
	Navps              float64 `json:"navps"`
	DividendYield      float64 `json:"dividend_yield"`
	IssueDate          int64   `json:"issue_date"`
	Timestamp          int64   `json:"timestamp"`
	Volume             float64 `json:"volume"`
	Amount             float64 `json:"amount"`
	TurnoverRate       float64 `json:"turnover_rate"`
	Amplitude          float64 `json:"amplitude"`
	LimitUp            float64 `json:"limit_up"`
	LimitDown          float64 `json:"limit_down"`
	AvgPrice           float64 `json:"avg_price"`
	VolumeRatio        float64 `json:"volume_ratio"`
	// 52 week high/low
	High52w float64 `json:"high52w"`
	Low52w  float64 `json:"low52w"`
}

// Market represents the market status information
type Market struct {
	Status   string `json:"status"`
	Region   string `json:"region"`
	TimeZone string `json:"time_zone"`
}

// GetQuoteDetail retrieves detailed quote information for a single symbol
func (c *Client) GetQuoteDetail(ctx context.Context, symbol string) (*QuoteDetailData, *request.Record, error) {
	if symbol == "" {
		return nil, nil, fmt.Errorf("symbol required")
	}

	// Ensure symbol is in Xueqiu format (e.g., SH600000)
	xueqiuSymbol, err := toXueqiuSymbol(symbol)
	if err != nil {
		// If conversion fails, try using the symbol as is
		xueqiuSymbol = symbol
	}

	params := url.Values{}
	params.Set("symbol", xueqiuSymbol)
	params.Set("extend", "detail")

	reqURL := fmt.Sprintf("%s%s?%s", "https://stock.xueqiu.com", QuoteDetailAPI, params.Encode())

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Referer":    "https://xueqiu.com/",
	}
	if c.cookie != "" {
		headers["Cookie"] = c.cookie
	}
	if c.token != "" {
		headers["X-Token"] = c.token
	}

	req := request.Request{
		Method:  "GET",
		URL:     reqURL,
		Headers: headers,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result QuoteDetailResult
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, fmt.Errorf("unmarshal response: %w", err)
	}

	if result.ErrorCode != 0 {
		return nil, record, fmt.Errorf("xueqiu error: %d %s", result.ErrorCode, result.ErrorDescription)
	}

	return &result.Data, record, nil
}
