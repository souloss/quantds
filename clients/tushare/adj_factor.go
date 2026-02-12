package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// AdjFactorParams 复权因子查询参数。
// Tushare API: adj_factor
// 获取股票复权因子，可用于计算前复权/后复权价格。
type AdjFactorParams struct {
	TSCode    string // 股票代码 (e.g., "000001.SZ")
	TradeDate string // 交易日期 (YYYYMMDD)
	StartDate string // 起始日期 (YYYYMMDD)
	EndDate   string // 结束日期 (YYYYMMDD)
}

func (p *AdjFactorParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.TSCode != "" {
		m["ts_code"] = p.TSCode
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

// AdjFactorRow 复权因子数据行。
type AdjFactorRow struct {
	TSCode    string  // 股票代码
	TradeDate string  // 交易日期 (YYYYMMDD)
	AdjFactor float64 // 复权因子（前复权价 = 未复权价 × 复权因子 / 最新复权因子）
}

// GetAdjFactor 获取股票复权因子。
// 前复权价 = 收盘价 × adj_factor / 最新 adj_factor
// 后复权价 = 收盘价 × adj_factor
func (c *Client) GetAdjFactor(ctx context.Context, params *AdjFactorParams) ([]AdjFactorRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIAdjFactor, params.ToMap(), FieldsAdjFactor)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]AdjFactorRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, AdjFactorRow{
			TSCode:    getStr(idx, item, "ts_code"),
			TradeDate: getStr(idx, item, "trade_date"),
			AdjFactor: getFlt(idx, item, "adj_factor"),
		})
	}

	return rows, record, nil
}
