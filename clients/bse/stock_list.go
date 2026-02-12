package bse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/souloss/quantds/request"
)

type StockListParams struct {
	Page      int
	Typejb    string
	Xxfcbj    string
	Xxzqdm    string
	Sortfield string
	Sorttype  string
}

func (p *StockListParams) ToValues() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	v.Set("page", strconv.Itoa(p.Page))
	v.Set("typejb", p.Typejb)
	v.Set("xxfcbj[]", p.Xxfcbj)
	v.Set("xxzqdm", p.Xxzqdm)
	v.Set("sortfield", p.Sortfield)
	v.Set("sorttype", p.Sorttype)
	return v
}

type StockRow struct {
	StockCode   string  `json:"xxzqdm"`
	StockName   string  `json:"xxzqjc"`
	ListDate    string  `json:"xxssrq"`
	TotalShares float64 `json:"xxzgb"`
	FloatShares float64 `json:"xxltgb"`
	Industry    string  `json:"xxsshy"`
}

type StockListResult struct {
	Data       []StockRow
	TotalPages int
}

func (c *Client) GetStockListPage(ctx context.Context, params *StockListParams) (*StockListResult, *request.Record, error) {
	if params == nil {
		params = &StockListParams{}
	}
	if params.Typejb == "" {
		params.Typejb = "T"
	}
	if params.Xxfcbj == "" {
		params.Xxfcbj = "2"
	}
	if params.Sortfield == "" {
		params.Sortfield = "xxzqdm"
	}
	if params.Sorttype == "" {
		params.Sorttype = "asc"
	}

	form := params.ToValues()
	headers := make(map[string]string)
	for k, v := range c.headers {
		headers[k] = v
	}
	headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"

	req := request.Request{
		Method:  "POST",
		URL:     c.baseURL,
		Headers: headers,
		Body:    []byte(form.Encode()),
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	bodyStr := string(resp.Body)
	startIdx := strings.Index(bodyStr, "[")
	endIdx := strings.LastIndex(bodyStr, "]")
	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return nil, record, fmt.Errorf("invalid response format")
	}

	var result []struct {
		Content    []StockRow `json:"content"`
		TotalPages int        `json:"totalPages"`
	}
	if err := json.Unmarshal([]byte(bodyStr[startIdx:endIdx+1]), &result); err != nil {
		return nil, record, fmt.Errorf("unmarshal: %w", err)
	}
	if len(result) == 0 {
		return nil, record, fmt.Errorf("empty result")
	}

	return &StockListResult{Data: result[0].Content, TotalPages: result[0].TotalPages}, record, nil
}

func (c *Client) GetStockList(ctx context.Context) ([]StockRow, []*request.Record, error) {
	var allRows []StockRow
	var records []*request.Record

	page := 0
	for {
		page++
		result, record, err := c.GetStockListPage(ctx, &StockListParams{Page: page})
		records = append(records, record)
		if err != nil {
			return allRows, records, err
		}
		allRows = append(allRows, result.Data...)
		if page >= result.TotalPages {
			break
		}
	}

	return allRows, records, nil
}
