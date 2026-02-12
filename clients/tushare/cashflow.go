package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// CashflowParams 现金流量表查询参数。
// Tushare API: cashflow
// 获取上市公司现金流量表数据。
type CashflowParams struct {
	TSCode    string // 股票代码 (e.g., "000001.SZ")
	AnnDate   string // 公告日期 (YYYYMMDD)
	StartDate string // 报告期起始日期 (YYYYMMDD)
	EndDate   string // 报告期结束日期 (YYYYMMDD)
	Period    string // 报告期 (e.g., "20241231")
	ReportType string // 报告类型: 1=合并报表
}

func (p *CashflowParams) ToMap() map[string]string {
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
	if p.ReportType != "" {
		m["report_type"] = p.ReportType
	}
	return m
}

// CashflowRow 现金流量表数据行。
type CashflowRow struct {
	TSCode        string  // 股票代码
	AnnDate       string  // 公告日期
	FAnnDate      string  // 实际公告日期
	EndDate       string  // 报告期
	ReportType    string  // 报告类型
	CompType      string  // 公司类型
	NCashflowAct  float64 // 经营活动产生的现金流量净额（元）
	NCashflowInv  float64 // 投资活动产生的现金流量净额（元）
	NCashflowFnc  float64 // 筹资活动产生的现金流量净额（元）
	CCashEquEnd   float64 // 期末现金及现金等价物余额（元）
}

// GetCashflow 获取现金流量表数据。
// 建议使用 report_type="1" 获取合并报表数据。
func (c *Client) GetCashflow(ctx context.Context, params *CashflowParams) ([]CashflowRow, *request.Record, error) {
	data, record, err := c.post(ctx, APICashflow, params.ToMap(), FieldsCashflow)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]CashflowRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, CashflowRow{
			TSCode:       getStr(idx, item, "ts_code"),
			AnnDate:      getStr(idx, item, "ann_date"),
			FAnnDate:     getStr(idx, item, "f_ann_date"),
			EndDate:      getStr(idx, item, "end_date"),
			ReportType:   getStr(idx, item, "report_type"),
			CompType:     getStr(idx, item, "comp_type"),
			NCashflowAct: getFlt(idx, item, "n_cashflow_act"),
			NCashflowInv: getFlt(idx, item, "n_cashflow_inv_act"),
			NCashflowFnc: getFlt(idx, item, "n_cash_flows_fnc_act"),
			CCashEquEnd:  getFlt(idx, item, "c_cash_equ_end_period"),
		})
	}

	return rows, record, nil
}
