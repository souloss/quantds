// Package xueqiu provides stock profile/detail API.
//
// 雪球证券详情API，支持获取个股详细信息。
package xueqiu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/souloss/quantds/request"
)

const (
	QuoteDetailAPI = "/v5/stock/quote.json"
	ProfileAPI     = "/v5/stock/profile.json"
)

// QuoteDetailParams represents parameters for quote detail request
// 个股详情查询参数
type QuoteDetailParams struct {
	Symbol string // 股票代码，如 SH600000
	Extend bool   // 是否获取扩展信息
}

// QuoteDetailResult represents detailed quote information
// 个股详情信息
type QuoteDetailResult struct {
	// 基本信息
	Symbol   string // 标准代码
	Code     string // 股票代码
	Name     string // 股票名称
	Exchange string // 交易所
	Type     int    // 类型：1-A股，2-港股，3-美股

	// 行情数据
	Current      float64 // 当前价
	Percent      float64 // 涨跌幅(%)
	Change       float64 // 涨跌额
	High         float64 // 最高价
	Low          float64 // 最低价
	Open         float64 // 开盘价
	PreClose     float64 // 昨收价
	Volume       float64 // 成交量
	Amount       float64 // 成交额
	TurnoverRate float64 // 换手率(%)
	Amplitude    float64 // 振幅(%)

	// 估值数据
	PE             float64 // 市盈率（TTM）
	PB             float64 // 市净率
	PS             float64 // 市销率
	PCF            float64 // 市现率
	TotalMarketCap float64 // 总市值
	FloatMarketCap float64 // 流通市值
	DividendYield  float64 // 股息率(%)

	// 股本数据
	TotalShares float64 // 总股本
	FloatShares float64 // 流通股本

	// 公司信息
	Industry     string  // 行业
	Sector       string  // 板块
	ListDate     string  // 上市日期
	Province     string  // 省份
	City         string  // 城市
	MainBusiness string  // 主营业务
	Description  string  // 公司简介
	Website      string  // 公司网址
	Email        string  // 电子邮箱
	Office       string  // 办公地址
	Chairman     string  // 董事长
	Manager      string  // 总经理
	Secretary    string  // 董秘
	RegCapital   float64 // 注册资本
	SetupDate    string  // 成立日期
	Employees    int     // 员工人数

	// 状态
	Status      int  // 状态：1-正常，2-停牌，3-退市
	IsSuspended bool // 是否停牌
}

// GetQuoteDetail retrieves detailed quote information
// 获取个股详细信息
//
// API: /v5/stock/quote.json
// 方法: GET
// 认证: 需要Cookie或Token
//
// 参数说明：
//   - symbol: 股票代码，格式为 交易所+代码，如 SH600000、SZ000001
//   - extend: 是否获取扩展信息，默认true
//
// 返回字段：
//   - 基本行情：当前价、涨跌幅、成交量等
//   - 估值指标：PE、PB、市值等
//   - 公司信息：行业、上市日期、简介等
func (c *Client) GetQuoteDetail(ctx context.Context, params *QuoteDetailParams) (*QuoteDetailResult, *request.Record, error) {
	if params == nil || params.Symbol == "" {
		return nil, nil, fmt.Errorf("symbol required")
	}

	// 转换代码格式
	xueqiuSymbol, err := toXueqiuSymbol(params.Symbol)
	if err != nil {
		xueqiuSymbol = params.Symbol
	}

	query := url.Values{}
	query.Set("symbol", xueqiuSymbol)
	if params.Extend {
		query.Set("extend", "detail")
	} else {
		query.Set("extend", "")
	}

	reqURL := fmt.Sprintf("%s%s?%s", BaseURL, QuoteDetailAPI, query.Encode())

	req := request.Request{
		Method:  "GET",
		URL:     reqURL,
		Headers: c.buildHeaders(),
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseQuoteDetailResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

// GetProfile is an alias for GetQuoteDetail for compatibility
func (c *Client) GetProfile(ctx context.Context, params *QuoteDetailParams) (*QuoteDetailResult, *request.Record, error) {
	return c.GetQuoteDetail(ctx, params)
}

func parseQuoteDetailResponse(body []byte) (*QuoteDetailResult, error) {
	var raw struct {
		Data struct {
			Quote struct {
				Symbol         string      `json:"symbol"`
				Code           string      `json:"code"`
				Name           string      `json:"name"`
				Exchange       string      `json:"exchange"`
				Type           int         `json:"type"`
				Current        interface{} `json:"current"`
				Percent        interface{} `json:"percent"`
				Change         interface{} `json:"chg"`
				High           interface{} `json:"high"`
				Low            interface{} `json:"low"`
				Open           interface{} `json:"open"`
				PreClose       interface{} `json:"last_close"`
				Volume         interface{} `json:"volume"`
				Amount         interface{} `json:"amount"`
				TurnoverRate   interface{} `json:"turnover_rate"`
				Amplitude      interface{} `json:"amplitude"`
				PE             interface{} `json:"pe_ttm"`
				PB             interface{} `json:"pb"`
				PS             interface{} `json:"ps_ttm"`
				PCF            interface{} `json:"pcf_ratio"`
				TotalMarketCap interface{} `json:"market_capital"`
				FloatMarketCap interface{} `json:"float_market_capital"`
				DividendYield  interface{} `json:"dividend_yield"`
				TotalShares    interface{} `json:"total_shares"`
				FloatShares    interface{} `json:"float_shares"`
				Industry       string      `json:"industry"`
				Sector         string      `json:"sector"`
				ListDate       string      `json:"list_date"`
				Province       string      `json:"province"`
				City           string      `json:"city"`
				MainBusiness   string      `json:"main_business"`
				Description    string      `json:"description"`
				Website        string      `json:"website"`
				Email          string      `json:"email"`
				Office         string      `json:"office"`
				Chairman       string      `json:"chairman"`
				Manager        string      `json:"manager"`
				Secretary      string      `json:"secretary"`
				RegCapital     interface{} `json:"reg_capital"`
				SetupDate      string      `json:"setup_date"`
				Employees      interface{} `json:"employees"`
				Status         int         `json:"status"`
				IsSuspended    bool        `json:"is_suspended"`
			} `json:"quote"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	q := raw.Data.Quote
	return &QuoteDetailResult{
		Symbol:         q.Symbol,
		Code:           q.Code,
		Name:           q.Name,
		Exchange:       q.Exchange,
		Type:           q.Type,
		Current:        toFloat(q.Current),
		Percent:        toFloat(q.Percent),
		Change:         toFloat(q.Change),
		High:           toFloat(q.High),
		Low:            toFloat(q.Low),
		Open:           toFloat(q.Open),
		PreClose:       toFloat(q.PreClose),
		Volume:         toFloat(q.Volume),
		Amount:         toFloat(q.Amount),
		TurnoverRate:   toFloat(q.TurnoverRate),
		Amplitude:      toFloat(q.Amplitude),
		PE:             toFloat(q.PE),
		PB:             toFloat(q.PB),
		PS:             toFloat(q.PS),
		PCF:            toFloat(q.PCF),
		TotalMarketCap: toFloat(q.TotalMarketCap),
		FloatMarketCap: toFloat(q.FloatMarketCap),
		DividendYield:  toFloat(q.DividendYield),
		TotalShares:    toFloat(q.TotalShares),
		FloatShares:    toFloat(q.FloatShares),
		Industry:       q.Industry,
		Sector:         q.Sector,
		ListDate:       q.ListDate,
		Province:       q.Province,
		City:           q.City,
		MainBusiness:   q.MainBusiness,
		Description:    q.Description,
		Website:        q.Website,
		Email:          q.Email,
		Office:         q.Office,
		Chairman:       q.Chairman,
		Manager:        q.Manager,
		Secretary:      q.Secretary,
		RegCapital:     toFloat(q.RegCapital),
		SetupDate:      q.SetupDate,
		Employees:      toInt(q.Employees),
		Status:         q.Status,
		IsSuspended:    q.IsSuspended,
	}, nil
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		i, _ := strconv.Atoi(val)
		return i
	default:
		return 0
	}
}
