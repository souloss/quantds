package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/souloss/quantds/request"
)

const AnnouncementAPI = "https://np-anotice-stock.eastmoney.com/api/security/ann"

// AnnouncementParams represents parameters for announcement request
type AnnouncementParams struct {
	StockList string // Comma-separated stock codes
	PageIndex int    // Page index (starting from 1)
	PageSize  int    // Page size
	AnnType   string // Announcement types (e.g., "SHA,SZA,BJA")
}

// AnnouncementResult represents the announcement result
type AnnouncementResult struct {
	Data  []AnnouncementData
	Total int
}

// AnnouncementData represents a single announcement
type AnnouncementData struct {
	ArtCode     string // Article code
	Title       string // Title
	DisplayTime string // Display time
	NoticeDate  string // Notice date
	StockCode   string // Stock code
	StockName   string // Stock name
	ColumnName  string // Column/category name
}

// NewsParams is an alias for AnnouncementParams
type NewsParams = AnnouncementParams

// NewsResult is an alias for AnnouncementResult
type NewsResult = AnnouncementResult

// NewsRow is an alias for AnnouncementData
type NewsRow = AnnouncementData

// GetAnnouncements retrieves company announcements
func (c *Client) GetAnnouncements(ctx context.Context, params *AnnouncementParams) (*AnnouncementResult, *request.Record, error) {
	v := url.Values{}
	v.Set("cb", "jQuery112300")
	v.Set("sr", "-1")
	v.Set("client", "web")
	v.Set("fnode", "1")
	v.Set("snode", "1")

	if params.StockList != "" {
		v.Set("code", params.StockList)
	}
	if params.PageIndex > 0 {
		v.Set("pageNo", strconv.Itoa(params.PageIndex))
	} else {
		v.Set("pageNo", "1")
	}
	if params.PageSize > 0 {
		v.Set("pageSize", strconv.Itoa(params.PageSize))
	} else {
		v.Set("pageSize", "20")
	}
	if params.AnnType != "" {
		v.Set("ann_type", params.AnnType)
	} else {
		v.Set("ann_type", "SHA,SZA,BJA")
	}

	apiURL := AnnouncementAPI + "?" + v.Encode()

	req := request.Request{
		Method: "GET",
		URL:    apiURL,
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":    "https://data.eastmoney.com/notices/stock.html",
		},
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseAnnouncementResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

// GetNews is an alias for GetAnnouncements
func (c *Client) GetNews(ctx context.Context, params *NewsParams) (*NewsResult, *request.Record, error) {
	return c.GetAnnouncements(ctx, params)
}

type announcementResponse struct {
	Data struct {
		TotalHits int `json:"total_hits"`
		List      []struct {
			ArtCode     string `json:"art_code"`
			Title       string `json:"title"`
			DisplayTime string `json:"display_time"`
			NoticeDate  string `json:"notice_date"`
			Codes       []struct {
				StockCode string `json:"stock_code"`
				ShortName string `json:"short_name"`
			} `json:"codes"`
			Columns []struct {
				ColumnName string `json:"column_name"`
			} `json:"columns"`
		} `json:"list"`
	} `json:"data"`
}

func parseAnnouncementResponse(body []byte) (*AnnouncementResult, error) {
	jsonStr := string(body)
	if idx := strings.Index(jsonStr, "("); idx > 0 && strings.HasSuffix(jsonStr, ")") {
		jsonStr = jsonStr[idx+1 : len(jsonStr)-1]
	}

	var resp announcementResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return nil, err
	}

	rows := make([]AnnouncementData, 0, len(resp.Data.List))
	for _, item := range resp.Data.List {
		row := AnnouncementData{
			ArtCode:     item.ArtCode,
			Title:       item.Title,
			DisplayTime: item.DisplayTime,
			NoticeDate:  item.NoticeDate,
		}
		if len(item.Codes) > 0 {
			row.StockCode = item.Codes[0].StockCode
			row.StockName = item.Codes[0].ShortName
		}
		if len(item.Columns) > 0 {
			row.ColumnName = item.Columns[0].ColumnName
		}
		rows = append(rows, row)
	}

	return &AnnouncementResult{Data: rows, Total: resp.Data.TotalHits}, nil
}
