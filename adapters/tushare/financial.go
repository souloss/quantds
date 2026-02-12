package tushare

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/tushare"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/financial"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// FinancialAdapter 将 Tushare 的利润表、资产负债表、现金流量表和财务指标
// 聚合为统一的 financial 域类型。
type FinancialAdapter struct {
	client *tushare.Client
}

// NewFinancialAdapter 创建 Tushare 财务数据适配器。
func NewFinancialAdapter(client *tushare.Client) *FinancialAdapter {
	return &FinancialAdapter{client: client}
}

func (a *FinancialAdapter) Name() string                      { return Name }
func (a *FinancialAdapter) SupportedMarkets() []domain.Market { return supportedMarkets }

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

func (a *FinancialAdapter) Fetch(ctx context.Context, _ request.Client, req financial.Request) (financial.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	tsCode, err := tushare.ToTushareSymbol(req.Symbol)
	if err != nil {
		return financial.Response{}, trace, err
	}

	// 并行获取利润表、资产负债表、现金流量表、财务指标
	// 这里使用串行调用以保持简单性（Tushare 有频率限制）

	// 1. 利润表
	incomeRows, incRec, err := a.client.GetIncome(ctx, &tushare.IncomeParams{
		TSCode:     tsCode,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		ReportType: "1", // 合并报表
	})
	trace.AddRequest(incRec)
	if err != nil {
		return financial.Response{}, trace, err
	}

	// 2. 资产负债表
	bsRows, bsRec, err := a.client.GetBalanceSheet(ctx, &tushare.BalanceSheetParams{
		TSCode:     tsCode,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		ReportType: "1",
	})
	trace.AddRequest(bsRec)

	// 3. 现金流量表
	cfRows, cfRec, err := a.client.GetCashflow(ctx, &tushare.CashflowParams{
		TSCode:     tsCode,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		ReportType: "1",
	})
	trace.AddRequest(cfRec)

	// 4. 财务指标
	fiRows, fiRec, err := a.client.GetFinaIndicator(ctx, &tushare.FinaIndicatorParams{
		TSCode:    tsCode,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})
	trace.AddRequest(fiRec)

	// 以利润表为主键，按 end_date 关联其他数据
	bsMap := make(map[string]*tushare.BalanceSheetRow, len(bsRows))
	for i := range bsRows {
		bsMap[bsRows[i].EndDate] = &bsRows[i]
	}
	cfMap := make(map[string]*tushare.CashflowRow, len(cfRows))
	for i := range cfRows {
		cfMap[cfRows[i].EndDate] = &cfRows[i]
	}
	fiMap := make(map[string]*tushare.FinaIndicatorRow, len(fiRows))
	for i := range fiRows {
		fiMap[fiRows[i].EndDate] = &fiRows[i]
	}

	data := make([]financial.FinancialData, 0, len(incomeRows))
	for _, inc := range incomeRows {
		fd := financial.FinancialData{
			ReportDate:         parseEndDate(inc.EndDate),
			ReportPeriod:       guessReportPeriod(inc.EndDate),
			FiscalYear:         parseYear(inc.EndDate),
			FiscalPeriod:       parsePeriod(inc.EndDate),
			TotalRevenue:       inc.TotalRevenue,
			Revenue:            inc.Revenue,
			TotalOperatingCost: inc.TotalCogs,
			OperatingCost:      inc.OperCost,
			ResearchExpense:    inc.RDExp,
			SalesExpense:       inc.SellExp,
			AdminExpense:       inc.AdminExp,
			FinancialExpense:   inc.FinExp,
			OperatingProfit:    inc.OperProfit,
			TotalProfit:        inc.TotalProfit,
			NetProfit:          inc.NIncome,
			NetProfitParent:    inc.NIncomeAttrP,
			EPS:                inc.BasicEPS,
			DEPS:               inc.DilutedEPS,
		}

		// 合并资产负债表数据
		if bs, ok := bsMap[inc.EndDate]; ok {
			fd.TotalAssets = bs.TotalAssets
			fd.CurrentAssets = bs.TotalCurAssets
			fd.NonCurrentAssets = bs.TotalNCA
			fd.TotalLiabilities = bs.TotalLiab
			fd.CurrentLiabilities = bs.TotalCurLiab
			fd.NonCurrentLiabilities = bs.TotalNCL
			fd.TotalOwnerEquity = bs.TotalHldrEqy
			fd.CapitalReserve = bs.CapRese
			fd.SurplusReserve = bs.SurplusRese
			fd.UndistributedProfit = bs.UndistProfit
		}

		// 合并现金流量表数据
		if cf, ok := cfMap[inc.EndDate]; ok {
			fd.OperatingCashFlow = cf.NCashflowAct
			fd.InvestingCashFlow = cf.NCashflowInv
			fd.FinancingCashFlow = cf.NCashflowFnc
			fd.CashEquivalents = cf.CCashEquEnd
		}

		// 合并财务指标数据
		if fi, ok := fiMap[inc.EndDate]; ok {
			fd.ROE = fi.ROE
			fd.ROA = fi.ROA
			fd.GrossMargin = fi.GrossProfitMargin
			fd.NetMargin = fi.NetProfitMargin
			fd.AssetTurnover = fi.AssetsTurn
			fd.CurrentRatio = fi.CurrentRatio
			fd.QuickRatio = fi.QuickRatio
			fd.DebtToAsset = fi.DebtToAssets
		}

		data = append(data, fd)
	}

	trace.Finish()
	return financial.Response{
		Symbol: req.Symbol,
		Data:   data,
		Source: Name,
		Total:  len(data),
	}, trace, nil
}

// parseEndDate 将 YYYYMMDD 格式日期解析为 time.Time。
func parseEndDate(yyyymmdd string) time.Time {
	if len(yyyymmdd) != 8 {
		return time.Time{}
	}
	t, _ := time.ParseInLocation("20060102", yyyymmdd, time.Local)
	return t
}

// parseYear 从 YYYYMMDD 中提取年份。
func parseYear(yyyymmdd string) int {
	if len(yyyymmdd) < 4 {
		return 0
	}
	y := 0
	for _, c := range yyyymmdd[:4] {
		y = y*10 + int(c-'0')
	}
	return y
}

// parsePeriod 从 YYYYMMDD 中推断报告期序号 (1-4)。
func parsePeriod(yyyymmdd string) int {
	if len(yyyymmdd) < 6 {
		return 0
	}
	mm := yyyymmdd[4:6]
	switch mm {
	case "03":
		return 1 // Q1
	case "06":
		return 2 // Q2/半年报
	case "09":
		return 3 // Q3
	case "12":
		return 4 // Q4/年报
	default:
		return 0
	}
}

// guessReportPeriod 从 end_date 推断报告期类型。
func guessReportPeriod(yyyymmdd string) financial.ReportPeriod {
	if len(yyyymmdd) < 6 {
		return financial.PeriodQuarterly
	}
	mm := yyyymmdd[4:6]
	switch mm {
	case "12":
		return financial.PeriodAnnual
	case "06":
		return financial.PeriodSemiAnnual
	default:
		return financial.PeriodQuarterly
	}
}

var _ manager.Provider[financial.Request, financial.Response] = (*FinancialAdapter)(nil)
