// Package tushare provides a Go client for the Tushare Pro API.
//
// Tushare Pro (https://tushare.pro) 是一个面向中国金融市场的数据服务平台,
// 提供股票、基金、期货等多种金融数据。
//
// 环境变量:
//   - TUSHARE_TOKEN: API 访问令牌（必需）
//   - TUSHARE_BASE_URL: 自定义 API 地址（可选，默认 http://api.tushare.pro）
package tushare

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/souloss/quantds/request"
)

// DefaultBaseURL 是 Tushare Pro 的默认 API 地址。
const DefaultBaseURL = "http://api.tushare.pro"

// API 名称常量
const (
	APIStockBasic   = "stock_basic"   // 股票基本信息
	APIDaily        = "daily"         // 日线行情（未复权）
	APIWeekly       = "weekly"        // 周线行情
	APIMonthly      = "monthly"       // 月线行情
	APIDailyBasic   = "daily_basic"   // 每日指标（PE/PB/PS/换手率/市值等）
	APIAdjFactor    = "adj_factor"    // 复权因子
	APIStockCompany = "stock_company" // 上市公司基本信息
	APIIncome       = "income"        // 利润表
	APIBalanceSheet = "balancesheet"  // 资产负债表
	APICashflow     = "cashflow"      // 现金流量表
	APIFinaIndicator = "fina_indicator" // 财务指标
	APITradeCal     = "trade_cal"     // 交易日历
	APIDividend     = "dividend"      // 分红送股
	APIConcept      = "concept"       // 概念分类
	APIConceptDetail = "concept_detail" // 概念明细
)

// 字段常量定义各 API 请求的返回字段列表
const (
	FieldsStockBasic = "ts_code,symbol,name,area,industry,market,exchange,list_status,list_date,delist_date"
	FieldsDaily      = "ts_code,trade_date,open,high,low,close,pre_close,change,pct_chg,vol,amount"
	FieldsDailyBasic = "ts_code,trade_date,turnover_rate,turnover_rate_f,volume_ratio,pe,pe_ttm,pb,ps,ps_ttm,dv_ratio,dv_ttm,total_share,float_share,total_mv,circ_mv"
	FieldsAdjFactor  = "ts_code,trade_date,adj_factor"
	FieldsIncome     = "ts_code,ann_date,f_ann_date,end_date,report_type,comp_type,basic_eps,diluted_eps,total_revenue,revenue,total_cogs,oper_cost,sell_exp,admin_exp,fin_exp,rd_exp,oper_profit,total_profit,n_income,n_income_attr_p,ebit,ebitda"
	FieldsBalance    = "ts_code,ann_date,f_ann_date,end_date,report_type,comp_type,total_assets,total_cur_assets,total_nca,total_liab,total_cur_liab,total_ncl,total_hldr_eqy_exc_min_int,total_hldr_eqy_inc_min_int,cap_rese,surplus_rese,undist_profit,money_cap,accounts_receiv,inventories,fix_assets"
	FieldsCashflow   = "ts_code,ann_date,f_ann_date,end_date,report_type,comp_type,n_cashflow_act,n_cashflow_inv_act,n_cash_flows_fnc_act,c_cash_equ_end_period"
	FieldsFinaInd    = "ts_code,ann_date,end_date,roe,roe_waa,roa,netprofit_margin,grossprofit_margin,current_ratio,quick_ratio,debt_to_assets,turn_days,roa_yearly,roe_avg,assets_turn,op_income,ebit,ebitda"
	FieldsTradeCal   = "exchange,cal_date,is_open,pretrade_date"
	FieldsCompany    = "ts_code,chairman,manager,secretary,reg_capital,setup_date,province,city,introduction,website,email,office,employees,main_business,business_scope"
	FieldsDividend   = "ts_code,ann_date,div_proc,stk_div,stk_bo_rate,stk_co_rate,cash_div,cash_div_tax,record_date,ex_date,pay_date,div_listdate,imp_ann_date"
	FieldsConcept    = "code,name,src"
	FieldsConceptDetail = "id,concept_name,ts_code,name,in_date,out_date"
)

// Client 是 Tushare Pro API 客户端。
type Client struct {
	http    request.Client
	token   string
	baseURL string
}

// Option 定义客户端配置选项。
type Option func(*Client)

// WithToken 设置 API 访问令牌。
func WithToken(token string) Option {
	return func(c *Client) { c.token = token }
}

// WithBaseURL 设置自定义 API 地址。
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// NewClient 创建 Tushare 客户端。
// 自动从环境变量 TUSHARE_TOKEN 和 TUSHARE_BASE_URL 读取配置。
func NewClient(httpClient request.Client, opts ...Option) *Client {
	if httpClient == nil {
		httpClient = request.NewClient(request.DefaultConfig())
	}

	baseURL := os.Getenv("TUSHARE_BASE_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	c := &Client{
		http:    httpClient,
		token:   os.Getenv("TUSHARE_TOKEN"),
		baseURL: baseURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Token 返回当前 Token
func (c *Client) Token() string {
	return c.token
}

// Close 关闭客户端。
func (c *Client) Close() {
	c.http.Close()
}

// FormatDate 格式化日期 YYYYMMDD -> YYYY-MM-DD
func FormatDate(date string) string {
	if len(date) != 8 {
		return date
	}
	return date[:4] + "-" + date[4:6] + "-" + date[6:]
}

// apiRequest 是 Tushare API 请求体。
type apiRequest struct {
	APIName string            `json:"api_name"`
	Token   string            `json:"token"`
	Params  map[string]string `json:"params,omitempty"`
	Fields  string            `json:"fields,omitempty"`
}

// apiResponse 是 Tushare API 响应体。
type apiResponse struct {
	Code int              `json:"code"` // 0=成功, 其他=错误
	Msg  string           `json:"msg"`
	Data *apiResponseData `json:"data"`
}

// apiResponseData 是 Tushare API 响应数据。
type apiResponseData struct {
	Fields []string        `json:"fields"` // 字段名列表
	Items  [][]interface{} `json:"items"`  // 数据行（二维数组）
}

// post 发送 Tushare API 请求。
func (c *Client) post(ctx context.Context, apiName string, params map[string]string, fields string) (*apiResponseData, *request.Record, error) {
	if c.token == "" {
		return nil, nil, fmt.Errorf("tushare token is required (set TUSHARE_TOKEN env or use WithToken)")
	}

	reqBody := apiRequest{
		APIName: apiName,
		Token:   c.token,
		Params:  params,
		Fields:  fields,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal request: %w", err)
	}

	req := request.Request{
		Method: "POST",
		URL:    c.baseURL,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Body: jsonBytes,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("http status %d", resp.StatusCode)
	}

	var tsResp apiResponse
	if err := json.Unmarshal(resp.Body, &tsResp); err != nil {
		return nil, record, fmt.Errorf("decode response: %w", err)
	}

	if tsResp.Code != 0 {
		return nil, record, fmt.Errorf("tushare api error: %s (code %d)", tsResp.Msg, tsResp.Code)
	}

	if tsResp.Data == nil {
		return nil, record, fmt.Errorf("empty data response")
	}

	return tsResp.Data, record, nil
}

// fieldIndex 将字段名列表转换为 name→index 映射，用于快速查找。
func fieldIndex(fields []string) map[string]int {
	m := make(map[string]int, len(fields))
	for i, f := range fields {
		m[f] = i
	}
	return m
}

// getStr 从响应行中按字段名提取字符串值。
func getStr(idx map[string]int, row []interface{}, key string) string {
	if i, ok := idx[key]; ok && i >= 0 && i < len(row) {
		if v, ok := row[i].(string); ok {
			return v
		}
	}
	return ""
}

// getFlt 从响应行中按字段名提取浮点数值。
func getFlt(idx map[string]int, row []interface{}, key string) float64 {
	if i, ok := idx[key]; ok && i >= 0 && i < len(row) {
		switch v := row[i].(type) {
		case float64:
			return v
		case string:
			f, _ := strconv.ParseFloat(v, 64)
			return f
		}
	}
	return 0
}

// getInt 从响应行中按字段名提取整数值。
func getInt(idx map[string]int, row []interface{}, key string) int {
	if i, ok := idx[key]; ok && i >= 0 && i < len(row) {
		switch v := row[i].(type) {
		case float64:
			return int(v)
		case string:
			n, _ := strconv.Atoi(v)
			return n
		}
	}
	return 0
}

// intToStr 将整数转换为字符串（避免使用 strconv 的简单实现）。
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	return strconv.Itoa(n)
}
