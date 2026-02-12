package tencent

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/souloss/quantds/request"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	SpotAPI       = "http://qt.gtimg.cn/q="
	SpotPrefixVar = "v_"
)

type SpotParams struct {
	Symbols []string
}

type SpotResult struct {
	Data  []SpotQuote
	Count int
}

type SpotQuote struct {
	Symbol     string
	Name       string
	Latest     float64
	PreClose   float64
	Open       float64
	High       float64
	Low        float64
	Volume     float64
	Change     float64
	ChangeRate float64
	Time       string
}

func (c *Client) GetSpot(ctx context.Context, params *SpotParams) (*SpotResult, *request.Record, error) {
	if params == nil || len(params.Symbols) == 0 {
		return nil, nil, fmt.Errorf("symbols required")
	}

	symbols := make([]string, 0, len(params.Symbols))
	for _, s := range params.Symbols {
		sym, err := toTencentSymbol(s)
		if err != nil {
			continue
		}
		symbols = append(symbols, sym)
	}

	if len(symbols) == 0 {
		return nil, nil, fmt.Errorf("no valid symbols")
	}

	url := SpotAPI + strings.Join(symbols, ",")

	req := request.Request{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":    "https://gu.qq.com/",
		},
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	utf8Body, err := gbkToUtf8(resp.Body)
	if err != nil {
		return nil, record, fmt.Errorf("decode gbk: %w", err)
	}

	result, err := parseSpotRows(string(utf8Body))
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

func parseSpotRows(body string) (*SpotResult, error) {
	var quotes []SpotQuote

	for _, line := range strings.Split(body, ";") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		left := strings.TrimSpace(parts[0])
		if !strings.HasPrefix(left, SpotPrefixVar) {
			continue
		}

		symbol := strings.TrimPrefix(left, SpotPrefixVar)
		right := strings.Trim(parts[1], "\"")
		if right == "" {
			continue
		}

		fields := strings.Split(right, "~")
		if len(fields) < 45 {
			continue
		}

		quotes = append(quotes, SpotQuote{
			Symbol:     symbol,
			Name:       fields[1],
			Latest:     parseFloat(fields[3]),
			PreClose:   parseFloat(fields[4]),
			Open:       parseFloat(fields[5]),
			Volume:     parseFloat(fields[6]),
			High:       parseFloat(fields[33]),
			Low:        parseFloat(fields[34]),
			Change:     parseFloat(fields[31]),
			ChangeRate: parseFloat(fields[32]),
			Time:       fields[30],
		})
	}

	return &SpotResult{Data: quotes, Count: len(quotes)}, nil
}

func gbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}
