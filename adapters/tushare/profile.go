package tushare

import (
	"context"

	"github.com/souloss/quantds/clients/tushare"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/profile"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// ProfileAdapter 将 Tushare 的 stock_basic、stock_company、daily_basic
// 聚合为统一的 profile 域类型。
type ProfileAdapter struct {
	client *tushare.Client
}

// NewProfileAdapter 创建 Tushare 个股档案适配器。
func NewProfileAdapter(client *tushare.Client) *ProfileAdapter {
	return &ProfileAdapter{client: client}
}

func (a *ProfileAdapter) Name() string                      { return Name }
func (a *ProfileAdapter) SupportedMarkets() []domain.Market { return supportedMarkets }

func (a *ProfileAdapter) CanHandle(symbol string) bool {
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

func (a *ProfileAdapter) Fetch(ctx context.Context, _ request.Client, req profile.Request) (profile.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	tsCode, err := tushare.ToTushareSymbol(req.Symbol)
	if err != nil {
		return profile.Response{}, trace, err
	}

	p := profile.Profile{Symbol: req.Symbol}

	// 1. 股票基本信息
	basicRows, basicRec, err := a.client.GetStockBasic(ctx, &tushare.StockBasicParams{TSCode: tsCode})
	trace.AddRequest(basicRec)
	if err == nil && len(basicRows) > 0 {
		row := basicRows[0]
		p.Name = row.Name
		p.Industry = row.Industry
		listDate := row.ListDate
		if len(listDate) == 8 {
			listDate = listDate[:4] + "-" + listDate[4:6] + "-" + listDate[6:]
		}
		p.ListingDate = listDate
		if row.DelistDate != "" {
			delistDate := row.DelistDate
			if len(delistDate) == 8 {
				delistDate = delistDate[:4] + "-" + delistDate[4:6] + "-" + delistDate[6:]
			}
			p.DelistDate = delistDate
		}
		p.Sector = row.Area
		p.Currency = "CNY"
	}

	// 2. 公司信息
	compRows, compRec, err := a.client.GetStockCompany(ctx, &tushare.StockCompanyParams{TSCode: tsCode})
	trace.AddRequest(compRec)
	if err == nil && len(compRows) > 0 {
		row := compRows[0]
		p.Chairman = row.Chairman
		p.CEO = row.Manager
		p.Secretary = row.Secretary
		p.RegCapital = row.RegCapital
		p.SetupDate = formatCompanyDate(row.SetupDate)
		p.Province = row.Province
		p.City = row.City
		p.Introduction = row.Introduction
		p.Website = row.Website
		p.Email = row.Email
		p.Address = row.Office
		p.Employees = row.Employees
		p.MainBusiness = row.MainBusiness
		p.BusinessScope = row.BusinessScope
	}

	// 3. 最新每日指标（PE/PB/PS/市值/换手率等）
	dbRows, dbRec, err := a.client.GetDailyBasic(ctx, &tushare.DailyBasicParams{TSCode: tsCode})
	trace.AddRequest(dbRec)
	if err == nil && len(dbRows) > 0 {
		row := dbRows[0]
		p.TradeDate = formatCompanyDate(row.TradeDate)
		p.PE = row.PE
		p.PETrailing = row.PETTM
		p.PB = row.PB
		p.PS = row.PS
		p.DivYield = row.DvRatio
		p.TurnoverRate = row.TurnoverRate
		p.VolumeRatio = row.VolumeRatio
		p.TotalShares = row.TotalShare * 10000   // 万股 → 股
		p.FloatShares = row.FloatShare * 10000    // 万股 → 股
		p.MarketCap = row.TotalMV * 10000         // 万元 → 元
		p.FloatCap = row.CircMV * 10000           // 万元 → 元
	}

	// 4. 财务指标（ROE/ROA/毛利率/净利率等）
	// 不指定日期，默认返回最新记录
	fiRows, fiRec, err := a.client.GetFinaIndicator(ctx, &tushare.FinaIndicatorParams{TSCode: tsCode})
	trace.AddRequest(fiRec)
	if err == nil && len(fiRows) > 0 {
		row := fiRows[0]
		p.ROE = row.ROE
		p.ROA = row.ROA
		p.GrossMargin = row.GrossProfitMargin
		p.NetMargin = row.NetProfitMargin
		// EPS/BPS 等如果 fina_indicator 里没有，可以从 income/balance sheet 计算，或者 daily_basic 可能有？
		// daily_basic 只有 PE/PB/PS
		// fina_indicator 有 basic_eps? Check fina_indicator.go
		// fina_indicator fields: roe, roe_waa, roa, netprofit_margin, grossprofit_margin, current_ratio, quick_ratio, debt_to_assets, turn_days, roa_yearly, roe_avg, assets_turn, op_income, ebit, ebitda
		// 没有 EPS/BPS。
		// 可以从 income (basic_eps) 获取 EPS?
		// 可以从 daily_basic 推算 BPS = Price / PB? 或者 TotalAssets / TotalShares?
		// 暂时先只填这些。
	}

	trace.Finish()
	return profile.Response{
		Data:   p,
		Source: Name,
	}, trace, nil
}

// formatCompanyDate 将 YYYYMMDD 格式转为 YYYY-MM-DD。
func formatCompanyDate(yyyymmdd string) string {
	if len(yyyymmdd) != 8 {
		return yyyymmdd
	}
	return yyyymmdd[:4] + "-" + yyyymmdd[4:6] + "-" + yyyymmdd[6:]
}

var _ manager.Provider[profile.Request, profile.Response] = (*ProfileAdapter)(nil)
