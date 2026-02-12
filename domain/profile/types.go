// Package profile provides security profile/details domain types.
//
// This package defines the request/response types for detailed security
// information including company profile, financial metrics, and trading
// statistics. It consolidates the former stockdetail and company domain types.
package profile

import "context"

// Request represents a security profile request.
type Request struct {
	Symbol string // 标的代码
}

// CacheKey returns the cache key for the request.
func (r Request) CacheKey() string {
	return "profile:" + r.Symbol
}

// Response represents a security profile response.
type Response struct {
	Data   Profile // 个股档案数据
	Source string  // 数据源名称
}

// Profile represents detailed information about a security,
// including trading data, valuation metrics, company information, and classification.
type Profile struct {
	// 基本信息
	Symbol      string // 标的代码
	Name        string // 简称
	FullName    string // 全称
	TradeDate   string // 最新交易日
	ListingDate string // 上市日期
	DelistDate  string // 退市日期
	Currency    string // 币种

	// 行情数据
	Open          float64 // 开盘价
	High          float64 // 最高价
	Low           float64 // 最低价
	Close         float64 // 收盘价（最新）
	PreClose      float64 // 昨收价
	Volume        float64 // 成交量
	Amount        float64 // 成交额
	TurnoverRate  float64 // 换手率 (%)
	ChangePercent float64 // 涨跌幅 (%)
	Amplitude     float64 // 振幅 (%)
	VolumeRatio   float64 // 量比

	// 估值指标
	PE         float64 // 市盈率（动态）
	PEStatic   float64 // 市盈率（静态）
	PETrailing float64 // 市盈率（TTM）
	PB         float64 // 市净率
	PS         float64 // 市销率
	PCF        float64 // 市现率

	// 盈利指标
	ROE         float64 // 净资产收益率 (%)
	ROA         float64 // 总资产收益率 (%)
	GrossMargin float64 // 毛利率 (%)
	NetMargin   float64 // 净利率 (%)
	EPS         float64 // 每股收益
	BPS         float64 // 每股净资产
	DPS         float64 // 每股股利
	DivYield    float64 // 股息率 (%)

	// 股本结构
	TotalShares       float64 // 总股本
	FloatShares       float64 // 流通股本
	MarketCap         float64 // 总市值
	FloatCap          float64 // 流通市值
	MarketCapCategory string  // 市值分类 (Large/Mid/Small)

	// 财务概要（最新一期）
	Revenue     float64 // 营业收入
	NetProfit   float64 // 净利润
	TotalAssets float64 // 总资产
	TotalEquity float64 // 总权益
	TotalDebts  float64 // 总负债
	OperatingCF float64 // 经营活动现金流
	InvestingCF float64 // 投资活动现金流
	FinancingCF float64 // 筹资活动现金流

	// 行业分类
	Sector        string   // 板块
	Industry      string   // 行业
	SubIndustry   string   // 子行业
	SwSector      string   // 申万一级
	SwIndustry    string   // 申万二级
	SwSubIndustry string   // 申万三级
	Concept       string   // 概念题材
	ConceptTags   []string // 概念标签列表

	// 公司信息
	Province      string // 所属省份
	City          string // 所属城市
	Address       string // 注册地址
	Website       string // 公司网址
	Phone         string // 联系电话
	Fax           string // 传真
	Email         string // 电子邮箱
	Chairman      string // 董事长
	CEO           string // 总经理/CEO
	Secretary     string // 董秘
	RegCapital    float64 // 注册资本（万元）
	SetupDate     string  // 成立日期
	Employees     int     // 员工人数
	Introduction  string  // 公司简介
	MainBusiness  string  // 主营业务
	BusinessScope string  // 经营范围
}

// Source defines the interface for profile data providers.
type Source interface {
	Name() string
	Fetch(ctx context.Context, req Request) (Response, error)
	HealthCheck(ctx context.Context) error
}
