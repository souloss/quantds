package alphavantage

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/souloss/quantds/request"
)

type KlineParams struct {
	Symbol   string
	Function string // TIME_SERIES_DAILY, TIME_SERIES_WEEKLY, TIME_SERIES_MONTHLY
	Size     string // compact (100) or full (20+ years)
}

type KlineResult struct {
	Symbol   string
	Timezone string
	Data     []KlineData
	Count    int
}

type KlineData struct {
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func (c *Client) GetDailyTimeSeries(ctx context.Context, params *KlineParams) (*KlineResult, *request.Record, error) {
	fn := params.Function
	if fn == "" {
		fn = "TIME_SERIES_DAILY"
	}
	size := params.Size
	if size == "" {
		size = "compact"
	}

	url := fmt.Sprintf("%s%s?function=%s&symbol=%s&outputsize=%s&apikey=%s",
		BaseURL, QueryAPI, fn, params.Symbol, size, c.apiKey)

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

	result, err := parseDailyResponse(resp.Body, params.Symbol)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

func parseDailyResponse(body []byte, symbol string) (*KlineResult, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	if errMsg, ok := raw["Error Message"]; ok {
		return nil, fmt.Errorf("API error: %s", string(errMsg))
	}
	if note, ok := raw["Note"]; ok {
		return nil, fmt.Errorf("API rate limit: %s", string(note))
	}

	var tsKey string
	for k := range raw {
		if strings.Contains(k, "Time Series") {
			tsKey = k
			break
		}
	}
	if tsKey == "" {
		return &KlineResult{Symbol: symbol}, nil
	}

	var timeSeries map[string]map[string]string
	if err := json.Unmarshal(raw[tsKey], &timeSeries); err != nil {
		return nil, err
	}

	dates := make([]string, 0, len(timeSeries))
	for d := range timeSeries {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	data := make([]KlineData, 0, len(dates))
	for _, date := range dates {
		vals := timeSeries[date]
		open, _ := strconv.ParseFloat(vals["1. open"], 64)
		high, _ := strconv.ParseFloat(vals["2. high"], 64)
		low, _ := strconv.ParseFloat(vals["3. low"], 64)
		close_, _ := strconv.ParseFloat(vals["4. close"], 64)
		volume, _ := strconv.ParseFloat(vals["5. volume"], 64)
		data = append(data, KlineData{
			Date: date, Open: open, High: high, Low: low, Close: close_, Volume: volume,
		})
	}

	return &KlineResult{
		Symbol: symbol,
		Data:   data,
		Count:  len(data),
	}, nil
}
