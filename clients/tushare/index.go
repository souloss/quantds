package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

const (
	APIIndexBasic = "index_basic"
	APIIndexDaily = "index_daily"
)

const (
	FieldsIndexBasic = "ts_code,name,fullname,market,publisher,index_type,category,base_date,base_point,list_date,weight_rule,desc,exp_date"
	FieldsIndexDaily = "ts_code,trade_date,close,open,high,low,pre_close,change,pct_chg,vol,amount"
)

type IndexBasicParams struct {
	Market string // MSCI, CSI, SSE, SZSE, CICC, SW, CNI, OTH
}

type IndexBasicRow struct {
	TSCode    string
	Name      string
	Fullname  string
	Market    string
	Publisher string
	IndexType string
	Category  string
	BaseDate  string
	ListDate  string
}

func (c *Client) GetIndexBasic(ctx context.Context, params *IndexBasicParams) ([]IndexBasicRow, *request.Record, error) {
	m := make(map[string]string)
	if params != nil && params.Market != "" {
		m["market"] = params.Market
	}

	data, record, err := c.post(ctx, APIIndexBasic, m, FieldsIndexBasic)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]IndexBasicRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, IndexBasicRow{
			TSCode:    getStr(idx, item, "ts_code"),
			Name:      getStr(idx, item, "name"),
			Fullname:  getStr(idx, item, "fullname"),
			Market:    getStr(idx, item, "market"),
			Publisher: getStr(idx, item, "publisher"),
			IndexType: getStr(idx, item, "index_type"),
			Category:  getStr(idx, item, "category"),
			BaseDate:  getStr(idx, item, "base_date"),
			ListDate:  getStr(idx, item, "list_date"),
		})
	}

	return rows, record, nil
}

type IndexDailyParams struct {
	TSCode    string
	StartDate string
	EndDate   string
}

type IndexDailyRow struct {
	TSCode    string
	TradeDate string
	Open      float64
	High      float64
	Low       float64
	Close     float64
	PreClose  float64
	Change    float64
	PctChg    float64
	Vol       float64
	Amount    float64
}

func (c *Client) GetIndexDaily(ctx context.Context, params *IndexDailyParams) ([]IndexDailyRow, *request.Record, error) {
	m := make(map[string]string)
	if params.TSCode != "" {
		m["ts_code"] = params.TSCode
	}
	if params.StartDate != "" {
		m["start_date"] = params.StartDate
	}
	if params.EndDate != "" {
		m["end_date"] = params.EndDate
	}

	data, record, err := c.post(ctx, APIIndexDaily, m, FieldsIndexDaily)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]IndexDailyRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, IndexDailyRow{
			TSCode:    getStr(idx, item, "ts_code"),
			TradeDate: getStr(idx, item, "trade_date"),
			Open:      getFlt(idx, item, "open"),
			High:      getFlt(idx, item, "high"),
			Low:       getFlt(idx, item, "low"),
			Close:     getFlt(idx, item, "close"),
			PreClose:  getFlt(idx, item, "pre_close"),
			Change:    getFlt(idx, item, "change"),
			PctChg:    getFlt(idx, item, "pct_chg"),
			Vol:       getFlt(idx, item, "vol"),
			Amount:    getFlt(idx, item, "amount"),
		})
	}

	return rows, record, nil
}
