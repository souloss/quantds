package szse

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"github.com/souloss/quantds/request"
	"github.com/xuri/excelize/v2"
)

type StockListParams struct {
	ShowType  string
	CatalogID string
	TabKey    string
	Random    string
}

func (p *StockListParams) ToQuery() map[string]string {
	if p == nil {
		return nil
	}
	return map[string]string{
		"SHOWTYPE":  p.ShowType,
		"CATALOGID": p.CatalogID,
		"TABKEY":    p.TabKey,
		"random":    p.Random,
	}
}

type StockListResult struct {
	Data [][]string
}

func (c *Client) GetStockList(ctx context.Context, params *StockListParams) (*StockListResult, *request.Record, error) {
	if params == nil {
		params = &StockListParams{}
	}
	if params.CatalogID == "" {
		params.CatalogID = "1110"
	}
	if params.TabKey == "" {
		params.TabKey = "tab1"
	}
	if params.ShowType == "" {
		params.ShowType = "xlsx"
	}

	q := make(url.Values)
	for k, v := range params.ToQuery() {
		if v != "" {
			q.Set(k, v)
		}
	}

	req := request.Request{
		Method:  "GET",
		URL:     c.baseURL + "?" + q.Encode(),
		Headers: c.headers,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	f, err := excelize.OpenReader(bytes.NewReader(resp.Body))
	if err != nil {
		return nil, record, fmt.Errorf("open excel: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, record, fmt.Errorf("get rows: %w", err)
	}

	return &StockListResult{Data: rows}, record, nil
}
