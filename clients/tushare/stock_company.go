package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// StockCompanyParams 上市公司基本信息查询参数。
// Tushare API: stock_company
// 获取上市公司基本信息，包括公司治理、注册资本等。
type StockCompanyParams struct {
	TSCode   string // 股票代码 (e.g., "000001.SZ")
	Exchange string // 交易所: SSE=上交所, SZSE=深交所, BSE=北交所
}

func (p *StockCompanyParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.TSCode != "" {
		m["ts_code"] = p.TSCode
	}
	if p.Exchange != "" {
		m["exchange"] = p.Exchange
	}
	return m
}

// StockCompanyRow 上市公司基本信息数据行。
type StockCompanyRow struct {
	TSCode        string  // 股票代码
	Chairman      string  // 法人代表/董事长
	Manager       string  // 总经理
	Secretary     string  // 董秘
	RegCapital    float64 // 注册资本（万元）
	SetupDate     string  // 成立日期 (YYYYMMDD)
	Province      string  // 所属省份
	City          string  // 所属城市
	Introduction  string  // 公司简介
	Website       string  // 公司网站
	Email         string  // 电子邮箱
	Office        string  // 办公地址
	Employees     int     // 员工人数
	MainBusiness  string  // 主营业务
	BusinessScope string  // 经营范围
}

// GetStockCompany 获取上市公司基本信息。
// 包含法人代表、注册资本、成立日期、主营业务、经营范围等。
func (c *Client) GetStockCompany(ctx context.Context, params *StockCompanyParams) ([]StockCompanyRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIStockCompany, params.ToMap(), FieldsCompany)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]StockCompanyRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, StockCompanyRow{
			TSCode:        getStr(idx, item, "ts_code"),
			Chairman:      getStr(idx, item, "chairman"),
			Manager:       getStr(idx, item, "manager"),
			Secretary:     getStr(idx, item, "secretary"),
			RegCapital:    getFlt(idx, item, "reg_capital"),
			SetupDate:     getStr(idx, item, "setup_date"),
			Province:      getStr(idx, item, "province"),
			City:          getStr(idx, item, "city"),
			Introduction:  getStr(idx, item, "introduction"),
			Website:       getStr(idx, item, "website"),
			Email:         getStr(idx, item, "email"),
			Office:        getStr(idx, item, "office"),
			Employees:     getInt(idx, item, "employees"),
			MainBusiness:  getStr(idx, item, "main_business"),
			BusinessScope: getStr(idx, item, "business_scope"),
		})
	}

	return rows, record, nil
}
