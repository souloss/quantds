package tushare

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/souloss/quantds/clients/tushare"
	"github.com/souloss/quantds/domain/announcement"
)

func TestAnnouncementAdapter_Fetch(t *testing.T) {
	client := tushare.NewClient()
	adapter := NewAnnouncementAdapter(client)

	t.Run("CanHandle", func(t *testing.T) {
		if !adapter.CanHandle("000001.SZ") {
			t.Error("Should handle SZ stock")
		}
		if !adapter.CanHandle("") {
			t.Error("Should handle empty symbol (market query)")
		}
	})

	t.Run("Fetch", func(t *testing.T) {
		// Only run if TUSHARE_TOKEN is set
		if client.Token() == "" {
			t.Skip("TUSHARE_TOKEN not set, skipping integration test")
		}

		ctx := context.Background()
		now := time.Now()
		startTime := now.AddDate(0, 0, -30)
		
		req := announcement.Request{
			Symbol:    "000001.SZ",
			StartTime: &startTime,
			EndTime:   &now,
		}

		resp, trace, err := adapter.Fetch(ctx, nil, req)
		if err != nil {
			msg := err.Error()
			if strings.Contains(msg, "token") || strings.Contains(msg, "40101") || strings.Contains(msg, "-1") {
				t.Skipf("Token issue, skipping: %v", err)
			}
			t.Fatalf("Fetch failed: %v", err)
		}

		if trace == nil {
			t.Error("Trace should not be nil")
		}

		t.Logf("Found %d announcements", len(resp.Data))
		for i, a := range resp.Data {
			if i >= 5 {
				break
			}
			t.Logf("Announcement: [%s] %s", a.PublishTime, a.Title)
			if a.Code == "" {
				t.Error("Code should not be empty")
			}
		}
	})
}
