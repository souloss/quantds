package bse

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
	Code        string  // 证券代码
	Name        string  // 证券简称
	ListDate    string  // 上市日期
	TotalShares float64 // 总股本
	FloatShares float64 // 流通股本
	Industry    string  // 所属行业
}

// GetInstruments retrieves list of securities from BSE (Beijing Stock Exchange)
// This method automatically handles pagination and returns all stocks
func (c *Client) GetInstruments(ctx context.Context, params *InstrumentParams) ([]StockRow, []*request.Record, error) {
	// If params is nil or empty, use the existing GetStockList
	if params == nil || (params.Page == 0 && params.Typejb == "" && params.Xxfcbj == "") {
		return c.GetStockList(ctx)
	}

	// Otherwise, use GetStockListPage with params
	result, record, err := c.GetStockListPage(ctx, params)
	if err != nil {
		return nil, []*request.Record{record}, err
	}

	return result.Data, []*request.Record{record}, nil
}

// GetInstrumentsPage retrieves a single page of securities from BSE
// Returns InstrumentResult for consistency with other exchanges
func (c *Client) GetInstrumentsPage(ctx context.Context, params *InstrumentParams) (*InstrumentResult, *request.Record, error) {
	return c.GetStockListPage(ctx, params)
}

// ToInstrumentData converts StockRow to InstrumentData for easier access
func (r *StockRow) ToInstrumentData() InstrumentData {
	return InstrumentData{
		Code:        r.StockCode,
		Name:        r.StockName,
		ListDate:    r.ListDate,
		TotalShares: r.TotalShares,
		FloatShares: r.FloatShares,
		Industry:    r.Industry,
	}
}
