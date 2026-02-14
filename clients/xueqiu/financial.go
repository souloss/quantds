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
	IncomeAPI    = "/v5/stock/finance/%s/income.json"
	BalanceAPI   = "/v5/stock/finance/%s/balance.json"
	CashFlowAPI  = "/v5/stock/finance/%s/cash_flow.json"
	IndicatorAPI = "/v5/stock/finance/%s/indicator.json"
)

// FinancialParams represents parameters for financial request
type FinancialParams struct {
	Symbol string
	Type   string // Q4=Annual, Q3=Q3, Q2=Interim, Q1=Q1, 0=All
	Count  int
}

// FinancialItem represents a generic financial data item
type FinancialItem struct {
	ReportName string                 `json:"report_name"`
	ReportDate int64                  `json:"report_date"`
	Values     map[string]interface{} `json:"values"`
}

// GetIncome retrieves income statement
func (c *Client) GetIncome(ctx context.Context, params *FinancialParams) ([]FinancialItem, *request.Record, error) {
	return c.getFinancial(ctx, IncomeAPI, params)
}

// GetBalance retrieves balance sheet
func (c *Client) GetBalance(ctx context.Context, params *FinancialParams) ([]FinancialItem, *request.Record, error) {
	return c.getFinancial(ctx, BalanceAPI, params)
}

// GetCashFlow retrieves cash flow statement
func (c *Client) GetCashFlow(ctx context.Context, params *FinancialParams) ([]FinancialItem, *request.Record, error) {
	return c.getFinancial(ctx, CashFlowAPI, params)
}

// GetIndicator retrieves financial indicators
func (c *Client) GetIndicator(ctx context.Context, params *FinancialParams) ([]FinancialItem, *request.Record, error) {
	return c.getFinancial(ctx, IndicatorAPI, params)
}

func (c *Client) getFinancial(ctx context.Context, apiTmpl string, params *FinancialParams) ([]FinancialItem, *request.Record, error) {
	if params.Symbol == "" {
		return nil, nil, fmt.Errorf("symbol required")
	}

	xueqiuSymbol, err := toXueqiuSymbol(params.Symbol)
	if err != nil {
		xueqiuSymbol = params.Symbol
	}

	api := fmt.Sprintf(apiTmpl, xueqiuSymbol)
	
	query := url.Values{}
	query.Set("symbol", xueqiuSymbol)
	if params.Type != "" {
		query.Set("type", params.Type)
	} else {
		query.Set("type", "all")
	}
	if params.Count > 0 {
		query.Set("count", strconv.Itoa(params.Count))
	} else {
		query.Set("count", "10")
	}
	query.Set("is_detail", "true")

	reqURL := fmt.Sprintf("%s%s?%s", BaseURL, api, query.Encode())

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

	items, err := parseFinancialResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}
	return items, record, nil
}

func parseFinancialResponse(body []byte) ([]FinancialItem, error) {
	// Since fields are dynamic, we need map[string]interface{} unmarshal for items.
	// But "list" is an array of objects.
	// We can't define struct easily.
	
	// Let's use map parsing.
	var rawMap map[string]interface{}
	if err := json.Unmarshal(body, &rawMap); err != nil {
		return nil, err
	}

	data, ok := rawMap["data"].(map[string]interface{})
	if !ok || data == nil {
		return nil, nil
	}

	list, ok := data["list"].([]interface{})
	if !ok {
		return nil, nil
	}

	items := make([]FinancialItem, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		fi := FinancialItem{
			Values: make(map[string]interface{}),
		}

		for k, v := range m {
			if k == "report_name" {
				if s, ok := v.(string); ok {
					fi.ReportName = s
				}
			} else if k == "report_date" {
				if f, ok := v.(float64); ok {
					fi.ReportDate = int64(f)
				}
			} else {
				// Value is usually [value, yoy]
				// We just store the value (first element) or the whole thing?
				// Let's store the raw array or just the value.
				// For simplicity, let's store the raw value for now.
				fi.Values[k] = v
			}
		}
		items = append(items, fi)
	}

	return items, nil
}
