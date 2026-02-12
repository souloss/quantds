package tushare

import (
	"context"
	"fmt"
	"time"

	"github.com/souloss/quantds/request"
)

type DailyParams struct {
	TSCode    string
	TradeDate string
	StartDate string
	EndDate   string
}

func (p *DailyParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.TSCode != "" {
		m["ts_code"] = p.TSCode
	}
	if p.TradeDate != "" {
		m["trade_date"] = p.TradeDate
	}
	if p.StartDate != "" {
		m["start_date"] = p.StartDate
	}
	if p.EndDate != "" {
		m["end_date"] = p.EndDate
	}
	return m
}

type KlineParams struct {
	Symbol    string
	StartDate string
	EndDate   string
	Period    string
	Adjust    string
}

type KlineResult struct {
	Data  []KlineBar
	Count int
}

type KlineBar struct {
	TSCode     string
	Date       string
	Open       float64
	High       float64
	Low        float64
	Close      float64
	PreClose   float64
	Change     float64
	ChangeRate float64
	Volume     float64
	Amount     float64
}

func (c *Client) GetKline(ctx context.Context, params *KlineParams) (*KlineResult, *request.Record, error) {
	if params == nil {
		return nil, nil, fmt.Errorf("params required")
	}

	tsCode, err := ToTushareSymbol(params.Symbol)
	if err != nil {
		return nil, nil, err
	}

	apiName := APIDaily
	switch params.Period {
	case "weekly":
		apiName = APIWeekly
	case "monthly":
		apiName = APIMonthly
	}

	dailyParams := &DailyParams{
		TSCode:    tsCode,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	}

	data, record, err := c.post(ctx, apiName, dailyParams.ToMap(), FieldsDaily)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	bars := make([]KlineBar, 0, len(data.Items))
	for _, item := range data.Items {
		bars = append(bars, KlineBar{
			TSCode:     getStr(idx, item, "ts_code"),
			Date:       formatDate(getStr(idx, item, "trade_date")),
			Open:       getFlt(idx, item, "open"),
			High:       getFlt(idx, item, "high"),
			Low:        getFlt(idx, item, "low"),
			Close:      getFlt(idx, item, "close"),
			PreClose:   getFlt(idx, item, "pre_close"),
			Change:     getFlt(idx, item, "change"),
			ChangeRate: getFlt(idx, item, "pct_chg"),
			Volume:     getFlt(idx, item, "vol"),
			Amount:     getFlt(idx, item, "amount"),
		})
	}

	return &KlineResult{Data: bars, Count: len(bars)}, record, nil
}

func ToPeriod(tf string) string {
	switch tf {
	case "1w":
		return "weekly"
	case "1M":
		return "monthly"
	default:
		return "daily"
	}
}

func ParseDate(dateStr string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02", dateStr, timeLoc)
	return t
}

var timeLoc, _ = time.LoadLocation("Asia/Shanghai")

func formatDate(yyyymmdd string) string {
	if len(yyyymmdd) != 8 {
		return yyyymmdd
	}
	return yyyymmdd[:4] + "-" + yyyymmdd[4:6] + "-" + yyyymmdd[6:]
}
