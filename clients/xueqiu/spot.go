package xueqiu

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

const SpotAPI = "/v5/stock/realtime/quotec.json"

// SpotParams represents parameters for real-time spot request
type SpotParams struct {
	Symbols []string
}

// SpotResult represents the spot data result
type SpotResult struct {
	Data  []SpotQuote
	Count int
}

// SpotQuote represents a single real-time quote
type SpotQuote struct {
	Symbol     string
	Name       string
	Latest     float64
	Open       float64
	High       float64
	Low        float64
	PreClose   float64
	Change     float64
	ChangeRate float64
	Volume     float64
	Turnover   float64
	Timestamp  int64
}

// GetSpot retrieves real-time spot quotes
func (c *Client) GetSpot(ctx context.Context, params *SpotParams) (*SpotResult, *request.Record, error) {
	if params == nil || len(params.Symbols) == 0 {
		return nil, nil, fmt.Errorf("symbols required")
	}

	// Convert symbols to xueqiu format
	symbols := make([]string, 0, len(params.Symbols))
	for _, s := range params.Symbols {
		sym, err := toXueqiuSymbol(s)
		if err != nil {
			continue
		}
		symbols = append(symbols, sym)
	}

	if len(symbols) == 0 {
		return nil, nil, fmt.Errorf("no valid symbols")
	}

	symbolStr := ""
	for i, s := range symbols {
		if i > 0 {
			symbolStr += ","
		}
		symbolStr += s
	}

	url := fmt.Sprintf("https://stock.xueqiu.com%s?symbol=%s", SpotAPI, symbolStr)

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
		URL:     url,
		Headers: headers,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseSpotResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

func parseSpotResponse(body []byte) (*SpotResult, error) {
	var resp struct {
		Data []struct {
			Symbol  string  `json:"symbol"`
			Current float64 `json:"current"`
			Open    float64 `json:"open"`
			High    float64 `json:"high"`
			Low     float64 `json:"low"`
			Last    float64 `json:"last_close"`
			Chg     float64 `json:"chg"`
			Percent float64 `json:"percent"`
			Volume  float64 `json:"volume"`
			Amount  float64 `json:"amount"`
			Time    int64   `json:"time"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	quotes := make([]SpotQuote, 0, len(resp.Data))
	for _, d := range resp.Data {
		quotes = append(quotes, SpotQuote{
			Symbol:     d.Symbol,
			Latest:     d.Current,
			Open:       d.Open,
			High:       d.High,
			Low:        d.Low,
			PreClose:   d.Last,
			Change:     d.Chg,
			ChangeRate: d.Percent,
			Volume:     d.Volume,
			Turnover:   d.Amount,
			Timestamp:  d.Time,
		})
	}

	return &SpotResult{Data: quotes, Count: len(quotes)}, nil
}
