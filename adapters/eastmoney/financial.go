package eastmoney

import (
	"context"
	"fmt"
	"time"

	"github.com/souloss/quantds/clients/eastmoney"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/financial"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// FinancialAdapter adapts Eastmoney financial data
type FinancialAdapter struct {
	client *eastmoney.Client
}

// NewFinancialAdapter creates a new financial adapter
func NewFinancialAdapter(client *eastmoney.Client) *FinancialAdapter {
	return &FinancialAdapter{client: client}
}

// Name returns the adapter name
func (a *FinancialAdapter) Name() string {
	return Name
}

// SupportedMarkets returns supported markets
func (a *FinancialAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

// CanHandle checks if the adapter can handle the symbol
func (a *FinancialAdapter) CanHandle(symbol string) bool {
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

// Fetch retrieves financial data
func (a *FinancialAdapter) Fetch(ctx context.Context, _ request.Client, req financial.Request) (financial.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	// 1. Income Statement
	incomeParams := &eastmoney.FinancialParams{
		ReportName: eastmoney.ReportIncome,
		Code:       req.Symbol,
		PageNumber: 1,
		PageSize:   100,
	}
	incomeRes, incRec, err := a.client.GetFinancials(ctx, incomeParams)
	trace.AddRequest(incRec)
	if err != nil {
		return financial.Response{}, trace, err
	}

	// 2. Balance Sheet
	bsParams := &eastmoney.FinancialParams{
		ReportName: eastmoney.ReportBalanceSheet,
		Code:       req.Symbol,
		PageNumber: 1,
		PageSize:   100,
	}
	bsRes, bsRec, bsErr := a.client.GetFinancials(ctx, bsParams)
	trace.AddRequest(bsRec)

	// 3. Cash Flow
	cfParams := &eastmoney.FinancialParams{
		ReportName: eastmoney.ReportCashflow,
		Code:       req.Symbol,
		PageNumber: 1,
		PageSize:   100,
	}
	cfRes, cfRec, cfErr := a.client.GetFinancials(ctx, cfParams)
	trace.AddRequest(cfRec)

	// Index by Report Date — 忽略 BS/CF 的错误，仍然返回可用的收入数据
	bsMap := make(map[string]map[string]interface{})
	if bsErr == nil && bsRes != nil {
		for _, row := range bsRes.Data {
			date := getString(row, "REPORT_DATE")
			if date != "" {
				bsMap[date] = row
			}
		}
	}

	cfMap := make(map[string]map[string]interface{})
	if cfErr == nil && cfRes != nil {
		for _, row := range cfRes.Data {
			date := getString(row, "REPORT_DATE")
			if date != "" {
				cfMap[date] = row
			}
		}
	}

	data := make([]financial.FinancialData, 0, len(incomeRes.Data))
	for _, row := range incomeRes.Data {
		date := getString(row, "REPORT_DATE")
		item := financial.FinancialData{
			ReportDate:         parseReportDate(date),
			TotalRevenue:       getFloat(row, "TOTAL_OPERATE_INCOME"),
			Revenue:            getFloat(row, "OPERATE_INCOME"),
			TotalOperatingCost: getFloat(row, "TOTAL_OPERATE_COST"),
			OperatingCost:      getFloat(row, "OPERATE_COST"),
			ResearchExpense:    getFloat(row, "RESEARCH_EXPENSE"),
			SalesExpense:       getFloat(row, "SALE_EXPENSE"),
			AdminExpense:       getFloat(row, "MANAGE_EXPENSE"),
			FinancialExpense:   getFloat(row, "FINANCE_EXPENSE"),
			OperatingProfit:    getFloat(row, "OPERATE_PROFIT"),
			TotalProfit:        getFloat(row, "TOTAL_PROFIT"),
			NetProfit:          getFloat(row, "NETPROFIT"),
			NetProfitParent:    getFloat(row, "PARENT_NETPROFIT"),
			NetProfitDeducted:  getFloat(row, "DEDUCT_PARENT_NETPROFIT"),
			EPS:                getFloat(row, "BASIC_EPS"),
			DEPS:               getFloat(row, "DILUTED_EPS"),
		}

		// Merge Balance Sheet
		if bsRow, ok := bsMap[date]; ok {
			item.TotalAssets = getFloat(bsRow, "TOTAL_ASSETS")
			item.CurrentAssets = getFloat(bsRow, "TOTAL_CURRENT_ASSETS")
			item.NonCurrentAssets = getFloat(bsRow, "TOTAL_NONCURRENT_ASSETS")
			item.TotalLiabilities = getFloat(bsRow, "TOTAL_LIABILITIES")
			item.CurrentLiabilities = getFloat(bsRow, "TOTAL_CURRENT_LIABILITIES")
			item.NonCurrentLiabilities = getFloat(bsRow, "TOTAL_NONCURRENT_LIABILITIES")
			item.TotalOwnerEquity = getFloat(bsRow, "TOTAL_EQUITY") // 归母权益? Need check field mapping. usually TOTAL_PARENT_EQUITY or TOTAL_EQUITY
			item.TotalEquity = getFloat(bsRow, "TOTAL_SHARE") // Share capital
			item.CapitalReserve = getFloat(bsRow, "CAPITAL_RESERVE")
			item.SurplusReserve = getFloat(bsRow, "SURPLUS_RESERVE")
			item.UndistributedProfit = getFloat(bsRow, "UNDISTRIBUTED_PROFIT")
		}

		// Merge Cash Flow
		if cfRow, ok := cfMap[date]; ok {
			item.OperatingCashFlow = getFloat(cfRow, "NETCASH_OPERATE")
			item.InvestingCashFlow = getFloat(cfRow, "NETCASH_INVEST")
			item.FinancingCashFlow = getFloat(cfRow, "NETCASH_FINANCE")
		}

		data = append(data, item)
	}

	trace.Finish()
	return financial.Response{
		Symbol: req.Symbol,
		Data:   data,
		Source: Name,
	}, trace, nil
}

func getString(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case string:
			return val
		case float64:
			return fmt.Sprintf("%.0f", val)
		}
	}
	return ""
}

func getFloat(data map[string]interface{}, key string) float64 {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case string:
			var f float64
			fmt.Sscanf(val, "%f", &f)
			return f
		}
	}
	return 0
}

func parseReportDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}
	t, _ := time.Parse("2006-01-02 15:04:05", dateStr)
	return t
}

var _ manager.Provider[financial.Request, financial.Response] = (*FinancialAdapter)(nil)
