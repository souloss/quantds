package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

type StockBasicParams struct {
	TSCode   string
	Exchange string
	Status   string
	Limit    int
	Offset   int
}

func (p *StockBasicParams) ToMap() map[string]string {
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
	if p.Status != "" {
		m["list_status"] = p.Status
	}
	if p.Limit > 0 {
		m["limit"] = intToStr(p.Limit)
	}
	if p.Offset > 0 {
		m["offset"] = intToStr(p.Offset)
	}
	return m
}

type StockBasicRow struct {
	TSCode     string
	Symbol     string
	Name       string
	Area       string
	Industry   string
	Market     string
	Exchange   string
	ListStatus string
	ListDate   string
	DelistDate string
}

func (c *Client) GetStockBasic(ctx context.Context, params *StockBasicParams) ([]StockBasicRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIStockBasic, params.ToMap(), FieldsStockBasic)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]StockBasicRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, StockBasicRow{
			TSCode:     getStr(idx, item, "ts_code"),
			Symbol:     getStr(idx, item, "symbol"),
			Name:       getStr(idx, item, "name"),
			Area:       getStr(idx, item, "area"),
			Industry:   getStr(idx, item, "industry"),
			Market:     getStr(idx, item, "market"),
			Exchange:   getStr(idx, item, "exchange"),
			ListStatus: getStr(idx, item, "list_status"),
			ListDate:   getStr(idx, item, "list_date"),
			DelistDate: getStr(idx, item, "delist_date"),
		})
	}

	return rows, record, nil
}
