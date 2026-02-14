package yahoo

import (
        "context"
        "encoding/json"
        "fmt"
        "strings"
        "time"

        "github.com/souloss/quantds/request"
)

// QuoteParams represents parameters for real-time quote request
type QuoteParams struct {
        Symbols []string // List of stock symbols (e.g., ["AAPL", "MSFT", "GOOGL"])
}

// QuoteResult represents the real-time quote result
type QuoteResult struct {
        Quotes []QuoteData // List of quotes
        Count  int         // Number of quotes
}

// QuoteData represents a single real-time quote
type QuoteData struct {
        Symbol           string    // Stock symbol
        Name             string    // Company name
        Exchange         string    // Exchange name
        Market           string    // Market (US)
        Latest           float64   // Latest price
        Open             float64   // Opening price
        High             float64   // Highest price (daily)
        Low              float64   // Lowest price (daily)
        PreClose         float64   // Previous close
        Change           float64   // Price change
        ChangeRate       float64   // Change rate (%)
        Volume           float64   // Trading volume
        AvgVolume        float64   // Average volume
        MarketCap        float64   // Market capitalization
        PE               float64   // P/E ratio
        High52Week       float64   // 52-week high
        Low52Week        float64   // 52-week low
        BidPrice         float64   // Best bid price
        BidSize          float64   // Best bid size
        AskPrice         float64   // Best ask price
        AskSize          float64   // Best ask size
        Timestamp        time.Time // Quote timestamp
        RegularMarketEnd time.Time // Market close time
}

// GetQuote retrieves real-time quotes for multiple symbols
func (c *Client) GetQuote(ctx context.Context, params *QuoteParams) (*QuoteResult, *request.Record, error) {
        if len(params.Symbols) == 0 {
                return &QuoteResult{}, nil, nil
        }

        symbols := strings.Join(params.Symbols, ",")
        url := fmt.Sprintf("%s%s?symbols=%s", BaseURL, QuoteAPI, symbols)

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

        result, err := parseQuoteResponse(resp.Body)
        if err != nil {
                return nil, record, err
        }

        return result, record, nil
}

// quoteResponse represents the Yahoo Finance quote API response
type quoteResponse struct {
        QuoteResponse struct {
                Result []struct {
                        Symbol                     string  `json:"symbol"`
                        ShortName                  string  `json:"shortName"`
                        LongName                   string  `json:"longName"`
                        Exchange                   string  `json:"fullExchangeName"`
                        Market                     string  `json:"market"`
                        RegularMarketPrice         float64 `json:"regularMarketPrice"`
                        RegularMarketOpen          float64 `json:"regularMarketOpen"`
                        DayHigh                    float64 `json:"regularMarketDayHigh"`
                        DayLow                     float64 `json:"regularMarketDayLow"`
                        PreviousClose              float64 `json:"regularMarketPreviousClose"`
                        RegularMarketChange        float64 `json:"regularMarketChange"`
                        RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
                        RegularMarketVolume        float64 `json:"regularMarketVolume"`
                        AverageDailyVolume3Month   float64 `json:"averageDailyVolume3Month"`
                        MarketCap                  float64 `json:"marketCap"`
                        TrailingPE                 float64 `json:"trailingPE"`
                        FiftyTwoWeekHigh           float64 `json:"fiftyTwoWeekHigh"`
                        FiftyTwoWeekLow            float64 `json:"fiftyTwoWeekLow"`
                        Bid                        float64 `json:"bid"`
                        BidSize                    float64 `json:"bidSize"`
                        Ask                        float64 `json:"ask"`
                        AskSize                    float64 `json:"askSize"`
                        RegularMarketTime          int64   `json:"regularMarketTime"`
                } `json:"result"`
                Error interface{} `json:"error"`
        } `json:"quoteResponse"`
}

func parseQuoteResponse(body []byte) (*QuoteResult, error) {
        var resp quoteResponse
        if err := json.Unmarshal(body, &resp); err != nil {
                return nil, err
        }

        quotes := make([]QuoteData, 0, len(resp.QuoteResponse.Result))
        for _, q := range resp.QuoteResponse.Result {
                name := q.ShortName
                if name == "" {
                        name = q.LongName
                }

                quotes = append(quotes, QuoteData{
                        Symbol:     q.Symbol,
                        Name:       name,
                        Exchange:   q.Exchange,
                        Market:     q.Market,
                        Latest:     q.RegularMarketPrice,
                        Open:       q.RegularMarketOpen,
                        High:       q.DayHigh,
                        Low:        q.DayLow,
                        PreClose:   q.PreviousClose,
                        Change:     q.RegularMarketChange,
                        ChangeRate: q.RegularMarketChangePercent,
                        Volume:     q.RegularMarketVolume,
                        AvgVolume:  q.AverageDailyVolume3Month,
                        MarketCap:  q.MarketCap,
                        PE:         q.TrailingPE,
                        High52Week: q.FiftyTwoWeekHigh,
                        Low52Week:  q.FiftyTwoWeekLow,
                        BidPrice:   q.Bid,
                        BidSize:    q.BidSize,
                        AskPrice:   q.Ask,
                        AskSize:    q.AskSize,
                        Timestamp:  time.Unix(q.RegularMarketTime, 0),
                })
        }

        return &QuoteResult{
                Quotes: quotes,
                Count:  len(quotes),
        }, nil
}

// SpotParams is an alias for QuoteParams
type SpotParams = QuoteParams

// SpotResult is an alias for QuoteResult
type SpotResult = QuoteResult

// GetSpot is an alias for GetQuote
func (c *Client) GetSpot(ctx context.Context, params *SpotParams) (*SpotResult, *request.Record, error) {
        return c.GetQuote(ctx, params)
}
