package tushare

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/souloss/quantds/request"
)

// newTestClient 创建测试用客户端，需要设置 TUSHARE_TOKEN 环境变量。
func newTestClient(t *testing.T) *Client {
	t.Helper()
	token := os.Getenv("TUSHARE_TOKEN")
	if token == "" {
		t.Skip("TUSHARE_TOKEN not set")
	}
	return NewClient(request.NewClient(request.DefaultConfig()), WithToken(token))
}

// skipOnTokenError 当遇到 token 相关错误时跳过测试（而非失败）。
func skipOnTokenError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	msg := err.Error()
	if strings.Contains(msg, "token") || strings.Contains(msg, "40101") || strings.Contains(msg, "-1") {
		t.Skipf("Token issue, skipping: %v", err)
	}
}

func TestClient_GetDailyBasic(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetDailyBasic(ctx, &DailyBasicParams{
		TSCode: "000001.SZ",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetDailyBasic() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d daily_basic rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: date=%s, PE=%.2f, PB=%.2f, PS=%.2f, TotalMV=%.2f万元, TurnoverRate=%.2f%%",
			r.TradeDate, r.PE, r.PB, r.PS, r.TotalMV, r.TurnoverRate)
	}
}

func TestClient_GetAdjFactor(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetAdjFactor(ctx, &AdjFactorParams{
		TSCode:    "000001.SZ",
		StartDate: "20240101",
		EndDate:   "20240131",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetAdjFactor() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d adj_factor rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: date=%s, adj_factor=%.4f", r.TradeDate, r.AdjFactor)
	}
}

func TestClient_GetIncome(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetIncome(ctx, &IncomeParams{
		TSCode:     "000001.SZ",
		ReportType: "1",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetIncome() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d income rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: end_date=%s, revenue=%.2f, net_income=%.2f, EPS=%.4f",
			r.EndDate, r.Revenue, r.NIncome, r.BasicEPS)
	}
}

func TestClient_GetBalanceSheet(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetBalanceSheet(ctx, &BalanceSheetParams{
		TSCode:     "000001.SZ",
		ReportType: "1",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetBalanceSheet() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d balance_sheet rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: end_date=%s, total_assets=%.2f, total_liab=%.2f, equity=%.2f",
			r.EndDate, r.TotalAssets, r.TotalLiab, r.TotalHldrEqy)
	}
}

func TestClient_GetCashflow(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetCashflow(ctx, &CashflowParams{
		TSCode:     "000001.SZ",
		ReportType: "1",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetCashflow() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d cashflow rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: end_date=%s, operating=%.2f, investing=%.2f, financing=%.2f",
			r.EndDate, r.NCashflowAct, r.NCashflowInv, r.NCashflowFnc)
	}
}

func TestClient_GetFinaIndicator(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetFinaIndicator(ctx, &FinaIndicatorParams{
		TSCode: "000001.SZ",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetFinaIndicator() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d fina_indicator rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("First: end_date=%s, ROE=%.2f%%, ROA=%.2f%%, GrossMargin=%.2f%%",
			r.EndDate, r.ROE, r.ROA, r.GrossProfitMargin)
	}
}

func TestClient_GetStockCompany(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetStockCompany(ctx, &StockCompanyParams{
		TSCode: "000001.SZ",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetStockCompany() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d company rows", len(rows))
	if len(rows) > 0 {
		r := rows[0]
		t.Logf("Company: chairman=%s, reg_capital=%.2f万, province=%s, employees=%d",
			r.Chairman, r.RegCapital, r.Province, r.Employees)
	}
}

func TestNewClient_EnvVars(t *testing.T) {
	// 测试环境变量读取
	origURL := os.Getenv("TUSHARE_BASE_URL")
	origToken := os.Getenv("TUSHARE_TOKEN")
	defer func() {
		os.Setenv("TUSHARE_BASE_URL", origURL)
		os.Setenv("TUSHARE_TOKEN", origToken)
	}()

	os.Setenv("TUSHARE_BASE_URL", "http://custom.api.test")
	os.Setenv("TUSHARE_TOKEN", "test_token_123")

	client := NewClient(nil)
	if client.baseURL != "http://custom.api.test" {
		t.Errorf("Expected baseURL from env, got %s", client.baseURL)
	}
	if client.token != "test_token_123" {
		t.Errorf("Expected token from env, got %s", client.token)
	}
}

func TestNewClient_Options(t *testing.T) {
	// 测试 Option 覆盖环境变量
	client := NewClient(nil, WithBaseURL("http://override.test"), WithToken("override_token"))
	if client.baseURL != "http://override.test" {
		t.Errorf("Expected baseURL from option, got %s", client.baseURL)
	}
	if client.token != "override_token" {
		t.Errorf("Expected token from option, got %s", client.token)
	}
}
