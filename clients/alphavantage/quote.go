package alphavantage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/souloss/quantds/request"
)

type QuoteParams struct {
	Symbol string
}

type QuoteResult struct {
	Symbol        string
	Open          float64
	High          float64
	Low           float64
	Price         float64
	Volume        float64
	PreviousClose float64
	Change        float64
	ChangePercent string
	LatestDay     string
}

func (c *Client) GetQuote(ctx context.Context, params *QuoteParams) (*QuoteResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		BaseURL, QueryAPI, params.Symbol, c.apiKey)

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

	result, err := parseGlobalQuoteResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type globalQuoteResponse struct {
	GlobalQuote map[string]string `json:"Global Quote"`
}

func parseGlobalQuoteResponse(body []byte) (*QuoteResult, error) {
	var resp globalQuoteResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	q := resp.GlobalQuote
	open, _ := strconv.ParseFloat(q["02. open"], 64)
	high, _ := strconv.ParseFloat(q["03. high"], 64)
	low, _ := strconv.ParseFloat(q["04. low"], 64)
	price, _ := strconv.ParseFloat(q["05. price"], 64)
	volume, _ := strconv.ParseFloat(q["06. volume"], 64)
	prevClose, _ := strconv.ParseFloat(q["08. previous close"], 64)
	change, _ := strconv.ParseFloat(q["09. change"], 64)

	return &QuoteResult{
		Symbol:        q["01. symbol"],
		Open:          open,
		High:          high,
		Low:           low,
		Price:         price,
		Volume:        volume,
		PreviousClose: prevClose,
		Change:        change,
		ChangePercent: q["10. change percent"],
		LatestDay:     q["07. latest trading day"],
	}, nil
}
