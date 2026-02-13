package szse

import (
	"context"

	"github.com/souloss/quantds/request"
)

// InstrumentParams is an alias for StockListParams
type InstrumentParams = StockListParams

// InstrumentResult is an alias for StockListResult
type InstrumentResult = StockListResult

// InstrumentData represents a single security information parsed from Excel data
type InstrumentData struct {
	Board     string // 板块
	Code      string // 公司代码
	Name      string // 公司简称
	FullName  string // 公司全称
	ACode     string // A股代码
	AName     string // A股简称
	ListDate  string // A股上市日期
	BCode     string // B股代码
	BName     string // B股简称
	BListDate string // B股上市日期
	Area      string // 地区
	Province  string // 省份
	City      string // 城市
	Industry  string // 行业
	Website   string // 公司网站
}

// GetInstruments retrieves list of securities from SZSE (Shenzhen Stock Exchange)
// This is a convenience wrapper around GetStockList that provides a more generic interface
func (c *Client) GetInstruments(ctx context.Context, params *InstrumentParams) (*InstrumentResult, *request.Record, error) {
	// Directly call GetStockList since types are compatible
	return c.GetStockList(ctx, params)
}

// ParseInstruments parses the Excel data into structured InstrumentData
// The caller should use this after calling GetInstruments
func ParseInstruments(rows [][]string) []InstrumentData {
	if len(rows) <= 1 {
		return nil
	}

	// Skip header row (first row)
	instruments := make([]InstrumentData, 0, len(rows)-1)
	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		if len(row) < 7 {
			continue // skip malformed rows
		}

		inst := InstrumentData{
			Board: safeGet(row, 0),
			Code:  safeGet(row, 1),
			Name:  safeGet(row, 2),
		}

		// Try to get A股代码 and A股简称
		if aCode := safeGet(row, 4); aCode != "" {
			inst.ACode = aCode
			inst.AName = safeGet(row, 5)
			inst.ListDate = safeGet(row, 6)
			inst.Code = aCode
			inst.Name = inst.AName
		}

		// Get additional fields if available
		if len(row) > 16 {
			inst.Province = safeGet(row, 14)
			inst.City = safeGet(row, 15)
			inst.Industry = safeGet(row, 16)
		}

		if inst.Code != "" {
			instruments = append(instruments, inst)
		}
	}

	return instruments
}

func safeGet(row []string, idx int) string {
	if idx < len(row) {
		return row[idx]
	}
	return ""
}
