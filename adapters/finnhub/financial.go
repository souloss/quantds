package finnhub

import (
	"context"
	"strconv"
	"time"

	"github.com/souloss/quantds/clients/finnhub"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/financial"
	"github.com/souloss/quantds/manager"
)

// FinancialAdapter adapts Finnhub financial data for US stocks
type FinancialAdapter struct {
	client *finnhub.Client
}

// NewFinancialAdapter creates a new financial adapter
func NewFinancialAdapter(client *finnhub.Client) *FinancialAdapter {
	return &FinancialAdapter{client: client}
}

func (a *FinancialAdapter) Name() string                      { return Name }
func (a *FinancialAdapter) SupportedMarkets() []domain.Market { return usMarkets }

func (a *FinancialAdapter) CanHandle(symbol string) bool {
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

func (a *FinancialAdapter) Fetch(ctx context.Context, req financial.Request) (financial.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	var sym domain.Symbol
	if err := sym.Parse(req.Symbol); err != nil {
		return financial.Response{}, trace, err
	}

	// Map domain report type to Finnhub statement type
	statement := mapReportTypeToStatement(req.ReportType)
	frequency := "annual"
	if req.Period == financial.PeriodQuarterly {
		frequency = "quarterly"
	}

	result, record, err := a.client.GetFinancials(ctx, &finnhub.FinancialsParams{
		Symbol:    sym.Code,
		Statement: statement,
		Frequency: frequency,
	})
	trace.AddRequest(record)

	if err != nil {
		return financial.Response{}, trace, err
	}

	data := make([]financial.FinancialData, 0, len(result.Data))
	for _, d := range result.Data {
		fd := financial.FinancialData{
			FiscalYear: d.Year,
		}

		// Parse fiscal period from quarter string
		if d.Quarter != "" {
			if q, err := strconv.Atoi(d.Quarter); err == nil {
				fd.FiscalPeriod = q
			}
		}

		// Parse report date from EndDate
		if d.EndDate != "" {
			if t, err := time.Parse("2006-01-02", d.EndDate); err == nil {
				fd.ReportDate = t
			}
		}

		if req.Period == financial.PeriodQuarterly {
			fd.ReportPeriod = financial.PeriodQuarterly
		} else {
			fd.ReportPeriod = financial.PeriodAnnual
		}

		// Map Finnhub financial items to domain fields
		mapFinancialItems(d.Items, &fd, req.ReportType)

		data = append(data, fd)
	}

	trace.Finish()
	return financial.Response{
		Symbol:      req.Symbol,
		Data:        data,
		Source:      Name,
		Total:       len(data),
		DataVersion: 1,
	}, trace, nil
}

// mapReportTypeToStatement maps domain report type to Finnhub statement parameter
func mapReportTypeToStatement(rt financial.ReportType) string {
	switch rt {
	case financial.ReportTypeBalance:
		return "bs"
	case financial.ReportTypeIncome:
		return "ic"
	case financial.ReportTypeCashflow:
		return "cf"
	default:
		return "ic"
	}
}

// mapFinancialItems maps Finnhub key-value items to domain FinancialData fields
func mapFinancialItems(items map[string]float64, fd *financial.FinancialData, rt financial.ReportType) {
	switch rt {
	case financial.ReportTypeIncome:
		fd.TotalRevenue = items["totalRevenue"]
		fd.Revenue = items["revenue"]
		fd.CostOfRevenue = items["costOfRevenue"]
		fd.GrossProfit = items["grossProfit"]
		fd.ResearchExpense = items["researchAndDevelopment"]
		fd.SellingExpense = items["sellingGeneralAndAdministrative"]
		fd.OperatingProfit = items["operatingIncome"]
		fd.TotalProfit = items["incomeBeforeTax"]
		fd.NetProfit = items["netIncome"]
		fd.NetProfitParent = items["netIncomeContinuousOperations"]
		fd.InterestExpense = items["interestExpense"]
		fd.EPS = items["basicEPS"]
		fd.DEPS = items["dilutedEPS"]

	case financial.ReportTypeBalance:
		fd.TotalAssets = items["totalAssets"]
		fd.CurrentAssets = items["totalCurrentAssets"]
		fd.NonCurrentAssets = items["totalNonCurrentAssets"]
		fd.TotalLiabilities = items["totalLiabilities"]
		fd.CurrentLiabilities = items["totalCurrentLiabilities"]
		fd.NonCurrentLiabilities = items["totalNonCurrentLiabilities"]
		fd.TotalOwnerEquity = items["totalShareholderEquity"]
		fd.CashAndEquivalents = items["cashAndShortTermInvestments"]
		fd.Inventory = items["inventory"]
		fd.AccountsReceivable = items["netReceivables"]
		fd.LongTermDebt = items["longTermDebt"]
		fd.ShortTermDebt = items["shortTermDebt"]

	case financial.ReportTypeCashflow:
		fd.OperatingCashFlow = items["operatingCashFlow"]
		fd.InvestingCashFlow = items["netInvestingCashFlow"]
		fd.FinancingCashFlow = items["netFinancingCashFlow"]
		fd.Depreciation = items["depreciationAndAmortization"]
		fd.CapitalExpenditure = items["capitalExpenditure"]
		fd.FreeCashFlow = items["freeCashFlow"]
	}
}

var _ manager.Provider[financial.Request, financial.Response] = (*FinancialAdapter)(nil)
