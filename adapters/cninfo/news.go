package cninfo

import (
	"context"
	"fmt"
	"time"

	"github.com/souloss/quantds/clients/cninfo"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/announcement"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// AnnouncementAdapter adapts CNINFO announcement/news data.
type AnnouncementAdapter struct {
	client *cninfo.Client
}

// NewAnnouncementAdapter creates a new announcement adapter.
func NewAnnouncementAdapter(client *cninfo.Client) *AnnouncementAdapter {
	return &AnnouncementAdapter{client: client}
}

func (a *AnnouncementAdapter) Name() string                      { return Name }
func (a *AnnouncementAdapter) SupportedMarkets() []domain.Market { return supportedMarkets }

func (a *AnnouncementAdapter) CanHandle(symbol string) bool {
	if symbol == "" {
		return true // Support querying all announcements
	}
	var sym domain.Symbol
	if err := sym.Parse(symbol); err != nil {
		return false
	}
	for _, m := range supportedMarkets {
		if sym.Market == m {
			return true
		}
	}
	return false
}

func (a *AnnouncementAdapter) Fetch(ctx context.Context, _ request.Client, req announcement.Request) (announcement.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	pageNum := req.PageIndex
	if pageNum <= 0 {
		pageNum = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// Build stock parameter (code + orgId)
	stock := ""
	if req.Symbol != "" {
		code := req.Symbol
		// Extract code from "000001.SZ" format
		if idx := len(code) - 3; idx > 0 && code[idx] == '.' {
			code = code[:idx]
		}
		orgID, _, _ := a.client.GetOrgIDForCode(ctx, code)
		if orgID != "" {
			stock = code + "," + orgID
		} else {
			stock = code
		}
	}

	// Build date range filter
	seDate := ""
	if req.StartTime != nil && req.EndTime != nil {
		seDate = req.StartTime.Format("2006-01-02") + "~" + req.EndTime.Format("2006-01-02")
	}

	params := &cninfo.NewsQueryParams{
		PageNum:  pageNum,
		PageSize: pageSize,
		Stock:    stock,
		SeDate:   seDate,
		TabName:  "fulltext",
	}

	rows, total, record, err := a.client.QueryNews(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return announcement.Response{}, trace, err
	}

	items := make([]announcement.Announcement, 0, len(rows))
	for _, row := range rows {
		publishTime := ""
		if row.AnnouncementTime > 0 {
			publishTime = time.Unix(row.AnnouncementTime/1000, 0).Format("2006-01-02 15:04:05")
		}
		url := ""
		if row.AdjunctURL != "" {
			url = fmt.Sprintf("http://static.cninfo.com.cn/%s", row.AdjunctURL)
		}
		items = append(items, announcement.Announcement{
			Title:       row.AnnouncementTitle,
			PublishTime: publishTime,
			Source:      "cninfo",
			URL:         url,
			Type:        announcement.TypeNews,
			Code:        row.SecCode,
			Name:        row.SecName,
		})
	}

	hasMore := len(rows) >= pageSize

	trace.Finish()
	return announcement.Response{
		Symbol:     req.Symbol,
		Data:       items,
		Source:     Name,
		HasMore:    hasMore,
		TotalCount: total,
	}, trace, nil
}

var _ manager.Provider[announcement.Request, announcement.Response] = (*AnnouncementAdapter)(nil)
