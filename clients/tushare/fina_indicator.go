package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// FinaIndicatorParams 财务指标查询参数。
// Tushare API: fina_indicator
// 获取上市公司主要财务指标，如 ROE/ROA/毛利率/净利率等。
type FinaIndicatorParams struct {
	TSCode    string // 股票代码 (e.g., "000001.SZ")
	AnnDate   string // 公告日期 (YYYYMMDD)
	StartDate string // 报告期起始日期 (YYYYMMDD)
	EndDate   string // 报告期结束日期 (YYYYMMDD)
	Period    string // 报告期 (e.g., "20241231")
}

func (p *FinaIndicatorParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.TSCode != "" {
		m["ts_code"] = p.TSCode
	}
	if p.AnnDate != "" {
		m["ann_date"] = p.AnnDate
	}
	if p.StartDate != "" {
		m["start_date"] = p.StartDate
	}
	if p.EndDate != "" {
		m["end_date"] = p.EndDate
	}
	if p.Period != "" {
		m["period"] = p.Period
	}
	return m
}

// FinaIndicatorRow 财务指标数据行。
type FinaIndicatorRow struct {
	TSCode          string  // 股票代码
	AnnDate         string  // 公告日期
	EndDate         string  // 报告期
	ROE             float64 // 净资产收益率 (%)
	ROEWAA          float64 // 加权平均净资产收益率 (%)
	ROA             float64 // 总资产报酬率 (%)
	NetProfitMargin float64 // 销售净利率 (%)
	GrossProfitMargin float64 // 销售毛利率 (%)
	CurrentRatio    float64 // 流动比率
	QuickRatio      float64 // 速动比率
	DebtToAssets    float64 // 资产负债率 (%)
	TurnDays        float64 // 营业周期（天）
	ROAYearly       float64 // 年化总资产净利率 (%)
	ROEAvg          float64 // 平均净资产收益率 (%)
	AssetsTurn      float64 // 总资产周转率
	OpIncome        float64 // 经营活动净收益（元）
	EBIT            float64 // 息税前利润（元）
	EBITDA          float64 // 息税折旧摊销前利润（元）
}

// GetFinaIndicator 获取财务指标数据。
// 包含 ROE/ROA/毛利率/净利率/资产负债率/流动比率/速动比率等关键指标。
func (c *Client) GetFinaIndicator(ctx context.Context, params *FinaIndicatorParams) ([]FinaIndicatorRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIFinaIndicator, params.ToMap(), FieldsFinaInd)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]FinaIndicatorRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, FinaIndicatorRow{
			TSCode:            getStr(idx, item, "ts_code"),
			AnnDate:           getStr(idx, item, "ann_date"),
			EndDate:           getStr(idx, item, "end_date"),
			ROE:               getFlt(idx, item, "roe"),
			ROEWAA:            getFlt(idx, item, "roe_waa"),
			ROA:               getFlt(idx, item, "roa"),
			NetProfitMargin:   getFlt(idx, item, "netprofit_margin"),
			GrossProfitMargin: getFlt(idx, item, "grossprofit_margin"),
			CurrentRatio:      getFlt(idx, item, "current_ratio"),
			QuickRatio:        getFlt(idx, item, "quick_ratio"),
			DebtToAssets:      getFlt(idx, item, "debt_to_assets"),
			TurnDays:          getFlt(idx, item, "turn_days"),
			ROAYearly:         getFlt(idx, item, "roa_yearly"),
			ROEAvg:            getFlt(idx, item, "roe_avg"),
			AssetsTurn:        getFlt(idx, item, "assets_turn"),
			OpIncome:          getFlt(idx, item, "op_income"),
			EBIT:              getFlt(idx, item, "ebit"),
			EBITDA:            getFlt(idx, item, "ebitda"),
		})
	}

	return rows, record, nil
}
