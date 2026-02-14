package eastmoney

import (
	"context"
	"testing"
)

// TestClient_GetAnnouncements tests retrieving company announcements
// API Rule: No authentication required
// Geo-Restriction: May be blocked in some regions
func TestClient_GetAnnouncements(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &AnnouncementParams{
		StockList: "000001", // Ping An Bank
		PageIndex: 1,
		PageSize:  5,
	}

	result, record, err := client.GetAnnouncements(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Announcement Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d announcements (total: %d)", len(result.Data), result.Total)

	if len(result.Data) == 0 {
		t.Log("Warning: No announcements returned")
		return
	}

	for i, ann := range result.Data {
		t.Logf("Announcement[%d]: code=%s, name=%s, title=%s, date=%s",
			i, ann.StockCode, ann.StockName, truncate(ann.Title, 50), ann.NoticeDate)
	}
}

// TestClient_GetNews tests the GetNews alias for GetAnnouncements
func TestClient_GetNews(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &NewsParams{
		PageIndex: 1,
		PageSize:  10,
		AnnType:   "SHA,SZA", // Shanghai and Shenzhen A-shares
	}

	result, record, err := client.GetNews(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("News Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d news items (total: %d)", len(result.Data), result.Total)

	for i, news := range result.Data {
		t.Logf("News[%d]: %s - %s", i, news.StockCode, truncate(news.Title, 40))
	}
}

// TestClient_GetAnnouncements_AllMarkets tests announcements from all markets
func TestClient_GetAnnouncements_AllMarkets(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &AnnouncementParams{
		PageIndex: 1,
		PageSize:  5,
		AnnType:   "SHA,SZA,BJA", // All A-share markets
	}

	result, _, err := client.GetAnnouncements(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("All Markets Announcements: %d items", len(result.Data))

	for _, ann := range result.Data {
		t.Logf("  [%s] %s: %s", ann.StockCode, ann.StockName, truncate(ann.Title, 30))
	}
}

// Helper function to truncate string
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
