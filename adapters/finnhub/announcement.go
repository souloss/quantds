package finnhub

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/finnhub"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/announcement"
	"github.com/souloss/quantds/manager"
)

// AnnouncementAdapter adapts Finnhub company news for US stocks
type AnnouncementAdapter struct {
	client *finnhub.Client
}

// NewAnnouncementAdapter creates a new announcement adapter
func NewAnnouncementAdapter(client *finnhub.Client) *AnnouncementAdapter {
	return &AnnouncementAdapter{client: client}
}

func (a *AnnouncementAdapter) Name() string                      { return Name }
func (a *AnnouncementAdapter) SupportedMarkets() []domain.Market { return usMarkets }

func (a *AnnouncementAdapter) CanHandle(symbol string) bool {
	var sym domain.Symbol
	if err := sym.Parse(symbol); err != nil {
		return false
	}
	for _, m := range usMarkets {
		if sym.Market == m {
			return true
		}
	}
	return false
}

func (a *AnnouncementAdapter) Fetch(ctx context.Context, req announcement.Request) (announcement.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	var sym domain.Symbol
	if err := sym.Parse(req.Symbol); err != nil {
		return announcement.Response{}, trace, err
	}

	// Calculate date range
	from := req.StartTime
	to := req.EndTime
	if to.IsZero() {
		to = time.Now()
	}
	if from.IsZero() {
		from = to.AddDate(0, 0, -7) // Default: last 7 days
	}

	result, record, err := a.client.GetCompanyNews(ctx, &finnhub.NewsParams{
		Symbol: sym.Code,
		From:   from.Format("2006-01-02"),
		To:     to.Format("2006-01-02"),
	})
	trace.AddRequest(record)

	if err != nil {
		return announcement.Response{}, trace, err
	}

	// Limit results based on PageSize
	limit := len(result.News)
	if req.PageSize > 0 && limit > req.PageSize {
		limit = req.PageSize
	}

	items := make([]announcement.Announcement, 0, limit)
	for i := 0; i < limit; i++ {
		n := result.News[i]
		items = append(items, announcement.Announcement{
			ID:          string(rune(n.ID)),
			Title:       n.Headline,
			Content:     n.Summary,
			Summary:     n.Summary,
			PublishTime: n.Datetime.Format("2006-01-02 15:04"),
			Source:      n.Source,
			URL:         n.URL,
			Type:        announcement.TypeNews,
			Category:    announcement.CategoryMedia,
			Code:        sym.Code,
			Name:        sym.Code,
		})
	}

	trace.Finish()
	return announcement.Response{
		Symbol:      req.Symbol,
		Data:        items,
		Source:      Name,
		HasMore:     len(result.News) > limit,
		TotalCount:  len(result.News),
		PageSize:    req.PageSize,
		DataVersion: 1,
	}, trace, nil
}

var _ manager.Provider[announcement.Request, announcement.Response] = (*AnnouncementAdapter)(nil)
