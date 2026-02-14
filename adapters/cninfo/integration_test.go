package cninfo

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/souloss/quantds/clients/cninfo"
	"github.com/souloss/quantds/domain/announcement"
	"github.com/souloss/quantds/domain/instrument"
	"github.com/souloss/quantds/request"
)

// skipOnAPIError 当遇到已知 API 服务错误时跳过测试。
func skipOnAPIError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	msg := err.Error()
	if strings.Contains(msg, "client error") ||
		strings.Contains(msg, "authentication") ||
		strings.Contains(msg, "status 4") ||
		strings.Contains(msg, "status 5") {
		t.Skipf("Skipping due to API error: %v", err)
	}
}

func TestIntegration_InstrumentAdapter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := cninfo.NewClient(cninfo.WithHTTPClient(request.NewClient(request.DefaultConfig())))
	defer client.Close()

	adapter := NewInstrumentAdapter(client)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test fetching all instruments
	resp, trace, err := adapter.Fetch(ctx, nil, instrument.Request{
		PageSize: 100,
	})

	if err != nil {
		skipOnAPIError(t, err)
		t.Fatalf("Fetch failed: %v", err)
	}

	if trace == nil {
		t.Error("Expected non-nil trace")
	}

	if resp.Source != Name {
		t.Errorf("Expected source '%s', got '%s'", Name, resp.Source)
	}

	if len(resp.Data) == 0 {
		t.Error("Expected non-empty instrument list")
	}

	t.Logf("Fetched %d instruments", len(resp.Data))

	// Verify first few entries have required fields
	for i, inst := range resp.Data[:min(3, len(resp.Data))] {
		if inst.Code == "" {
			t.Errorf("Instrument[%d] has empty code", i)
		}
		if inst.Name == "" {
			t.Errorf("Instrument[%d] has empty name", i)
		}
		if inst.Symbol == "" {
			t.Errorf("Instrument[%d] has empty symbol", i)
		}
		t.Logf("Instrument[%d]: %s - %s (%s)", i, inst.Code, inst.Name, inst.Exchange)
	}
}

func TestIntegration_AnnouncementAdapter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := cninfo.NewClient(cninfo.WithHTTPClient(request.NewClient(request.DefaultConfig())))
	defer client.Close()

	adapter := NewAnnouncementAdapter(client)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test fetching announcements for a specific stock
	startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	resp2, trace2, err := adapter.Fetch(ctx, nil, announcement.Request{
		Symbol:    "000001.SZ",
		PageSize:  10,
		PageIndex: 1,
		StartTime: &startTime,
		EndTime:   &endTime,
	})

	if err != nil {
		skipOnAPIError(t, err)
		t.Fatalf("Fetch failed: %v", err)
	}

	if trace2 == nil {
		t.Error("Expected non-nil trace")
	}

	if resp2.Source != Name {
		t.Errorf("Expected source '%s', got '%s'", Name, resp2.Source)
	}

	if resp2.Symbol != "000001.SZ" {
		t.Errorf("Expected symbol '000001.SZ', got '%s'", resp2.Symbol)
	}

	t.Logf("Fetched %d announcements, total: %d, hasMore: %v", 
		len(resp2.Data), resp2.TotalCount, resp2.HasMore)

	// Verify first few entries
	for i, ann := range resp2.Data[:min(3, len(resp2.Data))] {
		if ann.Title == "" {
			t.Errorf("Announcement[%d] has empty title", i)
		}
		t.Logf("Announcement[%d]: %s - %s", i, ann.Code, ann.Title)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
