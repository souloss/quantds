package tushare

import (
	"context"
	"fmt"
	"time"

	"github.com/souloss/quantds/clients/tushare"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/announcement"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

type AnnouncementAdapter struct {
	client *tushare.Client
}

func NewAnnouncementAdapter(client *tushare.Client) *AnnouncementAdapter {
	return &AnnouncementAdapter{client: client}
}

func (a *AnnouncementAdapter) Name() string {
	return Name
}

func (a *AnnouncementAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

func (a *AnnouncementAdapter) CanHandle(symbol string) bool {
	if symbol == "" {
		return true // 支持按市场查询
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

	params := &tushare.AnnouncementParams{
		TsCode: req.Symbol,
	}
	if req.StartTime != nil {
		params.StartDate = req.StartTime.Format("20060102")
	}
	if req.EndTime != nil {
		params.EndDate = req.EndTime.Format("20060102")
	}

	rows, record, err := a.client.GetAnnouncements(ctx, params)
	trace.AddRequest(record)
	if err != nil {
		return announcement.Response{}, trace, err
	}

	anns := make([]announcement.Announcement, 0, len(rows))
	for _, row := range rows {
		// AnnDate format YYYYMMDD
		pubTime := row.AnnDate
		if t, err := time.Parse("20060102", row.AnnDate); err == nil {
			pubTime = t.Format("2006-01-02")
		}

		anns = append(anns, announcement.Announcement{
			ID:          fmt.Sprintf("%s-%s", row.TsCode, row.AnnDate), // 简易ID生成
			Title:       row.Title,
			Content:     row.Content,
			PublishTime: pubTime,
			Source:      Name,
			Code:        row.TsCode,
			Type:        announcement.TypeReport, // 默认类型，Tushare disclosure 主要是公告
			Category:    announcement.CategoryCompany,
		})
	}

	trace.Finish()
	return announcement.Response{
		Symbol:     req.Symbol,
		Data:       anns,
		Source:     Name,
		TotalCount: len(anns),
	}, trace, nil
}

var _ manager.Provider[announcement.Request, announcement.Response] = (*AnnouncementAdapter)(nil)
