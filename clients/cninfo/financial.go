// Package cninfo provides financial data APIs.
//
// 巨潮资讯财务数据API，通过深证信数据服务平台获取上市公司财务报表数据。
package cninfo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/souloss/quantds/request"
)

const (
	FinancialAPI = "/api/sysapi/p_sysapi1076"
)

// FinancialParams represents parameters for financial data request
// 财务数据查询参数
type FinancialParams struct {
	StockCode  string // 股票代码，如 000001
	Code       string // 组织机构代码（可选）
	ReportDate string // 报告期，格式 YYYY-MM-DD（可选）
	PageNum    int    // 页码，从1开始
	PageSize   int    // 每页数量
}

// FinancialResult represents financial data result
// 财务数据结果
type FinancialResult struct {
	Data  []FinancialRow
	Total int
}

// FinancialRow represents a single financial data row
// 财务数据行
type FinancialRow struct {
	// 基本信息
	StockCode  string // 股票代码
	StockName  string // 股票简称
	ReportDate string // 报告期
	AnnDate    string // 公告日期

	// 利润表
	TotalRevenue    float64 // 营业总收入
	Revenue         float64 // 营业收入
	TotalCost       float64 // 营业总成本
	OperatingCost   float64 // 营业成本
	GrossProfit     float64 // 毛利润
	OperatingProfit float64 // 营业利润
	TotalProfit     float64 // 利润总额
	NetProfit       float64 // 净利润
	NetProfitParent float64 // 归母净利润
	BasicEPS        float64 // 基本每股收益
	DilutedEPS      float64 // 稀释每股收益

	// 资产负债表
	TotalAssets        float64 // 资产总计
	TotalLiabilities   float64 // 负债合计
	TotalEquity        float64 // 所有者权益合计
	CurrentAssets      float64 // 流动资产合计
	CurrentLiabilities float64 // 流动负债合计

	// 现金流量表
	OperatingCashFlow float64 // 经营活动现金流
	InvestingCashFlow float64 // 投资活动现金流
	FinancingCashFlow float64 // 筹资活动现金流
	NetCashIncrease   float64 // 现金净增加额
}

// GetFinancialData retrieves financial data from Cninfo
// 获取上市公司财务数据
//
// API: /api/sysapi/p_sysapi1076
// 方法: POST
// 认证: 无需认证，但需要正确的请求头
//
// 参数说明：
//   - stockCode: 股票代码，如 000001
//   - reportDate: 报告期，格式 YYYY-MM-DD（可选）
//   - pageNum: 页码，从1开始
//   - pageSize: 每页数量
//
// 限制：
//   - 需要正确的Referer和User-Agent
//   - 高频请求可能被限流
func (c *Client) GetFinancialData(ctx context.Context, params *FinancialParams) (*FinancialResult, *request.Record, error) {
	if params == nil || params.StockCode == "" {
		return nil, nil, fmt.Errorf("stockCode required")
	}

	// 默认值
	if params.PageNum <= 0 {
		params.PageNum = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	form := url.Values{}
	form.Set("scode", params.StockCode)
	if params.Code != "" {
		form.Set("code", params.Code)
	}
	if params.ReportDate != "" {
		form.Set("reportDate", params.ReportDate)
	}
	form.Set("pageNum", strconv.Itoa(params.PageNum))
	form.Set("pageSize", strconv.Itoa(params.PageSize))

	var raw struct {
		Data struct {
			List []struct {
				StockCode          string `json:"scode"`
				StockName          string `json:"sname"`
				ReportDate         string `json:"reportdate"`
				AnnDate            string `json:"declaredate"`
				TotalRevenue       string `json:"totaloperatereve"`
				Revenue            string `json:"operatereve"`
				TotalCost          string `json:"totaloperatecost"`
				OperatingCost      string `json:"operatecost"`
				OperatingProfit    string `json:"operateprofit"`
				TotalProfit        string `json:"totalprofit"`
				NetProfit          string `json:"netprofit"`
				NetProfitParent    string `json:"parentnetprofit"`
				BasicEPS           string `json:"basiceps"`
				TotalAssets        string `json:"totalassets"`
				TotalLiabilities   string `json:"totalliability"`
				TotalEquity        string `json:"totalshareholder"`
				CurrentAssets      string `json:"totalcurrentassets"`
				CurrentLiabilities string `json:"totalcurrentliability"`
				OperatingCF        string `json:"netoperatecashflow"`
				InvestingCF        string `json:"netinvestcashflow"`
				FinancingCF        string `json:"netfinancecashflow"`
			} `json:"list"`
			Total int `json:"total"`
		} `json:"data"`
	}

	record, err := c.do(ctx, "POST", FinancialAPI, form, &raw)
	if err != nil {
		return nil, record, err
	}

	rows := make([]FinancialRow, 0, len(raw.Data.List))
	for _, item := range raw.Data.List {
		row := FinancialRow{
			StockCode:          item.StockCode,
			StockName:          item.StockName,
			ReportDate:         item.ReportDate,
			AnnDate:            item.AnnDate,
			TotalRevenue:       parseFinancialFloat(item.TotalRevenue),
			Revenue:            parseFinancialFloat(item.Revenue),
			TotalCost:          parseFinancialFloat(item.TotalCost),
			OperatingCost:      parseFinancialFloat(item.OperatingCost),
			OperatingProfit:    parseFinancialFloat(item.OperatingProfit),
			TotalProfit:        parseFinancialFloat(item.TotalProfit),
			NetProfit:          parseFinancialFloat(item.NetProfit),
			NetProfitParent:    parseFinancialFloat(item.NetProfitParent),
			BasicEPS:           parseFinancialFloat(item.BasicEPS),
			TotalAssets:        parseFinancialFloat(item.TotalAssets),
			TotalLiabilities:   parseFinancialFloat(item.TotalLiabilities),
			TotalEquity:        parseFinancialFloat(item.TotalEquity),
			CurrentAssets:      parseFinancialFloat(item.CurrentAssets),
			CurrentLiabilities: parseFinancialFloat(item.CurrentLiabilities),
			OperatingCashFlow:  parseFinancialFloat(item.OperatingCF),
			InvestingCashFlow:  parseFinancialFloat(item.InvestingCF),
			FinancingCashFlow:  parseFinancialFloat(item.FinancingCF),
		}
		row.GrossProfit = row.Revenue - row.OperatingCost
		rows = append(rows, row)
	}

	return &FinancialResult{
		Data:  rows,
		Total: raw.Data.Total,
	}, record, nil
}

// ProfileParams represents parameters for company profile request
// 公司概况查询参数
type ProfileParams struct {
	StockCode string // 股票代码
}

// ProfileResult represents company profile data
// 公司概况数据
type ProfileResult struct {
	StockCode    string // 股票代码
	StockName    string // 股票简称
	FullName     string // 公司全称
	EnglishName  string // 英文名称
	ListDate     string // 上市日期
	Province     string // 省份
	City         string // 城市
	Industry     string // 行业
	Website      string // 公司网址
	Email        string // 电子邮箱
	Office       string // 办公地址
	RegAddress   string // 注册地址
	Chairman     string // 董事长
	Manager      string // 总经理
	Secretary    string // 董秘
	RegCapital   string // 注册资本
	SetupDate    string // 成立日期
	Employees    string // 员工人数
	MainBusiness string // 主营业务
	Description  string // 公司简介
}

// GetProfile retrieves company profile from Cninfo
// 获取上市公司概况信息
//
// API: /api/sysapi/p_sysapi1073
// 方法: POST
//
// 参数说明：
//   - stockCode: 股票代码
func (c *Client) GetProfile(ctx context.Context, params *ProfileParams) (*ProfileResult, *request.Record, error) {
	if params == nil || params.StockCode == "" {
		return nil, nil, fmt.Errorf("stockCode required")
	}

	form := url.Values{}
	form.Set("scode", params.StockCode)

	var raw struct {
		Data struct {
			StockCode    string `json:"scode"`
			StockName    string `json:"sname"`
			FullName     string `json:"fullnamex"`
			EnglishName  string `json:"englishname"`
			ListDate     string `json:"listingdate"`
			Province     string `json:"province"`
			City         string `json:"cityname"`
			Industry     string `json:"industryname"`
			Website      string `json:"website"`
			Email        string `json:"email"`
			Office       string `json:"officceaddress"`
			RegAddress   string `json:"registeredaddress"`
			Chairman     string `json:"chairman"`
			Manager      string `json:"manager"`
			Secretary    string `json:"secretary"`
			RegCapital   string `json:"registeredcapital"`
			SetupDate    string `json:"setupdate"`
			Employees    string `json:"staffnum"`
			MainBusiness string `json:"mainbusiness"`
			Description  string `json:"brief"`
		} `json:"data"`
	}

	record, err := c.do(ctx, "POST", "/api/sysapi/p_sysapi1073", form, &raw)
	if err != nil {
		return nil, record, err
	}

	return &ProfileResult{
		StockCode:    raw.Data.StockCode,
		StockName:    raw.Data.StockName,
		FullName:     raw.Data.FullName,
		EnglishName:  raw.Data.EnglishName,
		ListDate:     raw.Data.ListDate,
		Province:     raw.Data.Province,
		City:         raw.Data.City,
		Industry:     raw.Data.Industry,
		Website:      raw.Data.Website,
		Email:        raw.Data.Email,
		Office:       raw.Data.Office,
		RegAddress:   raw.Data.RegAddress,
		Chairman:     raw.Data.Chairman,
		Manager:      raw.Data.Manager,
		Secretary:    raw.Data.Secretary,
		RegCapital:   raw.Data.RegCapital,
		SetupDate:    raw.Data.SetupDate,
		Employees:    raw.Data.Employees,
		MainBusiness: raw.Data.MainBusiness,
		Description:  raw.Data.Description,
	}, record, nil
}

func parseFinancialFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// do is a helper method to make API requests
func (c *Client) doAPI(ctx context.Context, method, path string, form url.Values, result any) (*request.Record, error) {
	u := c.baseURL + path
	var body []byte
	var contentType string

	if form != nil {
		body = []byte(form.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	req := request.Request{
		Method: method,
		URL:    u,
		Headers: map[string]string{
			"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Accept":           "application/json, text/javascript, */*",
			"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8",
			"Referer":          c.baseURL,
			"X-Requested-With": "XMLHttpRequest",
		},
	}
	if len(body) > 0 {
		req.Headers["Content-Type"] = contentType
		req.Body = body
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return record, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return record, fmt.Errorf("http status %d", resp.StatusCode)
	}

	if result != nil {
		if err := json.Unmarshal(resp.Body, result); err != nil {
			return record, fmt.Errorf("decode response: %w", err)
		}
	}

	return record, nil
}
