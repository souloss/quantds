package cninfo

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/souloss/quantds/request"
)

type StockListRow struct {
	Code     string `json:"code"`
	Pinyin   string `json:"pinyin"`
	Category string `json:"category"`
	OrgID    string `json:"orgId"`
	Name     string `json:"zwjc"`
}

type stockListResponse struct {
	StockList []StockListRow `json:"stockList"`
}

func (c *Client) GetStockList(ctx context.Context) ([]StockListRow, *request.Record, error) {
	var data stockListResponse
	record, err := c.do(ctx, "GET", "/new/data/szse_stock.json", nil, &data)
	if err != nil {
		return nil, record, err
	}
	return data.StockList, record, nil
}

type OrgIDParams struct {
	KeyWord string
	MaxNum  int
}

func (p *OrgIDParams) ToValues() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.KeyWord != "" {
		v.Set("keyWord", p.KeyWord)
	}
	if p.MaxNum > 0 {
		v.Set("maxNum", strconv.Itoa(p.MaxNum))
	}
	return v
}

type OrgIDRow struct {
	Code  string `json:"code"`
	OrgID string `json:"orgId"`
}

func (c *Client) GetOrgID(ctx context.Context, params *OrgIDParams) ([]OrgIDRow, *request.Record, error) {
	if params == nil || params.KeyWord == "" {
		return nil, nil, fmt.Errorf("keyWord required")
	}

	query := params.ToValues()
	path := "/new/information/topSearch/query?" + query.Encode()

	var rows []OrgIDRow
	record, err := c.do(ctx, "POST", path, nil, &rows)
	if err != nil {
		return nil, record, err
	}
	return rows, record, nil
}

func (c *Client) GetOrgIDForCode(ctx context.Context, code string) (string, *request.Record, error) {
	rows, record, err := c.GetOrgID(ctx, &OrgIDParams{KeyWord: code, MaxNum: 5})
	if err != nil {
		return "", record, err
	}
	for _, r := range rows {
		if r.Code == code {
			return r.OrgID, record, nil
		}
	}
	if len(rows) > 0 {
		return rows[0].OrgID, record, nil
	}
	return "", record, fmt.Errorf("orgId not found for code %s", code)
}
