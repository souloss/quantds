package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/souloss/quantds/request"
)

type StockListParams struct {
	StockType     string
	RegProvince   string
	CsrcCode      string
	StockCode     string
	SqlID         string
	CompanyStatus string
	Type          string
	IsPagination  string
	CacheSize     string
	BeginPage     string
	PageSize      string
	PageNo        string
	EndPage       string
}

func (p *StockListParams) ToQuery() map[string]string {
	if p == nil {
		return nil
	}
	return map[string]string{
		"STOCK_TYPE":         p.StockType,
		"REG_PROVINCE":       p.RegProvince,
		"CSRC_CODE":          p.CsrcCode,
		"STOCK_CODE":         p.StockCode,
		"sqlId":              p.SqlID,
		"COMPANY_STATUS":     p.CompanyStatus,
		"type":               p.Type,
		"isPagination":       p.IsPagination,
		"pageHelp.cacheSize": p.CacheSize,
		"pageHelp.beginPage": p.BeginPage,
		"pageHelp.pageSize":  p.PageSize,
		"pageHelp.pageNo":    p.PageNo,
		"pageHelp.endPage":   p.EndPage,
	}
}

type StockRow struct {
	CompanyCode string `json:"COMPANY_CODE"`
	CompanyAbbr string `json:"COMPANY_ABBR"`
	SecNameCn   string `json:"SEC_NAME_CN"`
	ListDate    string `json:"LIST_DATE"`
	TotalShares string `json:"TOTAL_SHARES"`
	FloatShares string `json:"FLOW_SHARES"`
	Industry    string `json:"INDUSTRY_NAME"`
}

type StockListResult struct {
	Data  []StockRow
	Total int
}

func (c *Client) GetStockList(ctx context.Context, params *StockListParams) (*StockListResult, *request.Record, error) {
	if params == nil {
		params = &StockListParams{}
	}
	if params.SqlID == "" {
		params.SqlID = SqlID
	}
	if params.CompanyStatus == "" {
		params.CompanyStatus = Status
	}
	if params.Type == "" {
		params.Type = "inParams"
	}
	if params.IsPagination == "" {
		params.IsPagination = "true"
	}
	if params.CacheSize == "" {
		params.CacheSize = "1"
	}
	if params.BeginPage == "" {
		params.BeginPage = "1"
	}
	if params.PageSize == "" {
		params.PageSize = "10000"
	}
	if params.PageNo == "" {
		params.PageNo = "1"
	}
	if params.EndPage == "" {
		params.EndPage = "1"
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

	if strings.Contains(string(resp.Body), "<html") {
		return nil, record, fmt.Errorf("unexpected html response")
	}

	var result struct {
		Result []StockRow `json:"result"`
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, fmt.Errorf("unmarshal: %w", err)
	}

	return &StockListResult{Data: result.Result, Total: len(result.Result)}, record, nil
}
