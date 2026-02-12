package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// BalanceSheetParams 资产负债表查询参数。
// Tushare API: balancesheet
// 获取上市公司资产负债表数据。
type BalanceSheetParams struct {
	TSCode    string // 股票代码 (e.g., "000001.SZ")
	AnnDate   string // 公告日期 (YYYYMMDD)
	StartDate string // 报告期起始日期 (YYYYMMDD)
	EndDate   string // 报告期结束日期 (YYYYMMDD)
	Period    string // 报告期 (e.g., "20241231")
	ReportType string // 报告类型: 1=合并报表
}

func (p *BalanceSheetParams) ToMap() map[string]string {
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

// BalanceSheetRow 资产负债表数据行。
type BalanceSheetRow struct {
	TSCode         string  // 股票代码
	AnnDate        string  // 公告日期
	FAnnDate       string  // 实际公告日期
	EndDate        string  // 报告期
	ReportType     string  // 报告类型
	CompType       string  // 公司类型
	TotalAssets    float64 // 总资产（元）
	TotalCurAssets float64 // 流动资产合计（元）
	TotalNCA       float64 // 非流动资产合计（元）
	TotalLiab      float64 // 总负债（元）
	TotalCurLiab   float64 // 流动负债合计（元）
	TotalNCL       float64 // 非流动负债合计（元）
	TotalHldrEqy   float64 // 股东权益合计（不含少数股东）（元）
	TotalHldrEqyInc float64 // 股东权益合计（含少数股东）（元）
	CapRese        float64 // 资本公积（元）
	SurplusRese    float64 // 盈余公积（元）
	UndistProfit   float64 // 未分配利润（元）
	MoneyCap       float64 // 货币资金（元）
	AccountsReceiv float64 // 应收账款（元）
	Inventories    float64 // 存货（元）
	FixAssets      float64 // 固定资产（元）
}

// GetBalanceSheet 获取资产负债表数据。
// 建议使用 report_type="1" 获取合并报表数据。
func (c *Client) GetBalanceSheet(ctx context.Context, params *BalanceSheetParams) ([]BalanceSheetRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIBalanceSheet, params.ToMap(), FieldsBalance)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]BalanceSheetRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, BalanceSheetRow{
			TSCode:          getStr(idx, item, "ts_code"),
			AnnDate:         getStr(idx, item, "ann_date"),
			FAnnDate:        getStr(idx, item, "f_ann_date"),
			EndDate:         getStr(idx, item, "end_date"),
			ReportType:      getStr(idx, item, "report_type"),
			CompType:        getStr(idx, item, "comp_type"),
			TotalAssets:     getFlt(idx, item, "total_assets"),
			TotalCurAssets:  getFlt(idx, item, "total_cur_assets"),
			TotalNCA:        getFlt(idx, item, "total_nca"),
			TotalLiab:       getFlt(idx, item, "total_liab"),
			TotalCurLiab:    getFlt(idx, item, "total_cur_liab"),
			TotalNCL:        getFlt(idx, item, "total_ncl"),
			TotalHldrEqy:    getFlt(idx, item, "total_hldr_eqy_exc_min_int"),
			TotalHldrEqyInc: getFlt(idx, item, "total_hldr_eqy_inc_min_int"),
			CapRese:         getFlt(idx, item, "cap_rese"),
			SurplusRese:     getFlt(idx, item, "surplus_rese"),
			UndistProfit:    getFlt(idx, item, "undist_profit"),
			MoneyCap:        getFlt(idx, item, "money_cap"),
			AccountsReceiv:  getFlt(idx, item, "accounts_receiv"),
			Inventories:     getFlt(idx, item, "inventories"),
			FixAssets:       getFlt(idx, item, "fix_assets"),
		})
	}

	return rows, record, nil
}
