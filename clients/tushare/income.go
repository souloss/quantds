package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// IncomeParams 利润表查询参数。
// Tushare API: income
// 获取上市公司利润表数据。
type IncomeParams struct {
	TSCode    string // 股票代码 (e.g., "000001.SZ")
	AnnDate   string // 公告日期 (YYYYMMDD)
	StartDate string // 报告期起始日期 (YYYYMMDD)
	EndDate   string // 报告期结束日期 (YYYYMMDD)
	Period    string // 报告期 (e.g., "20241231")
	ReportType string // 报告类型: 1=合并报表, 4=调整后合并报表
}

func (p *IncomeParams) ToMap() map[string]string {
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

// IncomeRow 利润表数据行。
type IncomeRow struct {
	TSCode       string  // 股票代码
	AnnDate      string  // 公告日期
	FAnnDate     string  // 实际公告日期
	EndDate      string  // 报告期 (e.g., "20241231")
	ReportType   string  // 报告类型
	CompType     string  // 公司类型: 1=一般工商业, 2=银行, 3=保险, 4=证券
	BasicEPS     float64 // 基本每股收益（元）
	DilutedEPS   float64 // 稀释每股收益（元）
	TotalRevenue float64 // 营业总收入（元）
	Revenue      float64 // 营业收入（元）
	TotalCogs    float64 // 营业总成本（元）
	OperCost     float64 // 营业成本（元）
	SellExp      float64 // 销售费用（元）
	AdminExp     float64 // 管理费用（元）
	FinExp       float64 // 财务费用（元）
	RDExp        float64 // 研发费用（元）
	OperProfit   float64 // 营业利润（元）
	TotalProfit  float64 // 利润总额（元）
	NIncome      float64 // 净利润（元）
	NIncomeAttrP float64 // 归属母公司股东的净利润（元）
	EBIT         float64 // 息税前利润（元）
	EBITDA       float64 // 息税折旧摊销前利润（元）
}

// GetIncome 获取利润表数据。
// 建议使用 report_type="1" 获取合并报表数据。
func (c *Client) GetIncome(ctx context.Context, params *IncomeParams) ([]IncomeRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIIncome, params.ToMap(), FieldsIncome)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]IncomeRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, IncomeRow{
			TSCode:       getStr(idx, item, "ts_code"),
			AnnDate:      getStr(idx, item, "ann_date"),
			FAnnDate:     getStr(idx, item, "f_ann_date"),
			EndDate:      getStr(idx, item, "end_date"),
			ReportType:   getStr(idx, item, "report_type"),
			CompType:     getStr(idx, item, "comp_type"),
			BasicEPS:     getFlt(idx, item, "basic_eps"),
			DilutedEPS:   getFlt(idx, item, "diluted_eps"),
			TotalRevenue: getFlt(idx, item, "total_revenue"),
			Revenue:      getFlt(idx, item, "revenue"),
			TotalCogs:    getFlt(idx, item, "total_cogs"),
			OperCost:     getFlt(idx, item, "oper_cost"),
			SellExp:      getFlt(idx, item, "sell_exp"),
			AdminExp:     getFlt(idx, item, "admin_exp"),
			FinExp:       getFlt(idx, item, "fin_exp"),
			RDExp:        getFlt(idx, item, "rd_exp"),
			OperProfit:   getFlt(idx, item, "oper_profit"),
			TotalProfit:  getFlt(idx, item, "total_profit"),
			NIncome:      getFlt(idx, item, "n_income"),
			NIncomeAttrP: getFlt(idx, item, "n_income_attr_p"),
			EBIT:         getFlt(idx, item, "ebit"),
			EBITDA:       getFlt(idx, item, "ebitda"),
		})
	}

	return rows, record, nil
}
