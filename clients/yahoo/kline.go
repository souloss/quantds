package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/souloss/quantds/request"
)

// KlineParams represents parameters for K-line data request
type KlineParams struct {
	Symbol    string    // Stock symbol (e.g., "AAPL", "MSFT")
	Interval  string    // K-line interval: 1m, 5m, 15m, 30m, 60m, 1d, 1wk, 1mo
	Range     string    // Time range: 1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max
	StartDate time.Time // Start date (optional, used with EndDate)
	EndDate   time.Time // End date (optional, used with StartDate)
}

// KlineResult represents the K-line data result
type KlineResult struct {
	Symbol   string      // Stock symbol
	Timezone string      // Timezone
	Data     []KlineData // K-line data points
	Count    int         // Number of data points
}

// KlineData represents a single K-line (OHLCV) data point
type KlineData struct {
	Timestamp int64   // Unix timestamp in seconds
	Open      float64 // Opening price
	High      float64 // Highest price
	Low       float64 // Lowest price
	Close     float64 // Closing price
	Volume    float64 // Trading volume
}

// GetKline retrieves historical K-line data for a symbol
func (c *Client) GetKline(ctx context.Context, params *KlineParams) (*KlineResult, *request.Record, error) {
	if params.Interval == "" {
		params.Interval = Interval1d
	}
	if params.Range == "" && params.StartDate.IsZero() {
		params.Range = Range1y
	}

	url := fmt.Sprintf("%s%s/%s", BaseURL, ChartAPI, params.Symbol)

	query := fmt.Sprintf("?interval=%s", params.Interval)

	if params.Range != "" {
		query += fmt.Sprintf("&range=%s", params.Range)
	} else if !params.StartDate.IsZero() && !params.EndDate.IsZero() {
		// Convert to Unix timestamp
		period1 := params.StartDate.Unix()
		period2 := params.EndDate.Unix()
		query += fmt.Sprintf("&period1=%d&period2=%d", period1, period2)
	}

	req := request.Request{
		Method:  "GET",
		URL:     url + query,
		Headers: DefaultHeaders,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseKlineResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

// chartResponse represents the Yahoo Finance chart API response
type chartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol   string `json:"symbol"`
				Timezone string `json:"exchangeTimezoneName"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []float64 `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

func parseKlineResponse(body []byte) (*KlineResult, error) {
	var resp chartResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if len(resp.Chart.Result) == 0 {
		return &KlineResult{}, nil
	}

	data := resp.Chart.Result[0]
	if len(data.Timestamp) == 0 {
		return &KlineResult{
			Symbol:   data.Meta.Symbol,
			Timezone: data.Meta.Timezone,
		}, nil
	}

	var quotes []KlineData
	if len(data.Indicators.Quote) > 0 {
		q := data.Indicators.Quote[0]
		for i, ts := range data.Timestamp {
			if i >= len(q.Open) {
				break
			}
			quotes = append(quotes, KlineData{
				Timestamp: ts,
				Open:      q.Open[i],
				High:      q.High[i],
				Low:       q.Low[i],
				Close:     q.Close[i],
				Volume:    q.Volume[i],
			})
		}
	}

	return &KlineResult{
		Symbol:   data.Meta.Symbol,
		Timezone: data.Meta.Timezone,
		Data:     quotes,
		Count:    len(quotes),
	}, nil
}

// ToInterval converts domain timeframe to Yahoo interval
func ToInterval(tf string) string {
	switch tf {
	case "1m":
		return Interval1m
	case "5m":
		return Interval5m
	case "15m":
		return Interval15m
	case "30m":
		return Interval30m
	case "60m":
		return Interval60m
	case "1d", "":
		return Interval1d
	case "1w":
		return Interval1w
	case "1M":
		return Interval1M
	default:
		return Interval1d
	}
}

// ParseTimestamp converts Unix timestamp to time.Time
func ParseTimestamp(ts int64, tz string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	return time.Unix(ts, 0).In(loc)
}

// ParseTimestampString converts Unix timestamp string to time.Time
func ParseTimestampString(tsStr string, tz string) (time.Time, error) {
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return ParseTimestamp(ts, tz), nil
}
