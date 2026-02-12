package xueqiu

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/souloss/quantds/request"
)

const (
	KlineAPI = "/chart/kline.json"
)

type KlineParams struct {
	Symbol string
	Period string
	Count  int
}

type KlineResult struct {
	Data  []KlineBar
	Count int
}

type KlineBar struct {
	Timestamp     int64
	Volume        float64
	Open          float64
	High          float64
	Low           float64
	Close         float64
	Change        float64
	ChangePercent float64
	Turnover      float64
}

type xueqiuKlineResponse struct {
	DataList  [][]interface{} `json:"dataList"`
	ChartList struct {
		Timestamp int64   `json:"timestamp"`
		Volume    float64 `json:"volume"`
		Open      float64 `json:"open"`
		High      float64 `json:"high"`
		Low       float64 `json:"low"`
		Close     float64 `json:"close"`
		Chg       float64 `json:"chg"`
		Percent   float64 `json:"percent"`
		Turnover  float64 `json:"turnover"`
	} `json:"chartList"`
}

func (c *Client) GetKline(ctx context.Context, params *KlineParams) (*KlineResult, *request.Record, error) {
	symbol, err := toXueqiuSymbol(params.Symbol)
	if err != nil {
		return nil, nil, err
	}

	period := params.Period
	if period == "" {
		period = "day"
	}

	count := params.Count
	if count == 0 {
		count = 320
	}

	url := fmt.Sprintf("%s%s?symbol=%s&type=%s&count=%d",
		BaseURL, KlineAPI, symbol, period, count)

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Referer":    "https://xueqiu.com/",
	}
	if c.cookie != "" {
		headers["Cookie"] = c.cookie
	}
	if c.token != "" {
		headers["X-Token"] = c.token
	}

	req := request.Request{
		Method:  "GET",
		URL:     url,
		Headers: headers,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseXueqiuKline(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

func parseXueqiuKline(body []byte) (*KlineResult, error) {
	var rawResp struct {
		DataList [][]interface{} `json:"dataList"`
	}

	if err := json.Unmarshal(body, &rawResp); err != nil {
		return nil, err
	}

	if len(rawResp.DataList) == 0 {
		return &KlineResult{Data: nil, Count: 0}, nil
	}

	bars := make([]KlineBar, 0, len(rawResp.DataList))
	for _, raw := range rawResp.DataList {
		if len(raw) < 8 {
			continue
		}
		bar := KlineBar{
			Timestamp:     int64(getFloat(raw[0])),
			Volume:        getFloat(raw[1]),
			Open:          getFloat(raw[2]),
			High:          getFloat(raw[3]),
			Low:           getFloat(raw[4]),
			Close:         getFloat(raw[5]),
			Change:        getFloat(raw[6]),
			ChangePercent: getFloat(raw[7]),
			Turnover:      getFloat(raw[8]),
		}
		bars = append(bars, bar)
	}

	return &KlineResult{Data: bars, Count: len(bars)}, nil
}

func getFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		return 0
	}
}

func toXueqiuSymbol(symbol string) (string, error) {
	code, exchange, ok := parseSymbol(symbol)
	if !ok {
		return "", fmt.Errorf("invalid symbol: %s", symbol)
	}
	switch exchange {
	case "SH":
		return "SH" + code, nil
	case "SZ":
		return "SZ" + code, nil
	case "BJ":
		return "BJ" + code, nil
	default:
		return "SH" + code, nil
	}
}

func parseSymbol(symbol string) (code string, exchange string, ok bool) {
	if len(symbol) < 6 {
		return "", "", false
	}
	if strings.Contains(symbol, ".") {
		parts := strings.Split(symbol, ".")
		if len(parts) == 2 {
			return parts[0], strings.ToUpper(parts[1]), true
		}
	}
	return "", "", false
}

func ToPeriod(tf string) string {
	switch tf {
	case "1m":
		return "min1"
	case "5m":
		return "min5"
	case "15m":
		return "min15"
	case "30m":
		return "min30"
	case "60m":
		return "min60"
	case "1d", "":
		return "day"
	case "1w":
		return "week"
	case "1M":
		return "month"
	default:
		return "day"
	}
}

func ParseTimestamp(ts int64) time.Time {
	return time.Unix(ts/1000, 0).In(timeLoc)
}

var timeLoc, _ = time.LoadLocation("Asia/Shanghai")
