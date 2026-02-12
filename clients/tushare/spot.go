package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// SpotParams 实时行情/日线行情查询参数
// Tushare API: daily
type SpotParams struct {
	TsCode    string // 股票代码
	TradeDate string // 交易日期 (YYYYMMDD)
	StartDate string // 开始日期 (YYYYMMDD)
	EndDate   string // 结束日期 (YYYYMMDD)
}

func (p *SpotParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.TsCode != "" {
		m["ts_code"] = p.TsCode
	}
	if p.TradeDate != "" {
		m["trade_date"] = p.TradeDate
	}
	if p.StartDate != "" {
		m["start_date"] = p.StartDate
	}
	if p.EndDate != "" {
		m["end_date"] = p.EndDate
	}
	return m
}

// SpotRow 日线/行情数据行
type SpotRow struct {
	TsCode    string  // 股票代码
	TradeDate string  // 交易日期
	Open      float64 // 开盘价
	High      float64 // 最高价
	Low       float64 // 最低价
	Close     float64 // 收盘价
	PreClose  float64 // 昨收价
	Change    float64 // 涨跌额
	PctChg    float64 // 涨跌幅
	Vol       float64 // 成交量 (手)
	Amount    float64 // 成交额 (千元)
}

// GetSpot 获取行情数据 (基于 daily 接口)
// 注意：Tushare Pro 的 daily 接口通常在收盘后更新，盘中可能无法获取当日实时数据。
// 这里作为"Spot"的实现，通常用于获取最近一个交易日的收盘数据作为参考。
func (c *Client) GetSpot(ctx context.Context, params *SpotParams) ([]SpotRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIDaily, params.ToMap(), FieldsDaily)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]SpotRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, SpotRow{
			TsCode:    getStr(idx, item, "ts_code"),
			TradeDate: getStr(idx, item, "trade_date"),
			Open:      getFlt(idx, item, "open"),
			High:      getFlt(idx, item, "high"),
			Low:       getFlt(idx, item, "low"),
			Close:     getFlt(idx, item, "close"),
			PreClose:  getFlt(idx, item, "pre_close"),
			Change:    getFlt(idx, item, "change"),
			PctChg:    getFlt(idx, item, "pct_chg"),
			Vol:       getFlt(idx, item, "vol"),
			Amount:    getFlt(idx, item, "amount"),
		})
	}

	return rows, record, nil
}
