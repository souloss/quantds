package tencent

import (
        "context"
        "fmt"
        "strings"

        "github.com/souloss/quantds/request"
)

// QuotePrefixVar is the prefix for quote response variables
const QuotePrefixVar = "v_"

// QuoteParams represents parameters for real-time quote request
type QuoteParams struct {
        Symbols []string // List of stock symbols
}

// QuoteResult represents the real-time quote result
type QuoteResult struct {
        Data  []QuoteData
        Count int
}

// QuoteData represents a single real-time market quote
type QuoteData struct {
        Symbol     string  // Stock symbol
        Name       string  // Stock name
        Latest     float64 // Latest price
        PreClose   float64 // Previous closing price
        Open       float64 // Opening price
        High       float64 // Highest price
        Low        float64 // Lowest price
        Volume     float64 // Trading volume
        Change     float64 // Price change
        ChangeRate float64 // Change rate (%)
        Time       string  // Quote time
}

// GetQuotes retrieves real-time quotes for multiple stocks
func (c *Client) GetQuotes(ctx context.Context, params *QuoteParams) (*QuoteResult, *request.Record, error) {
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

        url := QuoteAPI + strings.Join(symbols, ",")

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

        utf8Body, err := gbkToUtf8(resp.Body)
        if err != nil {
                return nil, record, fmt.Errorf("decode gbk: %w", err)
        }

        result, err := parseQuoteRows(string(utf8Body))
        if err != nil {
                return nil, record, err
        }

        return result, record, nil
}

func parseQuoteRows(body string) (*QuoteResult, error) {
        var quotes []QuoteData

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
                if !strings.HasPrefix(left, QuotePrefixVar) {
                        continue
                }

                symbol := strings.TrimPrefix(left, QuotePrefixVar)
                right := strings.Trim(parts[1], "\"")
                if right == "" {
                        continue
                }

                fields := strings.Split(right, "~")
                if len(fields) < 45 {
                        continue
                }

                quotes = append(quotes, QuoteData{
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

        return &QuoteResult{Data: quotes, Count: len(quotes)}, nil
}
