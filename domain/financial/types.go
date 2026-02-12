// Package financial provides financial report data domain types.
//
// This package defines the request/response types for financial statement data
// including income statements, balance sheets, and cash flow statements.
package financial

import (
	"context"
	"time"
)

// ReportType represents the type of financial report
type ReportType string

const (
	ReportTypeIncome    ReportType = "INCOME"    // 利润表
	ReportTypeBalance   ReportType = "BALANCE"   // 资产负债表
	ReportTypeCashflow  ReportType = "CASHFLOW"  // 现金流量表
	ReportTypeIndicator ReportType = "INDICATOR" // 财务指标
)

// ReportPeriod represents the reporting period
type ReportPeriod string

const (
	PeriodAnnual     ReportPeriod = "ANNUAL"    // 年报
	PeriodSemiAnnual ReportPeriod = "SEMI"      // 半年报
	PeriodQuarterly  ReportPeriod = "QUARTERLY" // 季报
)

// Request represents a financial data request
type Request struct {
	Symbol     string       // Symbol code
	ReportType ReportType   // Type of report
	Period     ReportPeriod // Reporting period
	StartDate  string       // Start date (YYYY-MM-DD)
	EndDate    string       // End date (YYYY-MM-DD)
	PageSize   int          // Page size
	PageNumber int          // Page number
}

// CacheKey returns the cache key for the request
func (r Request) CacheKey() string {
	return "financial:" + r.Symbol + ":" + string(r.ReportType)
}

// Response represents a financial data response
type Response struct {
	Symbol     string          // Symbol code
	Data       []FinancialData // Financial data list
	Source     string          // Data source name
	Total      int             // Total records
	PageNumber int             // Current page
	PageSize   int             // Page size
}

// FinancialData represents a single financial report record
type FinancialData struct {
	// Report Information
	ReportDate   time.Time    // Report date
	ReportPeriod ReportPeriod // Period type
	FiscalYear   int          // Fiscal year
	FiscalPeriod int          // Fiscal period (1-4)

	// Income Statement Items
	TotalRevenue       float64 // 营业总收入
	Revenue            float64 // 营业收入
	TotalOperatingCost float64 // 营业总成本
	OperatingCost      float64 // 营业成本
	GrossProfit        float64 // 毛利润
	ResearchExpense    float64 // 研发费用
	SalesExpense       float64 // 销售费用
	AdminExpense       float64 // 管理费用
	FinancialExpense   float64 // 财务费用
	OperatingProfit    float64 // 营业利润
	TotalProfit        float64 // 利润总额
	NetProfit          float64 // 净利润
	NetProfitParent    float64 // 归母净利润
	NetProfitDeducted  float64 // 扣非净利润
	EPS                float64 // 每股收益
	DEPS               float64 // 稀释每股收益

	// Balance Sheet Items
	TotalAssets           float64 // 总资产
	CurrentAssets         float64 // 流动资产
	NonCurrentAssets      float64 // 非流动资产
	TotalLiabilities      float64 // 总负债
	CurrentLiabilities    float64 // 流动负债
	NonCurrentLiabilities float64 // 非流动负债
	TotalOwnerEquity      float64 // 所有者权益
	TotalEquity           float64 // 总股本
	FloatEquity           float64 // 流通股本
	CapitalReserve        float64 // 资本公积
	SurplusReserve        float64 // 盈余公积
	UndistributedProfit   float64 // 未分配利润

	// Cash Flow Statement Items
	OperatingCashFlow float64 // 经营活动现金流
	InvestingCashFlow float64 // 投资活动现金流
	FinancingCashFlow float64 // 筹资活动现金流
	NetCashFlow       float64 // 现金净增加额
	CashEquivalents   float64 // 期末现金等价物

	// Financial Indicators
	ROE                float64 // 净资产收益率 (%)
	ROA                float64 // 总资产收益率 (%)
	GrossMargin        float64 // 毛利率 (%)
	NetMargin          float64 // 净利率 (%)
	AssetTurnover      float64 // 总资产周转率
	EquityMultiplier   float64 // 权益乘数
	CurrentRatio       float64 // 流动比率
	QuickRatio         float64 // 速动比率
	DebtToAsset        float64 // 资产负债率 (%)
	InterestCoverage   float64 // 利息保障倍数
	InventoryTurnover  float64 // 存货周转率
	ReceivableTurnover float64 // 应收账款周转率
}

// FinanceData is an alias for FinancialData for backward compatibility
type FinanceData = FinancialData

// Source defines the interface for financial data providers
type Source interface {
	Name() string
	Fetch(ctx context.Context, req Request) (Response, error)
	HealthCheck(ctx context.Context) error
}
