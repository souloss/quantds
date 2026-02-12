package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// DailyBasicParams 每日指标查询参数。
// Tushare API: daily_basic
// 获取每日交易指标，包括 PE/PB/PS/换手率/市值等。
type DailyBasicParams struct {
	TSCode    string // 股票代码 (e.g., "000001.SZ")
	TradeDate string // 交易日期 (YYYYMMDD)
	StartDate string // 起始日期 (YYYYMMDD)
	EndDate   string // 结束日期 (YYYYMMDD)
}

func (p *DailyBasicParams) ToMap() map[string]string {
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

// DailyBasicRow 每日指标数据行。
type DailyBasicRow struct {
	TSCode       string  // 股票代码
	TradeDate    string  // 交易日期 (YYYYMMDD)
	TurnoverRate float64 // 换手率 (%)
	TurnoverRateF float64 // 换手率（自由流通股本）(%)
	VolumeRatio  float64 // 量比
	PE           float64 // 市盈率（动态）
	PETTM        float64 // 市盈率（TTM）
	PB           float64 // 市净率
	PS           float64 // 市销率
	PSTTM        float64 // 市销率（TTM）
	DvRatio      float64 // 股息率 (%)
	DvTTM        float64 // 股息率（TTM）(%)
	TotalShare   float64 // 总股本（万股）
	FloatShare   float64 // 流通股本（万股）
	TotalMV      float64 // 总市值（万元）
	CircMV       float64 // 流通市值（万元）
}

// GetDailyBasic 获取每日交易指标。
// 每天更新一次，包含 PE/PB/PS/换手率/市值等估值与交易指标。
func (c *Client) GetDailyBasic(ctx context.Context, params *DailyBasicParams) ([]DailyBasicRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIDailyBasic, params.ToMap(), FieldsDailyBasic)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]DailyBasicRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, DailyBasicRow{
			TSCode:        getStr(idx, item, "ts_code"),
			TradeDate:     getStr(idx, item, "trade_date"),
			TurnoverRate:  getFlt(idx, item, "turnover_rate"),
			TurnoverRateF: getFlt(idx, item, "turnover_rate_f"),
			VolumeRatio:   getFlt(idx, item, "volume_ratio"),
			PE:            getFlt(idx, item, "pe"),
			PETTM:         getFlt(idx, item, "pe_ttm"),
			PB:            getFlt(idx, item, "pb"),
			PS:            getFlt(idx, item, "ps"),
			PSTTM:         getFlt(idx, item, "ps_ttm"),
			DvRatio:       getFlt(idx, item, "dv_ratio"),
			DvTTM:         getFlt(idx, item, "dv_ttm"),
			TotalShare:    getFlt(idx, item, "total_share"),
			FloatShare:    getFlt(idx, item, "float_share"),
			TotalMV:       getFlt(idx, item, "total_mv"),
			CircMV:        getFlt(idx, item, "circ_mv"),
		})
	}

	return rows, record, nil
}
