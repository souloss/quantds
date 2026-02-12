package binance

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
	Symbol    string    // Trading pair symbol (e.g., "BTCUSDT")
	Interval  string    // K-line interval: 1m, 5m, 15m, 30m, 1h, 4h, 1d, 1w, 1M
	Limit     int       // Number of data points (max 1000, default 500)
	StartTime time.Time // Start time (optional)
	EndTime   time.Time // End time (optional)
}

// KlineResult represents the K-line data result
type KlineResult struct {
	Symbol string      // Trading pair symbol
	Data   []KlineData // K-line data points
	Count  int         // Number of data points
}

// KlineData represents a single K-line (OHLCV) data point
type KlineData struct {
	OpenTime      int64   // K-line open time (Unix timestamp in ms)
	Open          float64 // Opening price
	High          float64 // Highest price
	Low           float64 // Lowest price
	Close         float64 // Closing price
	Volume        float64 // Trading volume (in base asset)
	CloseTime     int64   // K-line close time (Unix timestamp in ms)
	QuoteVol      float64 // Trading volume (in quote asset)
	Trades        int     // Number of trades
	TakerBuyBase  float64 // Taker buy volume (base asset)
	TakerBuyQuote float64 // Taker buy volume (quote asset)
}

// GetKline retrieves historical K-line data for a trading pair
func (c *Client) GetKline(ctx context.Context, params *KlineParams) (*KlineResult, *request.Record, error) {
	if params.Interval == "" {
		params.Interval = Interval1d
	}
	if params.Limit <= 0 || params.Limit > MaxKlineLimit {
		params.Limit = 500
	}

	url := fmt.Sprintf("%s%s?symbol=%s&interval=%s&limit=%d",
		BaseURL, KlineAPI, params.Symbol, params.Interval, params.Limit)

	if !params.StartTime.IsZero() {
		url += fmt.Sprintf("&startTime=%d", params.StartTime.UnixMilli())
	}
	if !params.EndTime.IsZero() {
		url += fmt.Sprintf("&endTime=%d", params.EndTime.UnixMilli())
	}

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

// Binance returns klines as array of arrays
type binanceKlineResponse [][]interface{}

func parseKlineResponse(body []byte, symbol string) (*KlineResult, error) {
	var resp binanceKlineResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	data := make([]KlineData, 0, len(resp))
	for _, k := range resp {
		if len(k) < 11 {
			continue
		}

		kline := KlineData{
			OpenTime:      int64(k[0].(float64)),
			Open:          parseFloat(k[1]),
			High:          parseFloat(k[2]),
			Low:           parseFloat(k[3]),
			Close:         parseFloat(k[4]),
			Volume:        parseFloat(k[5]),
			CloseTime:     int64(k[6].(float64)),
			QuoteVol:      parseFloat(k[7]),
			Trades:        int(k[8].(float64)),
			TakerBuyBase:  parseFloat(k[9]),
			TakerBuyQuote: parseFloat(k[10]),
		}
		data = append(data, kline)
	}

	return &KlineResult{
		Symbol: symbol,
		Data:   data,
		Count:  len(data),
	}, nil
}

func parseFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}

// ToInterval converts domain timeframe to Binance interval
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
	case "60m", "1h":
		return Interval1h
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

// ParseOpenTime converts Unix millisecond timestamp to time.Time
func ParseOpenTime(ms int64) time.Time {
	return time.UnixMilli(ms)
}

// FormatToUnixMilli converts time.Time to Unix millisecond timestamp
func FormatToUnixMilli(t time.Time) int64 {
	return t.UnixMilli()
}
