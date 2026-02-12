package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/souloss/quantds/request"
)

const CandleAPI = "/api/qt/stock/kline/get"

const FieldCandles = "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61"

// CandleParams represents parameters for candlestick data request
type CandleParams struct {
	Symbol    string // Stock symbol (e.g., "000001.SZ")
	StartDate string // Start date in format "YYYYMMDD"
	EndDate   string // End date in format "YYYYMMDD"
	Period    string // Period: 1,5,15,30,60,101(daily),102(weekly),103(monthly)
	Adjust    string // Adjustment: 0(none),1(qfq),2(hfq)
}

// CandleResult represents the candlestick data result
type CandleResult struct {
	Data  []CandleData
	Count int
}

// CandleData represents a single candlestick/OHLCV data point
type CandleData struct {
	Date         string  // Date string
	Open         float64 // Opening price
	Close        float64 // Closing price
	High         float64 // Highest price
	Low          float64 // Lowest price
	Volume       float64 // Trading volume
	Turnover     float64 // Trading turnover
	Amplitude    float64 // Price amplitude (%)
	ChangeRate   float64 // Change rate (%)
	Change       float64 // Price change
	TurnoverRate float64 // Turnover rate (%)
}

// KlineBar is an alias for CandleData for backward compatibility
type KlineBar = CandleData

// KlineParams is an alias for CandleParams for backward compatibility
type KlineParams = CandleParams

// KlineResult is an alias for CandleResult for backward compatibility
type KlineResult = CandleResult

// GetCandles retrieves historical candlestick data for a symbol
func (c *Client) GetCandles(ctx context.Context, params *CandleParams) (*CandleResult, *request.Record, error) {
	secid, err := toEastMoneySecid(params.Symbol)
	if err != nil {
		return nil, nil, err
	}

	period := params.Period
	if period == "" {
		period = "101"
	}

	adjust := params.Adjust
	if adjust == "" {
		adjust = "0"
	}

	url := fmt.Sprintf("%s%s?secid=%s&fields1=f1,f2,f3,f4,f5,f6&fields2=%s&klt=%s&fqt=%s&beg=%s&end=%s",
		BaseURL, CandleAPI, secid, FieldCandles, period, adjust,
		params.StartDate, params.EndDate)

	req := request.Request{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":    "https://quote.eastmoney.com/",
		},
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseCandleResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

// GetKline is an alias for GetCandles for backward compatibility
func (c *Client) GetKline(ctx context.Context, params *KlineParams) (*KlineResult, *request.Record, error) {
	return c.GetCandles(ctx, params)
}

type candleResponse struct {
	Data *struct {
		Klines []string `json:"klines"`
	} `json:"data"`
}

func parseCandleResponse(body []byte) (*CandleResult, error) {
	var resp candleResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Data == nil || len(resp.Data.Klines) == 0 {
		return &CandleResult{Data: nil, Count: 0}, nil
	}

	candles := make([]CandleData, 0, len(resp.Data.Klines))
	for _, line := range resp.Data.Klines {
		candle, err := parseCandleLine(line)
		if err != nil {
			continue
		}
		candles = append(candles, *candle)
	}

	return &CandleResult{Data: candles, Count: len(candles)}, nil
}

// parseKlineResponse is an alias for parseCandleResponse
func parseKlineResponse(body []byte) (*KlineResult, error) {
	return parseCandleResponse(body)
}

func parseCandleLine(line string) (*CandleData, error) {
	parts := strings.Split(line, ",")
	if len(parts) < 11 {
		return nil, fmt.Errorf("invalid candle line: %s", line)
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

	return &CandleData{
		Date:         parts[0],
		Open:         open,
		Close:        close_,
		High:         high,
		Low:          low,
		Volume:       volume,
		Turnover:     turnover,
		Amplitude:    amplitude,
		ChangeRate:   changeRate,
		Change:       change,
		TurnoverRate: turnoverRate,
	}, nil
}

// parseKlineLine is an alias for parseCandleLine
func parseKlineLine(line string) (*KlineBar, error) {
	return parseCandleLine(line)
}

// Period conversion functions

// ToPeriod converts timeframe string to EastMoney period code
func ToPeriod(tf string) string {
	switch tf {
	case "1m":
		return "1"
	case "5m":
		return "5"
	case "15m":
		return "15"
	case "30m":
		return "30"
	case "60m":
		return "60"
	case "1d", "":
		return "101"
	case "1w":
		return "102"
	case "1M":
		return "103"
	default:
		return "101"
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

// ParseDate parses date string to time.Time
func ParseDate(dateStr string) time.Time {
	if strings.Contains(dateStr, " ") {
		t, _ := time.ParseInLocation("2006-01-02 15:04", dateStr, timeLoc)
		return t
	}
	t, _ := time.ParseInLocation("2006-01-02", dateStr, timeLoc)
	return t
}

var timeLoc, _ = time.LoadLocation("Asia/Shanghai")
