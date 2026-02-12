package eastmoney

import (
	"context"
	"strings"

	"github.com/souloss/quantds/clients/eastmoney"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/announcement"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// AnnouncementAdapter adapts Eastmoney announcement data
type AnnouncementAdapter struct {
	client *eastmoney.Client
}

// NewAnnouncementAdapter creates a new announcement adapter
func NewAnnouncementAdapter(client *eastmoney.Client) *AnnouncementAdapter {
	return &AnnouncementAdapter{client: client}
}

// Name returns the adapter name
func (a *AnnouncementAdapter) Name() string {
	return Name
}

// SupportedMarkets returns supported markets
func (a *AnnouncementAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

// CanHandle checks if the adapter can handle the symbol
func (a *AnnouncementAdapter) CanHandle(symbol string) bool {
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

// Fetch retrieves announcement data
func (a *AnnouncementAdapter) Fetch(ctx context.Context, _ request.Client, req announcement.Request) (announcement.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	params := &eastmoney.AnnouncementParams{
		StockList: req.Symbol,
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
	}

	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageIndex <= 0 {
		params.PageIndex = 1
	}

	result, record, err := a.client.GetAnnouncements(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return announcement.Response{}, trace, err
	}

	announcements := make([]announcement.Announcement, 0, len(result.Data))
	for _, row := range result.Data {
		itemType := announcement.TypeNews
		if strings.Contains(row.ColumnName, "报告") {
			itemType = announcement.TypeReport
		}

		announcements = append(announcements, announcement.Announcement{
			Title:       row.Title,
			PublishTime: row.DisplayTime,
			Source:      row.ColumnName,
			Type:        itemType,
			Code:        row.StockCode,
			Name:        row.StockName,
		})
	}

	hasMore := len(result.Data) >= params.PageSize

	trace.Finish()
	return announcement.Response{
		Symbol:     req.Symbol,
		Data:       announcements,
		Source:     Name,
		HasMore:    hasMore,
		TotalCount: result.Total,
	}, trace, nil
}

var _ manager.Provider[announcement.Request, announcement.Response] = (*AnnouncementAdapter)(nil)
