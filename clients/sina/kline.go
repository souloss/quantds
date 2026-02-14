package sina

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
	KlineAPI = "/cn/api/json_v2.php/CN_MarketDataService.getKLineData"
)

type KlineParams struct {
	Symbol    string
	StartDate string
	EndDate   string
	Period    string
	Adjust    string
}

type KlineResult struct {
	Data  []KlineBar
	Count int
}

type KlineBar struct {
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

type sinaKlineResponse struct {
	Day []struct {
		D string `json:"d"`
		O string `json:"o"`
		H string `json:"h"`
		L string `json:"l"`
		C string `json:"c"`
		V string `json:"v"`
	} `json:"day,omitempty"`
}

func (c *Client) GetKline(ctx context.Context, params *KlineParams) (*KlineResult, *request.Record, error) {
	symbol, err := toSinaSymbol(params.Symbol)
	if err != nil {
		return nil, nil, err
	}

	period := params.Period
	if period == "" {
		period = "d"
	}

	u := fmt.Sprintf("%s%s?symbol=%s&scale=%s&datalen=%d",
		BaseURL, KlineAPI, symbol, period, 500)

	req := request.Request{
		Method: "GET",
		URL:    u,
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

	result, err := parseSinaKlineResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

func parseSinaKlineResponse(body []byte) (*KlineResult, error) {
	var resp sinaKlineResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if len(resp.Day) == 0 {
		return &KlineResult{Data: nil, Count: 0}, nil
	}

	bars := make([]KlineBar, 0, len(resp.Day))
	for _, d := range resp.Day {
		bar := KlineBar{
			Date:   d.D,
			Open:   parseFloat(d.O),
			High:   parseFloat(d.H),
			Low:    parseFloat(d.L),
			Close:  parseFloat(d.C),
			Volume: parseFloat(d.V),
		}
		bars = append(bars, bar)
	}

	return &KlineResult{Data: bars, Count: len(bars)}, nil
}

func toSinaSymbol(symbol string) (string, error) {
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
	case "5m":
		return "5"
	case "15m":
		return "15"
	case "30m":
		return "30"
	case "60m":
		return "60"
	case "1d", "":
		return "d"
	case "1w":
		return "w"
	case "1M":
		return "m"
	default:
		return "d"
	}
}

func ParseDate(dateStr string) time.Time {
	if strings.Contains(dateStr, " ") {
		t, _ := time.ParseInLocation("2006-01-02 15:04", dateStr, timeLoc)
		return t
	}
	t, _ := time.ParseInLocation("2006-01-02", dateStr, timeLoc)
	return t
}

var timeLoc, _ = time.LoadLocation("Asia/Shanghai")
