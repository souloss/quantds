package alphavantage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/souloss/quantds/request"
)

type ForexRateParams struct {
	FromCurrency string
	ToCurrency   string
}

type ForexRateResult struct {
	FromCurrency string
	ToCurrency   string
	ExchangeRate float64
	BidPrice     float64
	AskPrice     float64
	LastRefreshed string
}

func (c *Client) GetForexExchangeRate(ctx context.Context, params *ForexRateParams) (*ForexRateResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s?function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s&apikey=%s",
		BaseURL, QueryAPI, params.FromCurrency, params.ToCurrency, c.apiKey)

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

	result, err := parseForexRateResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type forexRateResponse struct {
	Rate map[string]string `json:"Realtime Currency Exchange Rate"`
}

func parseForexRateResponse(body []byte) (*ForexRateResult, error) {
	var resp forexRateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	r := resp.Rate
	rate, _ := strconv.ParseFloat(r["5. Exchange Rate"], 64)
	bid, _ := strconv.ParseFloat(r["8. Bid Price"], 64)
	ask, _ := strconv.ParseFloat(r["9. Ask Price"], 64)

	return &ForexRateResult{
		FromCurrency:  r["1. From_Currency Code"],
		ToCurrency:    r["3. To_Currency Code"],
		ExchangeRate:  rate,
		BidPrice:      bid,
		AskPrice:      ask,
		LastRefreshed: r["6. Last Refreshed"],
	}, nil
}
