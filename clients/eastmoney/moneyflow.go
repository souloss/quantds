package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/souloss/quantds/request"
)

const (
	MoneyFlowAPI        = "/api/qt/stock/fflow/get"
	MoneyFlowHistoryAPI = "/api/qt/stock/fflow/kline/get"
)

// Realtime money flow fields
// f62: 主力净流入
// f184: 主力净流入占比
// f66: 超大单净流入
// f69: 大单净流入
// f72: 中单净流入
// f75: 小单净流入
const MoneyFlowFields = "f62,f184,f66,f69,f72,f75"

// MoneyFlowParams represents parameters for money flow request
type MoneyFlowParams struct {
	Symbol string
}

// MoneyFlowData represents real-time money flow data
type MoneyFlowData struct {
	MainNetInflow      float64 `json:"main_net_inflow"`      // 主力净流入
	MainNetInflowRatio float64 `json:"main_net_inflow_ratio"` // 主力净流入占比
	SuperNetInflow     float64 `json:"super_net_inflow"`     // 超大单净流入
	LargeNetInflow     float64 `json:"large_net_inflow"`     // 大单净流入
	MediumNetInflow    float64 `json:"medium_net_inflow"`    // 中单净流入
	SmallNetInflow     float64 `json:"small_net_inflow"`     // 小单净流入
}

// GetMoneyFlow retrieves real-time money flow for a stock
func (c *Client) GetMoneyFlow(ctx context.Context, params *MoneyFlowParams) (*MoneyFlowData, *request.Record, error) {
	if params.Symbol == "" {
		return nil, nil, fmt.Errorf("symbol required")
	}

	secid, err := toEastMoneySecid(params.Symbol)
	if err != nil {
		return nil, nil, err
	}

	query := url.Values{}
	query.Set("lmt", "0")
	query.Set("klt", "1")
	query.Set("secid", secid)
	query.Set("fields", MoneyFlowFields)

	url := fmt.Sprintf("%s%s?%s", PushURL, MoneyFlowAPI, query.Encode())

	req := request.Request{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":    "https://data.eastmoney.com/",
		},
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := parseMoneyFlowResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}
	return data, record, nil
}

func parseMoneyFlowResponse(body []byte) (*MoneyFlowData, error) {
	var raw struct {
		Data *struct {
			MainNetInflow      float64 `json:"f62"`
			MainNetInflowRatio float64 `json:"f184"`
			SuperNetInflow     float64 `json:"f66"`
			LargeNetInflow     float64 `json:"f69"`
			MediumNetInflow    float64 `json:"f72"`
			SmallNetInflow     float64 `json:"f75"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	if raw.Data == nil {
		return nil, fmt.Errorf("no data found")
	}

	return &MoneyFlowData{
		MainNetInflow:      raw.Data.MainNetInflow,
		MainNetInflowRatio: raw.Data.MainNetInflowRatio,
		SuperNetInflow:     raw.Data.SuperNetInflow,
		LargeNetInflow:     raw.Data.LargeNetInflow,
		MediumNetInflow:    raw.Data.MediumNetInflow,
		SmallNetInflow:     raw.Data.SmallNetInflow,
	}, nil
}

// MoneyFlowHistoryParams represents parameters for money flow history
type MoneyFlowHistoryParams struct {
	Symbol string
	Limit  int // Number of days, default 100
}

// MoneyFlowHistoryItem represents a single day's money flow
type MoneyFlowHistoryItem struct {
	Date               string  `json:"date"`
	MainNetInflow      float64 `json:"main_net_inflow"`
	MainNetInflowRatio float64 `json:"main_net_inflow_ratio"`
	SuperNetInflow     float64 `json:"super_net_inflow"`
	LargeNetInflow     float64 `json:"large_net_inflow"`
	MediumNetInflow    float64 `json:"medium_net_inflow"`
	SmallNetInflow     float64 `json:"small_net_inflow"`
	ClosePrice         float64 `json:"close_price"`
	ChangePercent      float64 `json:"change_percent"`
}

// GetMoneyFlowHistory retrieves historical money flow
func (c *Client) GetMoneyFlowHistory(ctx context.Context, params *MoneyFlowHistoryParams) ([]MoneyFlowHistoryItem, *request.Record, error) {
	if params.Symbol == "" {
		return nil, nil, fmt.Errorf("symbol required")
	}

	secid, err := toEastMoneySecid(params.Symbol)
	if err != nil {
		return nil, nil, err
	}

	limit := params.Limit
	if limit <= 0 {
		limit = 100
	}

	query := url.Values{}
	query.Set("lmt", fmt.Sprintf("%d", limit))
	query.Set("klt", "101") // Daily
	query.Set("secid", secid)
	query.Set("fields1", "f1,f2,f3,f7")
	query.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63")

	url := fmt.Sprintf("%s%s?%s", PushURL, MoneyFlowHistoryAPI, query.Encode())

	req := request.Request{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":    "https://data.eastmoney.com/",
		},
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	items, err := parseMoneyFlowHistoryResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}
	return items, record, nil
}

func parseMoneyFlowHistoryResponse(body []byte) ([]MoneyFlowHistoryItem, error) {
	var raw struct {
		Data *struct {
			Klines []string `json:"klines"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	if raw.Data == nil {
		return nil, fmt.Errorf("no data found")
	}

	items := make([]MoneyFlowHistoryItem, 0, len(raw.Data.Klines))
	for _, line := range raw.Data.Klines {
		parts := strings.Split(line, ",")
		// Format: date, close, change_pct, main_net, main_ratio, super_net, super_ratio, large_net, large_ratio, medium_net, medium_ratio, small_net, small_ratio
		if len(parts) < 13 {
			continue
		}
		
		item := MoneyFlowHistoryItem{
			Date:               parts[0],
			ClosePrice:         parseFloat(parts[1]),
			ChangePercent:      parseFloat(parts[2]),
			MainNetInflow:      parseFloat(parts[3]),
			MainNetInflowRatio: parseFloat(parts[4]),
			SuperNetInflow:     parseFloat(parts[5]),
			LargeNetInflow:     parseFloat(parts[7]),
			MediumNetInflow:    parseFloat(parts[9]),
			SmallNetInflow:     parseFloat(parts[11]),
		}
		items = append(items, item)
	}

	return items, nil
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
