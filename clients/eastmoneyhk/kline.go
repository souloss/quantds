package eastmoneyhk

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/souloss/quantds/request"
)

// KlineAPI endpoint
const KlineAPI = "/api/qt/stock/kline/get"

// KlineParams represents parameters for K-line data request
type KlineParams struct {
	Symbol    string // Stock symbol (e.g., "00700.HK", "00941.HK")
	StartDate string // Start date in format "YYYYMMDD"
	EndDate   string // End date in format "YYYYMMDD"
	Period    string // Period: 1,5,15,30,60,101(daily),102(weekly),103(monthly)
	Adjust    string // Adjustment: 0(none),1(qfq),2(hfq)
}

// KlineResult represents the K-line data result
type KlineResult struct {
	Symbol string      // Stock symbol
	Data   []KlineData // K-line data points
	Count  int         // Number of data points
}

// KlineData represents a single K-line (OHLCV) data point
type KlineData struct {
	Date         string    // Date string
	Open         float64   // Opening price
	High         float64   // Highest price
	Low          float64   // Lowest price
	Close        float64   // Closing price
	Volume       float64   // Trading volume
	Turnover     float64   // Trading turnover
	Amplitude    float64   // Price amplitude (%)
	ChangeRate   float64   // Change rate (%)
	Change       float64   // Price change
	TurnoverRate float64   // Turnover rate (%)
	Timestamp    time.Time // Parsed timestamp
}

// GetKline retrieves historical K-line data for a HK stock
func (c *Client) GetKline(ctx context.Context, params *KlineParams) (*KlineResult, *request.Record, error) {
	secid, err := toHKSecid(params.Symbol)
	if err != nil {
		return nil, nil, err
	}

	period := params.Period
	if period == "" {
		period = Period1d
	}

	adjust := params.Adjust
	if adjust == "" {
		adjust = "0"
	}

	// Fields for HK stocks
	fields := "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61"

	url := fmt.Sprintf("%s%s?secid=%s&fields1=f1,f2,f3,f4,f5,f6&fields2=%s&klt=%s&fqt=%s&beg=%s&end=%s",
		BaseURL, KlineAPI, secid, fields, period, adjust,
		params.StartDate, params.EndDate)

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

	result, err := parseKlineResponse(resp.Body, params.Symbol)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type klineResponse struct {
	Data *struct {
		Klines []string `json:"klines"`
	} `json:"data"`
}

func parseKlineResponse(body []byte, symbol string) (*KlineResult, error) {
	var resp klineResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Data == nil || len(resp.Data.Klines) == 0 {
		return &KlineResult{Symbol: symbol, Data: nil, Count: 0}, nil
	}

	data := make([]KlineData, 0, len(resp.Data.Klines))
	for _, line := range resp.Data.Klines {
		kline, err := parseKlineLine(line)
		if err != nil {
			continue
		}
		data = append(data, *kline)
	}

	return &KlineResult{
		Symbol: symbol,
		Data:   data,
		Count:  len(data),
	}, nil
}

func parseKlineLine(line string) (*KlineData, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 11 {
		return nil, fmt.Errorf("invalid kline line: %s", line)
	}

	open, _ := strconv.ParseFloat(parts[1], 64)
	close_, _ := strconv.ParseFloat(parts[2], 64)
	high, _ := strconv.ParseFloat(parts[3], 64)
	low, _ := strconv.ParseFloat(parts[4], 64)
	volume, _ := strconv.ParseFloat(parts[5], 64)
	turnover, _ := strconv.ParseFloat(parts[6], 64)
	amplitude, _ := strconv.ParseFloat(parts[7], 64)
	changeRate, _ := strconv.ParseFloat(parts[8], 64)
	change, _ := strconv.ParseFloat(parts[9], 64)
	turnoverRate, _ := strconv.ParseFloat(parts[10], 64)

	// Parse date
	dateStr := parts[0]
	timestamp, _ := parseHKDate(dateStr)

	return &KlineData{
		Date:         dateStr,
		Open:         open,
		High:         high,
		Low:          low,
		Close:        close_,
		Volume:       volume,
		Turnover:     turnover,
		Amplitude:    amplitude,
		ChangeRate:   changeRate,
		Change:       change,
		TurnoverRate: turnoverRate,
		Timestamp:    timestamp,
	}, nil
}

// toHKSecid converts symbol to EastMoney secid format for HK stocks
// HK stocks use format "116.CODE" where 116 is the market ID
func toHKSecid(symbol string) (string, error) {
	code, ok := ParseHKSymbol(symbol)
	if !ok {
		return "", fmt.Errorf("invalid HK stock symbol: %s", symbol)
	}
	return fmt.Sprintf("116.%s", code), nil
}

// ParseHKSymbol parses a HK stock symbol
// Accepts formats: 00700, 00700.HK, 0700.HKEX
func ParseHKSymbol(symbol string) (code string, ok bool) {
	symbol = strings.TrimSpace(symbol)

	// Handle CODE.MARKET format
	if strings.Contains(symbol, ".") {
		parts := strings.Split(symbol, ".")
		if len(parts) >= 2 {
			code = parts[0]
			// Ensure 5-digit code
			code = padHKCode(code)
			return code, len(code) == 5
		}
		return "", false
	}

	// Plain code
	code = padHKCode(symbol)
	if len(code) == 5 {
		return code, true
	}
	return "", false
}

// padHKCode pads HK stock code to 5 digits
func padHKCode(code string) string {
	// Remove leading zeros and re-pad
	code = strings.TrimLeft(code, "0")
	if code == "" {
		return "00000"
	}
	// Pad to 5 digits
	for len(code) < 5 {
		code = "0" + code
	}
	return code
}

// parseHKDate parses HK date string
func parseHKDate(dateStr string) (time.Time, error) {
	if strings.Contains(dateStr, " ") {
		return time.ParseInLocation("2006-01-02 15:04", dateStr, timeLoc)
	}
	return time.ParseInLocation("2006-01-02", dateStr, timeLoc)
}

var timeLoc, _ = time.LoadLocation("Asia/Hong_Kong")

// ToPeriod converts domain timeframe to EastMoney period code
func ToPeriod(tf string) string {
	switch tf {
	case "1m":
		return Period1m
	case "5m":
		return Period5m
	case "15m":
		return Period15m
	case "30m":
		return Period30m
	case "60m":
		return Period60m
	case "1d", "":
		return Period1d
	case "1w":
		return Period1w
	case "1M":
		return Period1M
	default:
		return Period1d
	}
}

// ToAdjust converts adjustment string to EastMoney adjust code
func ToAdjust(adj string) string {
	switch adj {
	case "qfq":
		return "1"
	case "hfq":
		return "2"
	default:
		return "0"
	}
}
