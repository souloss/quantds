package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/souloss/quantds/request"
)

const FinancialAPI = "/api/data/v1/get"

// Report type constants
const (
	ReportBalanceSheet = "RPT_DMSK_FN_BALANCE"
	ReportIncome       = "RPT_DMSK_FN_INCOME"
	ReportCashflow     = "RPT_DMSK_FN_CASHFLOW"
	ReportIndicator    = "RPT_DMSK_FN_FINANCE"
)

// FinancialParams represents parameters for financial data request
type FinancialParams struct {
	ReportName string
	Code       string
	PageNumber int
	PageSize   int
}

// FinancialResult represents the financial data result
type FinancialResult struct {
	Data    []map[string]interface{}
	Success bool
	Code    int
	Message string
}

// FinanceParams is an alias for FinancialParams
type FinanceParams = FinancialParams

// GetFinancials retrieves financial report data
func (c *Client) GetFinancials(ctx context.Context, params *FinancialParams) (*FinancialResult, *request.Record, error) {
	if params.Code == "" || params.ReportName == "" {
		return nil, nil, fmt.Errorf("code and reportName required")
	}

	v := url.Values{}
	v.Set("reportName", params.ReportName)
	v.Set("columns", "ALL")
	v.Set("filter", fmt.Sprintf(`(SECURITY_CODE="%s")`, params.Code))
	if params.PageNumber > 0 {
		v.Set("pageNumber", strconv.Itoa(params.PageNumber))
	} else {
		v.Set("pageNumber", "1")
	}
	if params.PageSize > 0 {
		v.Set("pageSize", strconv.Itoa(params.PageSize))
	} else {
		v.Set("pageSize", "50")
	}
	v.Set("sortColumns", "REPORT_DATE")
	v.Set("sortTypes", "-1")

	apiURL := Datacenter + FinancialAPI + "?" + v.Encode()

	req := request.Request{
		Method: "GET",
		URL:    apiURL,
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":    "https://data.eastmoney.com/",
		},
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseFinancialResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

// GetFinance is an alias for GetFinancials
func (c *Client) GetFinance(ctx context.Context, params *FinanceParams) (*FinancialResult, *request.Record, error) {
	return c.GetFinancials(ctx, params)
}

// GetBalanceSheet retrieves balance sheet data
func (c *Client) GetBalanceSheet(ctx context.Context, code string, pageNum, pageSize int) (*FinancialResult, *request.Record, error) {
	return c.GetFinancials(ctx, &FinancialParams{
		ReportName: ReportBalanceSheet,
		Code:       code,
		PageNumber: pageNum,
		PageSize:   pageSize,
	})
}

// GetIncomeStatement retrieves income statement data
func (c *Client) GetIncomeStatement(ctx context.Context, code string, pageNum, pageSize int) (*FinancialResult, *request.Record, error) {
	return c.GetFinancials(ctx, &FinancialParams{
		ReportName: ReportIncome,
		Code:       code,
		PageNumber: pageNum,
		PageSize:   pageSize,
	})
}

// GetCashflowStatement retrieves cash flow statement data
func (c *Client) GetCashflowStatement(ctx context.Context, code string, pageNum, pageSize int) (*FinancialResult, *request.Record, error) {
	return c.GetFinancials(ctx, &FinancialParams{
		ReportName: ReportCashflow,
		Code:       code,
		PageNumber: pageNum,
		PageSize:   pageSize,
	})
}

// GetFinancialIndicator retrieves financial indicator data
func (c *Client) GetFinancialIndicator(ctx context.Context, code string, pageNum, pageSize int) (*FinancialResult, *request.Record, error) {
	return c.GetFinancials(ctx, &FinancialParams{
		ReportName: ReportIndicator,
		Code:       code,
		PageNumber: pageNum,
		PageSize:   pageSize,
	})
}

type financialResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Data []map[string]interface{} `json:"data"`
	} `json:"result"`
}

func parseFinancialResponse(body []byte) (*FinancialResult, error) {
	var resp financialResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if !resp.Success || resp.Result.Data == nil {
		if resp.Code == 9201 {
			return &FinancialResult{Data: []map[string]interface{}{}, Success: true}, nil
		}
		return &FinancialResult{
			Success: resp.Success,
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}

	return &FinancialResult{
		Data:    resp.Result.Data,
		Success: resp.Success,
		Code:    resp.Code,
		Message: resp.Message,
	}, nil
}
