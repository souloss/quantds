package tushare

import (
	"context"
	"fmt"
	"time"

	"github.com/souloss/quantds/request"
)

type TradeCalParams struct {
	Exchange  string
	StartDate string
	EndDate   string
	IsOpen    string
}

func (p *TradeCalParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.Exchange != "" {
		m["exchange"] = p.Exchange
	}
	if p.StartDate != "" {
		m["start_date"] = p.StartDate
	}
	if p.EndDate != "" {
		m["end_date"] = p.EndDate
	}
	if p.IsOpen != "" {
		m["is_open"] = p.IsOpen
	}
	return m
}

type TradeCalRow struct {
	Exchange     string
	CalDate      string
	IsOpen       bool
	PretradeDate string
}

func (c *Client) GetTradeCal(ctx context.Context, params *TradeCalParams) ([]TradeCalRow, *request.Record, error) {
	data, record, err := c.post(ctx, APITradeCal, params.ToMap(), FieldsTradeCal)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]TradeCalRow, 0, len(data.Items))
	for _, item := range data.Items {
		isOpen := getStr(idx, item, "is_open") == "1"
		rows = append(rows, TradeCalRow{
			Exchange:     getStr(idx, item, "exchange"),
			CalDate:      getStr(idx, item, "cal_date"),
			IsOpen:       isOpen,
			PretradeDate: getStr(idx, item, "pretrade_date"),
		})
	}

	return rows, record, nil
}

func (c *Client) GetLatestTradeDate(ctx context.Context) (string, error) {
	now := time.Now()
	params := &TradeCalParams{
		Exchange:  "SSE",
		StartDate: now.AddDate(0, 0, -10).Format("20060102"),
		EndDate:   now.Format("20060102"),
		IsOpen:    "1",
	}

	rows, _, err := c.GetTradeCal(ctx, params)
	if err != nil {
		return "", err
	}

	if len(rows) == 0 {
		return "", fmt.Errorf("no trading dates found")
	}

	return rows[len(rows)-1].CalDate, nil
}
