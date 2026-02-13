package sse

import (
	"context"

	"github.com/souloss/quantds/request"
)

// InstrumentParams is an alias for StockListParams
type InstrumentParams = StockListParams

// InstrumentResult is an alias for StockListResult
type InstrumentResult = StockListResult

// InstrumentData represents a single security information
// This is a convenience wrapper around StockRow
type InstrumentData struct {
	Code        string // 公司代码
	Name        string // 公司简称
	FullName    string // 公司全称
	ListDate    string // 上市日期
	TotalShares string // 总股本
	FloatShares string // 流通股本
	Industry    string // 行业名称
}

// GetInstruments retrieves list of securities from SSE (Shanghai Stock Exchange)
// This is a convenience wrapper around GetStockList that provides a more generic interface
func (c *Client) GetInstruments(ctx context.Context, params *InstrumentParams) (*InstrumentResult, *request.Record, error) {
	// Directly call GetStockList since types are compatible
	return c.GetStockList(ctx, params)
}

// ToInstrumentData converts StockRow to InstrumentData for easier access
func (r *StockRow) ToInstrumentData() InstrumentData {
	return InstrumentData{
		Code:        r.CompanyCode,
		Name:        r.CompanyAbbr,
		FullName:    r.SecNameCn,
		ListDate:    r.ListDate,
		TotalShares: r.TotalShares,
		FloatShares: r.FloatShares,
		Industry:    r.Industry,
	}
}
