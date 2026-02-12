// Package tushare provides realtime quote APIs.
//
// 实时行情相关接口，包括：
//   - realtime_quote: 实时盘口TICK快照（爬虫版）
//   - rt_k: 当日实时日线成交
package tushare

import (
	"context"
	"strconv"

	"github.com/souloss/quantds/request"
)

// RealtimeQuoteParams represents parameters for realtime tick quote request
// 实时盘口TICK快照参数
type RealtimeQuoteParams struct {
	TSCode     string // 股票代码，支持通配符*，如 6*.SH（选填）
	TradeDate  string // 交易日期，格式YYYYMMDD（选填）
	PageSize   int    // 分页大小（选填）
	PageNumber int    // 页码（选填）
}

// RealtimeQuoteRow represents a single realtime tick quote
// 实时盘口TICK快照数据
type RealtimeQuoteRow struct {
	TSCode        string  // 股票代码
	Symbol        string  // 股票代码（不含交易所后缀）
	Name          string  // 股票名称
	Trade         float64 // 当前价
	PriceChange   float64 // 涨跌额
	ChangePercent float64 // 涨跌幅(%)
	Buy           float64 // 买一价
	Sell          float64 // 卖一价
	High          float64 // 最高价
	Low           float64 // 最低价
	Open          float64 // 开盘价
	PreClose      float64 // 昨收价
	Volume        float64 // 成交量（手）
	Amount        float64 // 成交额（万元）
	TickTime      string  // 时间
}

// GetRealtimeQuote retrieves realtime tick quotes
// 获取实时盘口TICK快照数据（爬虫版）
//
// 接口：realtime_quote
// 描述：A股实时行情，本接口是tushare org版实时接口的顺延
// 数据来自网络，且不进入tushare服务器，属于爬虫接口
// 权限：0积分完全开放
//
// 限制：
//   - 需要将tushare升级到1.3.3版本以上
//   - 数据有延迟，非实时
//   - 高频请求可能被限流
func (c *Client) GetRealtimeQuote(ctx context.Context, params *RealtimeQuoteParams) ([]RealtimeQuoteRow, *request.Record, error) {
	if params == nil {
		params = &RealtimeQuoteParams{}
	}

	p := make(map[string]string)
	if params.TSCode != "" {
		p["ts_code"] = params.TSCode
	}
	if params.TradeDate != "" {
		p["trade_date"] = params.TradeDate
	}
	if params.PageSize > 0 {
		p["page_size"] = strconv.Itoa(params.PageSize)
	}
	if params.PageNumber > 0 {
		p["page_number"] = strconv.Itoa(params.PageNumber)
	}

	data, record, err := c.post(ctx, "realtime_quote", p, "")
	if err != nil {
		return nil, record, err
	}

	rows := make([]RealtimeQuoteRow, 0, len(data.Items))
	for _, item := range data.Items {
		if len(item) < 17 {
			continue
		}
		row := RealtimeQuoteRow{}
		if v, ok := item[0].(string); ok {
			row.TSCode = v
		}
		if v, ok := item[1].(string); ok {
			row.Symbol = v
		}
		if v, ok := item[2].(string); ok {
			row.Name = v
		}
		row.Trade = getFloat(item, 3)
		row.PriceChange = getFloat(item, 4)
		row.ChangePercent = getFloat(item, 5)
		row.Buy = getFloat(item, 6)
		row.Sell = getFloat(item, 7)
		row.High = getFloat(item, 8)
		row.Low = getFloat(item, 9)
		row.Open = getFloat(item, 10)
		row.PreClose = getFloat(item, 11)
		row.Volume = getFloat(item, 12)
		row.Amount = getFloat(item, 13)
		if v, ok := item[14].(string); ok {
			row.TickTime = v
		}
		rows = append(rows, row)
	}

	return rows, record, nil
}

// RtKParams represents parameters for realtime daily K-line request
// 当日实时日线成交参数
type RtKParams struct {
	TSCode    string // 股票代码，支持通配符*，如 6*.SH、000001.SZ（选填）
	TradeDate string // 交易日期，格式YYYYMMDD（选填）
}

// RtKRow represents a single realtime daily K-line data
// 当日实时日线成交数据
type RtKRow struct {
	TSCode    string  // 股票代码
	TradeDate string  // 交易日期
	Open      float64 // 开盘价
	High      float64 // 最高价
	Low       float64 // 最低价
	Close     float64 // 收盘价（当前最新价）
	PreClose  float64 // 昨收价
	Change    float64 // 涨跌额
	PctChg    float64 // 涨跌幅(%)
	Vol       float64 // 成交量（手）
	Amount    float64 // 成交额（千元）
}

// GetRtK retrieves realtime daily K-line data
// 获取当日实时日线行情数据
//
// 接口：rt_k
// 描述：获取实时日k线行情，支持按股票代码及股票代码通配符一次性提取全部股票实时日k线行情
// 权限：本接口是单独开权限的，需要用户有对应的权限才能访问
//
// 限量：单次最大可提取6000条数据，等同于一次提取全市场
//
// 参数说明：
//   - ts_code: 股票代码，支持通配符
//   - 单个股票：600000.SH、000001.SZ、430047.BJ
//   - 通配符：6*.SH（沪市主板）、0*.SZ（深市主板）、3*.SZ（创业板）
//   - trade_date: 交易日期，格式YYYYMMDD
func (c *Client) GetRtK(ctx context.Context, params *RtKParams) ([]RtKRow, *request.Record, error) {
	if params == nil {
		params = &RtKParams{}
	}

	p := make(map[string]string)
	if params.TSCode != "" {
		p["ts_code"] = params.TSCode
	}
	if params.TradeDate != "" {
		p["trade_date"] = params.TradeDate
	}

	data, record, err := c.post(ctx, "rt_k", p, "ts_code,trade_date,open,high,low,close,pre_close,change,pct_chg,vol,amount")
	if err != nil {
		return nil, record, err
	}

	rows := make([]RtKRow, 0, len(data.Items))
	for _, item := range data.Items {
		if len(item) < 11 {
			continue
		}
		row := RtKRow{}
		if v, ok := item[0].(string); ok {
			row.TSCode = v
		}
		if v, ok := item[1].(string); ok {
			row.TradeDate = v
		}
		row.Open = getFloat(item, 2)
		row.High = getFloat(item, 3)
		row.Low = getFloat(item, 4)
		row.Close = getFloat(item, 5)
		row.PreClose = getFloat(item, 6)
		row.Change = getFloat(item, 7)
		row.PctChg = getFloat(item, 8)
		row.Vol = getFloat(item, 9)
		row.Amount = getFloat(item, 10)
		rows = append(rows, row)
	}

	return rows, record, nil
}

// getFloat extracts float value from interface{} at given index
func getFloat(item []interface{}, idx int) float64 {
	if idx >= len(item) {
		return 0
	}
	switch v := item[idx].(type) {
	case float64:
		return v
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	}
	return 0
}
