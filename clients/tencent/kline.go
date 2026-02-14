package tencent

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
	// KlineAPI is the endpoint for K-line data (not used directly, URL is constructed in GetKline)
	KlineAPI = "/q=cn_kline/d"
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
	Date   string
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume float64
}

type tencentKlineResponse struct {
	Data string `json:"data"`
}

func (c *Client) GetKline(ctx context.Context, params *KlineParams) (*KlineResult, *request.Record, error) {
	symbol, err := toTencentSymbol(params.Symbol)
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

	url := fmt.Sprintf("https://web.sqt.gtimg.cn/q=cn_%s/k%s=%s", period, period, symbol)

	req := request.Request{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"User-Agent": DefaultUserAgent,
			"Referer":    DefaultReferer,
		},
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseTencentKline(resp.Body)
	if err != nil {
		return nil, record, err
	}

	if count > 0 && len(result.Data) > count {
		result.Data = result.Data[len(result.Data)-count:]
		result.Count = len(result.Data)
	}

	return result, record, nil
}

func parseTencentKline(body []byte) (*KlineResult, error) {
	bodyStr := string(body)
	startIdx := strings.Index(bodyStr, "kline:")
	if startIdx == -1 {
		return &KlineResult{Data: nil, Count: 0}, nil
	}

	jsonStart := strings.Index(bodyStr[startIdx:], "[")
	if jsonStart == -1 {
		return &KlineResult{Data: nil, Count: 0}, nil
	}

	jsonEnd := strings.LastIndex(bodyStr, "]")
	if jsonEnd == -1 {
		return &KlineResult{Data: nil, Count: 0}, nil
	}

	jsonStr := bodyStr[startIdx+jsonStart : jsonEnd+1]

	var rawBars [][]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &rawBars); err != nil {
		return nil, err
	}

	bars := make([]KlineBar, 0, len(rawBars))
	for _, raw := range rawBars {
		if len(raw) < 6 {
			continue
		}
		bar := KlineBar{
			Date:   fmt.Sprintf("%v", raw[0]),
			Open:   parseFloat(fmt.Sprintf("%v", raw[1])),
			Close:  parseFloat(fmt.Sprintf("%v", raw[2])),
			High:   parseFloat(fmt.Sprintf("%v", raw[3])),
			Low:    parseFloat(fmt.Sprintf("%v", raw[4])),
			Volume: parseFloat(fmt.Sprintf("%v", raw[5])),
		}
		bars = append(bars, bar)
	}

	return &KlineResult{Data: bars, Count: len(bars)}, nil
}

func toTencentSymbol(symbol string) (string, error) {
	code, exchange, ok := parseSymbol(symbol)
	if !ok {
		return "", fmt.Errorf("invalid symbol: %s", symbol)
	}
	switch exchange {
	case "SH":
		return "sh" + code, nil
	case "SZ":
		return "sz" + code, nil
	case "BJ":
		return "bj" + code, nil
	default:
		return "sh" + code, nil
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

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func ToPeriod(tf string) string {
	switch tf {
	case "1m":
		return "m1"
	case "5m":
		return "m5"
	case "15m":
		return "m15"
	case "30m":
		return "m30"
	case "60m":
		return "m60"
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

func ParseDate(dateStr string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02", dateStr, timeLoc)
	return t
}

var timeLoc, _ = time.LoadLocation("Asia/Shanghai")
