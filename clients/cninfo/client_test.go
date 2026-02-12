package cninfo

import (
	"context"
	"testing"
	"time"

	"github.com/souloss/quantds/request"
)

func TestClient_GetStockList(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetStockList(ctx)
	if err != nil {
		t.Fatalf("GetStockList() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	if len(rows) == 0 {
		t.Fatal("no rows returned")
	}

	t.Logf("Total stocks: %d", len(rows))

	for i, row := range rows[:min(3, len(rows))] {
		t.Logf("Row[%d]: code=%s, name=%s, orgId=%s", i, row.Code, row.Name, row.OrgID)
	}
}

func TestClient_GetOrgID(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, record, err := client.GetOrgID(ctx, &OrgIDParams{
		KeyWord: "000001",
		MaxNum:  5,
	})
	if err != nil {
		t.Fatalf("GetOrgID() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	if len(rows) == 0 {
		t.Fatal("no rows returned")
	}

	t.Logf("Found %d org IDs for keyword '000001'", len(rows))
	for i, row := range rows {
		t.Logf("Row[%d]: code=%s, orgId=%s", i, row.Code, row.OrgID)
	}
}

func TestClient_GetOrgIDForCode(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	orgID, record, err := client.GetOrgIDForCode(ctx, "000001")
	if err != nil {
		t.Fatalf("GetOrgIDForCode() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	if orgID == "" {
		t.Fatal("orgID is empty")
	}

	t.Logf("000001 -> orgId: %s", orgID)
}

func TestClient_QueryNews(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	orgID, _, err := client.GetOrgIDForCode(ctx, "000001")
	if err != nil {
		t.Fatalf("GetOrgIDForCode() error = %v", err)
	}

	rows, total, record, err := client.QueryNews(ctx, &NewsQueryParams{
		PageNum:  1,
		PageSize: 5,
		Stock:    "000001," + orgID,
		SeDate:   "2024-01-01~2025-12-31",
		TabName:  "fulltext",
	})
	if err != nil {
		t.Fatalf("QueryNews() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Total: %d, Returned: %d", total, len(rows))

	if len(rows) > 0 {
		for i, row := range rows[:min(3, len(rows))] {
			t.Logf("Row[%d]: title=%s, time=%d", i, row.AnnouncementTitle, row.AnnouncementTime)
		}
	}
}

func TestClient_QueryNewsByColumn(t *testing.T) {
	client := NewClient(request.NewClient(request.DefaultConfig()))
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, total, record, err := client.QueryNewsByColumn(ctx, &NewsQueryByColumnParams{
		Column:   "szse",
		PageNum:  1,
		PageSize: 5,
	})
	if err != nil {
		t.Fatalf("QueryNewsByColumn() error = %v", err)
	}

	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Column: szse, Total: %d, Returned: %d", total, len(rows))

	if len(rows) > 0 {
		for i, row := range rows[:min(3, len(rows))] {
			t.Logf("Row[%d]: code=%s, title=%s", i, row.SecCode, row.AnnouncementTitle)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
