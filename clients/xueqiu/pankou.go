package xueqiu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/souloss/quantds/request"
)

const (
	PankouAPI = "/v5/stock/realtime/pankou.json"
)

// PankouResult represents the response for order book (pankou)
type PankouResult struct {
	Data             PankouData `json:"data"`
	ErrorCode        int        `json:"error_code"`
	ErrorDescription string     `json:"error_description"`
}

// PankouData contains the order book levels
type PankouData struct {
	Symbol    string  `json:"symbol"`
	Timestamp int64   `json:"timestamp"`
	Bp1       float64 `json:"bp1"`
	Bc1       float64 `json:"bc1"`
	Bp2       float64 `json:"bp2"`
	Bc2       float64 `json:"bc2"`
	Bp3       float64 `json:"bp3"`
	Bc3       float64 `json:"bc3"`
	Bp4       float64 `json:"bp4"`
	Bc4       float64 `json:"bc4"`
	Bp5       float64 `json:"bp5"`
	Bc5       float64 `json:"bc5"`
	Sp1       float64 `json:"sp1"`
	Sc1       float64 `json:"sc1"`
	Sp2       float64 `json:"sp2"`
	Sc2       float64 `json:"sc2"`
	Sp3       float64 `json:"sp3"`
	Sc3       float64 `json:"sc3"`
	Sp4       float64 `json:"sp4"`
	Sc4       float64 `json:"sc4"`
	Sp5       float64 `json:"sp5"`
	Sc5       float64 `json:"sc5"`
}

// GetPankou retrieves the order book (pankou) for a symbol
func (c *Client) GetPankou(ctx context.Context, symbol string) (*PankouData, *request.Record, error) {
	if symbol == "" {
		return nil, nil, fmt.Errorf("symbol required")
	}

	xueqiuSymbol, err := toXueqiuSymbol(symbol)
	if err != nil {
		xueqiuSymbol = symbol
	}

	params := url.Values{}
	params.Set("symbol", xueqiuSymbol)

	reqURL := fmt.Sprintf("%s%s?%s", BaseURL, PankouAPI, params.Encode())

	req := request.Request{
		Method:  "GET",
		URL:     reqURL,
		Headers: c.buildHeaders(),
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result PankouResult
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, fmt.Errorf("unmarshal response: %w", err)
	}

	if result.ErrorCode != 0 {
		return nil, record, fmt.Errorf("xueqiu error: %d %s", result.ErrorCode, result.ErrorDescription)
	}

	return &result.Data, record, nil
}
