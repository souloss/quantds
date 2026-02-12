package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// DividendParams 分红送股查询参数。
// Tushare API: dividend
// 获取上市公司分红送股信息。
type DividendParams struct {
	TSCode     string // 股票代码 (e.g., "000001.SZ")
	AnnDate    string // 公告日期 (YYYYMMDD)
	ExDate     string // 除权除息日 (YYYYMMDD)
	RecordDate string // 股权登记日 (YYYYMMDD)
	ImpAnnDate string // 实施公告日 (YYYYMMDD)
}

func (p *DividendParams) ToMap() map[string]string {
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
	if p.ExDate != "" {
		m["ex_date"] = p.ExDate
	}
	if p.RecordDate != "" {
		m["record_date"] = p.RecordDate
	}
	if p.ImpAnnDate != "" {
		m["imp_ann_date"] = p.ImpAnnDate
	}
	return m
}

// DividendRow 分红送股数据行。
type DividendRow struct {
	TSCode      string  // 股票代码
	AnnDate     string  // 公告日期
	DivProc     string  // 实施进度
	StkDiv      float64 // 每股送转 (股)
	StkBoRate   float64 // 每股送股比例
	StkCoRate   float64 // 每股转增比例
	CashDiv     float64 // 每股分红 (税前)
	CashDivTax  float64 // 每股分红 (税后)
	RecordDate  string  // 股权登记日
	ExDate      string  // 除权除息日
	PayDate     string  // 派息日
	DivListDate string  // 红股上市日
	ImpAnnDate  string  // 实施公告日
}

// GetDividend 获取分红送股信息。
func (c *Client) GetDividend(ctx context.Context, params *DividendParams) ([]DividendRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIDividend, params.ToMap(), FieldsDividend)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]DividendRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, DividendRow{
			TSCode:      getStr(idx, item, "ts_code"),
			AnnDate:     getStr(idx, item, "ann_date"),
			DivProc:     getStr(idx, item, "div_proc"),
			StkDiv:      getFlt(idx, item, "stk_div"),
			StkBoRate:   getFlt(idx, item, "stk_bo_rate"),
			StkCoRate:   getFlt(idx, item, "stk_co_rate"),
			CashDiv:     getFlt(idx, item, "cash_div"),
			CashDivTax:  getFlt(idx, item, "cash_div_tax"),
			RecordDate:  getStr(idx, item, "record_date"),
			ExDate:      getStr(idx, item, "ex_date"),
			PayDate:     getStr(idx, item, "pay_date"),
			DivListDate: getStr(idx, item, "div_listdate"),
			ImpAnnDate:  getStr(idx, item, "imp_ann_date"),
		})
	}

	return rows, record, nil
}
