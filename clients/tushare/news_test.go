package tushare

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetNews tests retrieving news data
// API Rule: Requires Tushare token
func TestClient_GetNews(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get news from the last 24 hours
	rows, record, err := client.GetNews(ctx, &NewsParams{
		Src: "sina",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetNews() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d news rows", len(rows))
	for i, r := range rows {
		if i >= 5 {
			break
		}
		t.Logf("News[%d]: datetime=%s, title=%s", i, r.Datetime, truncate(r.Title, 50))
	}
}

// TestClient_GetCctvNews tests retrieving CCTV news
func TestClient_GetCctvNews(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Get CCTV news for recent date
	rows, record, err := client.GetCctvNews(ctx, &CctvNewsParams{
		Date: "20240101",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetCctvNews() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d CCTV news rows", len(rows))
	for i, r := range rows {
		if i >= 3 {
			break
		}
		t.Logf("CCTV[%d]: date=%s, title=%s", i, r.Date, truncate(r.Title, 40))
	}
}

// TestClient_GetAnnouncements tests retrieving company announcements
func TestClient_GetAnnouncements(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := context.Background()

	rows, record, err := client.GetAnnouncements(ctx, &AnnouncementParams{
		TsCode:    "000001.SZ",
		StartDate: "20240101",
		EndDate:   "20240131",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetAnnouncements() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d announcement rows for 000001.SZ", len(rows))
	for i, r := range rows {
		if i >= 5 {
			break
		}
		t.Logf("Announcement[%d]: date=%s, title=%s", i, r.AnnDate, truncate(r.Title, 40))
	}
}

// Helper function to truncate string
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
