package cninfo

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/souloss/quantds/request"
)

type NewsQueryParams struct {
	PageNum     int
	PageSize    int
	TabName     string
	IsHLTitle   bool
	SeDate      string
	Stock       string
	Column      string
	Plate       string
	ColumnTitle string
	Searchkey   string
	Secid       string
	Category    string
	Trade       string
	SortName    string
	SortType    string
}

func (p *NewsQueryParams) ToValues() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.PageNum > 0 {
		v.Set("pageNum", strconv.Itoa(p.PageNum))
	}
	if p.PageSize > 0 {
		v.Set("pageSize", strconv.Itoa(p.PageSize))
	}
	if p.TabName != "" {
		v.Set("tabName", p.TabName)
	}
	if p.IsHLTitle {
		v.Set("isHLtitle", "true")
	} else {
		v.Set("isHLtitle", "false")
	}
	if p.SeDate != "" {
		v.Set("seDate", p.SeDate)
	}
	v.Set("stock", p.Stock)
	v.Set("column", p.Column)
	v.Set("plate", p.Plate)
	v.Set("columnTitle", p.ColumnTitle)
	v.Set("searchkey", p.Searchkey)
	v.Set("secid", p.Secid)
	v.Set("category", p.Category)
	v.Set("trade", p.Trade)
	v.Set("sortName", p.SortName)
	v.Set("sortType", p.SortType)
	return v
}

type AnnouncementRow struct {
	ID                string `json:"id"`
	SecCode           string `json:"secCode"`
	SecName           string `json:"secName"`
	OrgID             string `json:"orgId"`
	AnnouncementID    string `json:"announcementId"`
	AnnouncementTitle string `json:"announcementTitle"`
	AnnouncementTime  int64  `json:"announcementTime"`
	AdjunctURL        string `json:"adjunctUrl"`
	ColumnID          string `json:"columnId"`
}

type announcementResponse struct {
	Announcements  []AnnouncementRow `json:"announcements"`
	TotalRecordNum int               `json:"totalRecordNum"`
}

func (c *Client) QueryNews(ctx context.Context, params *NewsQueryParams) ([]AnnouncementRow, int, *request.Record, error) {
	if params == nil {
		return nil, 0, nil, fmt.Errorf("params required")
	}

	var data announcementResponse
	record, err := c.do(ctx, "POST", "/new/hisAnnouncement/query", params.ToValues(), &data)
	if err != nil {
		return nil, 0, record, err
	}
	return data.Announcements, data.TotalRecordNum, record, nil
}

type NewsQueryByColumnParams struct {
	Column      string
	Plate       string
	ColumnTitle string
	SeDate      string
	PageNum     int
	PageSize    int
}

func (p *NewsQueryByColumnParams) ToValues() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	if p.PageNum > 0 {
		v.Set("pageNum", strconv.Itoa(p.PageNum))
	}
	if p.PageSize > 0 {
		v.Set("pageSize", strconv.Itoa(p.PageSize))
	}
	v.Set("tabName", "fulltext")
	v.Set("isHLtitle", "true")
	if p.SeDate != "" {
		v.Set("seDate", p.SeDate)
	}
	v.Set("column", p.Column)
	v.Set("plate", p.Plate)
	v.Set("columnTitle", p.ColumnTitle)
	v.Set("stock", "")
	v.Set("searchkey", "")
	v.Set("secid", "")
	v.Set("category", "")
	v.Set("trade", "")
	v.Set("sortName", "")
	v.Set("sortType", "")
	return v
}

func (c *Client) QueryNewsByColumn(ctx context.Context, params *NewsQueryByColumnParams) ([]AnnouncementRow, int, *request.Record, error) {
	if params == nil {
		return nil, 0, nil, fmt.Errorf("params required")
	}

	var data announcementResponse
	record, err := c.do(ctx, "POST", "/new/hisAnnouncement/query", params.ToValues(), &data)
	if err != nil {
		return nil, 0, record, err
	}
	return data.Announcements, data.TotalRecordNum, record, nil
}

func formatDateStr(ts int64) string {
	if ts == 0 {
		return ""
	}
	year := ts / 1000000000
	month := (ts / 100000000) % 100
	day := (ts / 1000000) % 100
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func parseAnnouncementURL(adjunctURL string) string {
	if adjunctURL == "" {
		return ""
	}
	if strings.HasPrefix(adjunctURL, "http") {
		return adjunctURL
	}
	return BaseURL + adjunctURL
}

func GetPDFURL(announcementID string) string {
	return fmt.Sprintf("http://static.cninfo.com.cn/%s.pdf", announcementID)
}
