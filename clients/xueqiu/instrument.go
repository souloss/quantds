// Package xueqiu provides stock list (instrument) API.
//
// 雪球证券列表API，支持获取沪深京A股列表。
package xueqiu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/souloss/quantds/request"
)

const (
	StockListAPI = "/v5/stock/screener/quote/list.json"
)

// StockListParams represents parameters for stock list request
// 股票列表查询参数
type StockListParams struct {
	Market    string // 市场代码：CN（A股）、HK（港股）、US（美股）
	Exchange  string // 交易所：SH（上交所）、SZ（深交所）、BJ（北交所）
	BoardType string // 板块类型：all（全部）、main（主板）、gem（创业板）、star（科创板）、bse（北交所）
	Order     string // 排序字段：symbol、current、percent、amount 等
	OrderBy   string // 排序方式：asc、desc
	Page      int    // 页码，从1开始
	Size      int    // 每页数量，最大90
}

// StockListResult represents the result of stock list query
// 股票列表查询结果
type StockListResult struct {
	Total   int         // 总数
	Count   int         // 当前页数量
	Page    int         // 当前页码
	MaxPage int         // 最大页码
	Items   []StockItem // 股票列表
}

// StockItem represents a single stock in the list
// 单只股票信息
type StockItem struct {
	Symbol         string  // 标准代码（如 SH600000）
	Code           string  // 股票代码（不含交易所）
	Name           string  // 股票名称
	Current        float64 // 当前价
	Percent        float64 // 涨跌幅(%)
	Change         float64 // 涨跌额
	High           float64 // 最高价
	Low            float64 // 最低价
	Open           float64 // 开盘价
	PreClose       float64 // 昨收价
	Volume         float64 // 成交量（手）
	Amount         float64 // 成交额
	TurnoverRate   float64 // 换手率(%)
	PE             float64 // 市盈率
	PB             float64 // 市净率
	TotalMarketCap float64 // 总市值
	FloatMarketCap float64 // 流通市值
	MarketCapital  float64 // 总市值（字段别名）
	Status         int     // 状态：1-正常，2-停牌，3-退市
	Type           int     // 类型
}

// GetStockList retrieves stock list from Xueqiu
// 获取股票列表
//
// API: /v5/stock/screener/quote/list.json
// 方法: GET
// 认证: 需要Cookie或Token
//
// 参数说明：
//   - market: 市场代码，CN表示A股
//   - exchange: 交易所，SH/SZ/BJ
//   - board: 板块类型，all/main/gem/star/bse
//   - order: 排序字段
//   - order_by: 排序方式
//   - page: 页码，从1开始
//   - size: 每页数量，最大90
//
// 限制：
//   - 需要认证（Cookie或Token）
//   - 单次最多返回90条
//   - 高频请求会被限流
func (c *Client) GetStockList(ctx context.Context, params *StockListParams) (*StockListResult, *request.Record, error) {
	if params == nil {
		params = &StockListParams{}
	}

	// 默认值
	if params.Market == "" {
		params.Market = "CN"
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 90
	}
	if params.Size > 90 {
		params.Size = 90
	}

	query := url.Values{}
	query.Set("market", params.Market)
	if params.Exchange != "" {
		query.Set("exchange", params.Exchange)
	}
	if params.BoardType != "" {
		query.Set("board", params.BoardType)
	}
	if params.Order != "" {
		query.Set("order", params.Order)
	}
	if params.OrderBy != "" {
		query.Set("order_by", params.OrderBy)
	}
	query.Set("page", strconv.Itoa(params.Page))
	query.Set("size", strconv.Itoa(params.Size))

	reqURL := fmt.Sprintf("%s%s?%s", BaseURL, StockListAPI, query.Encode())

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

	result, err := parseStockListResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

func parseStockListResponse(body []byte) (*StockListResult, error) {
	var raw struct {
		Total   int `json:"count"`
		MaxPage int `json:"maxPage"`
		Page    int `json:"page"`
		List    []struct {
			Symbol         string      `json:"symbol"`
			Name           string      `json:"name"`
			Current        interface{} `json:"current"`
			Percent        interface{} `json:"percent"`
			Change         interface{} `json:"chg"`
			High           interface{} `json:"high"`
			Low            interface{} `json:"low"`
			Open           interface{} `json:"open"`
			PreClose       interface{} `json:"last_close"`
			Volume         interface{} `json:"volume"`
			Amount         interface{} `json:"amount"`
			TurnoverRate   interface{} `json:"turnover_rate"`
			PE             interface{} `json:"pe_ttm"`
			PB             interface{} `json:"pb"`
			TotalMarketCap interface{} `json:"market_capital"`
			FloatMarketCap interface{} `json:"float_market_capital"`
			Status         int         `json:"status"`
			Type           int         `json:"type"`
		} `json:"list"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	items := make([]StockItem, 0, len(raw.List))
	for _, l := range raw.List {
		item := StockItem{
			Symbol:         l.Symbol,
			Name:           l.Name,
			Current:        toFloat(l.Current),
			Percent:        toFloat(l.Percent),
			Change:         toFloat(l.Change),
			High:           toFloat(l.High),
			Low:            toFloat(l.Low),
			Open:           toFloat(l.Open),
			PreClose:       toFloat(l.PreClose),
			Volume:         toFloat(l.Volume),
			Amount:         toFloat(l.Amount),
			TurnoverRate:   toFloat(l.TurnoverRate),
			PE:             toFloat(l.PE),
			PB:             toFloat(l.PB),
			TotalMarketCap: toFloat(l.TotalMarketCap),
			FloatMarketCap: toFloat(l.FloatMarketCap),
			MarketCapital:  toFloat(l.TotalMarketCap),
			Status:         l.Status,
			Type:           l.Type,
		}
		// 提取股票代码
		if len(l.Symbol) > 2 {
			item.Code = l.Symbol[2:]
		}
		items = append(items, item)
	}

	return &StockListResult{
		Total:   raw.Total,
		Count:   len(items),
		Page:    raw.Page,
		MaxPage: raw.MaxPage,
		Items:   items,
	}, nil
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}
