package finnhub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

// FinancialsParams represents parameters for financial data request
type FinancialsParams struct {
	Symbol    string // Stock symbol (e.g., "AAPL")
	Statement string // Statement type: "bs" (balance sheet), "ic" (income statement), "cf" (cash flow)
	Frequency string // Frequency: "annual" or "quarterly"
}

// FinancialsResult represents the financial data result
type FinancialsResult struct {
	Symbol    string
	Statement string
	Data      []FinancialDataPoint
}

// FinancialDataPoint represents a single financial data record
type FinancialDataPoint struct {
	Year      int
	Quarter   string // "Q1", "Q2", "Q3", "Q4", or "" for annual
	StartDate string
	EndDate   string
	Items     map[string]float64 // Key-value pairs of financial items
}

// GetFinancials retrieves financial statement data
func (c *Client) GetFinancials(ctx context.Context, params *FinancialsParams) (*FinancialsResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s?symbol=%s&statement=%s&frequency=%s&token=%s",
		BaseURL, FinancialsAPI, params.Symbol, params.Statement, params.Frequency, c.apiKey)

	req := request.Request{
		Method:  "GET",
		URL:     url,
		Headers: DefaultHeaders,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseFinancialsResponse(resp.Body, params.Symbol, params.Statement)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type finnhubFinancialsResponse struct {
	Symbol    string                     `json:"symbol"`
	Statement string                     `json:"statement"`
	Frequency string                     `json:"frequency"`
	Data      []finnhubFinancialDataItem `json:"data"`
}

type finnhubFinancialDataItem struct {
	Year      int                `json:"year"`
	Quarter   string             `json:"quarter"`
	StartDate string             `json:"startDate"`
	EndDate   string             `json:"endDate"`
	Items     map[string]float64 `json:"items"`
}

func parseFinancialsResponse(body []byte, symbol string, statement string) (*FinancialsResult, error) {
	var resp finnhubFinancialsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	dataPoints := make([]FinancialDataPoint, 0, len(resp.Data))
	for _, d := range resp.Data {
		dataPoints = append(dataPoints, FinancialDataPoint{
			Year:      d.Year,
			Quarter:   d.Quarter,
			StartDate: d.StartDate,
			EndDate:   d.EndDate,
			Items:     d.Items,
		})
	}

	return &FinancialsResult{
		Symbol:    symbol,
		Statement: statement,
		Data:      dataPoints,
	}, nil
}
