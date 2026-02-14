package sina

import (
        "context"
        "encoding/json"
        "fmt"
        "net/url"

        "github.com/souloss/quantds/request"
)

// MoneyFlowParams represents parameters for money flow request
type MoneyFlowParams struct {
        Symbol string
        Count  int
}

// MoneyFlowItem represents a single day's money flow
type MoneyFlowItem struct {
        Date           string  `json:"opendate"`
        Trade          float64 `json:"trade"`       // Close price
        ChangeRatio    float64 `json:"changeratio"` // Change percent
        Turnover       float64 `json:"turnover"`    // Turnover rate
        NetAmount      float64 `json:"netamount"`   // Net inflow amount (Main)
        RatioAmount    float64 `json:"ratioamount"` // Net inflow ratio
        R0Net          float64 `json:"r0_net"`      // Super Large
        R0Ratio        float64 `json:"r0_ratio"`
        R0xRatio       float64 `json:"r0x_ratio"`
        CntR0xRatio    float64 `json:"cnt_r0x_ratio"`
        CateR0xRatio   float64 `json:"cate_ra"`
        R1Net          float64 `json:"r1_net"` // Large
        R1Ratio        float64 `json:"r1_ratio"`
        R2Net          float64 `json:"r2_net"` // Medium
        R2Ratio        float64 `json:"r2_ratio"`
        R3Net          float64 `json:"r3_net"` // Small
        R3Ratio        float64 `json:"r3_ratio"`
}

// GetMoneyFlow retrieves money flow history
func (c *Client) GetMoneyFlow(ctx context.Context, params *MoneyFlowParams) ([]MoneyFlowItem, *request.Record, error) {
        if params.Symbol == "" {
                return nil, nil, fmt.Errorf("symbol required")
        }

        sym, err := toSinaSymbol(params.Symbol)
        if err != nil {
                return nil, nil, err
        }

        count := params.Count
        if count <= 0 {
                count = 20
        }

        query := url.Values{}
        query.Set("num", fmt.Sprintf("%d", count))
        query.Set("sort", "opendate")
        query.Set("asc", "0")
        query.Set("symbol", sym)

        url := fmt.Sprintf("%s?%s", MoneyFlowAPI, query.Encode())

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

        // Sina might return keys without quotes? 
        // Usually this API returns valid JSON: [{"opendate":"2023-10-10",...}]
        
        var items []MoneyFlowItem
        // Handle potential "null" or empty response
        if len(resp.Body) == 0 || string(resp.Body) == "null" {
                return nil, record, nil
        }

        // Sina keys are sometimes not quoted in other APIs, but json_v2.php usually returns valid JSON.
        // If unmarshal fails, we might need manual parsing or lenient parser.
        // But let's try standard unmarshal first.
        // Note: The float fields might be strings in JSON, need to handle that?
        // The struct defines float64, if JSON has strings "123.45", Go Unmarshal might fail or need custom UnmarshalJSON.
        // Sina often returns strings for numbers.
        
        // Let's use a temporary struct with interface{} or strings to be safe, or just try.
        // Or define custom unmarshaller.
        // To be safe, I'll use a helper to unmarshal via a map.
        
        var rawItems []map[string]interface{}
        if err := json.Unmarshal(resp.Body, &rawItems); err != nil {
                // If standard unmarshal fails, it might be due to keys without quotes (though v2 usually has them).
                // Or values are not standard.
                return nil, record, fmt.Errorf("unmarshal error: %w body: %s", err, string(resp.Body))
        }

        items = make([]MoneyFlowItem, 0, len(rawItems))
        for _, raw := range rawItems {
                items = append(items, MoneyFlowItem{
                        Date:        getString(raw, "opendate"),
                        Trade:       getFloat(raw, "trade"),
                        ChangeRatio: getFloat(raw, "changeratio"),
                        Turnover:    getFloat(raw, "turnover"),
                        NetAmount:   getFloat(raw, "netamount"),
                        RatioAmount: getFloat(raw, "ratioamount"),
                        R0Net:       getFloat(raw, "r0_net"),
                        R0Ratio:     getFloat(raw, "r0_ratio"),
                        R1Net:       getFloat(raw, "r1_net"),
                        R1Ratio:     getFloat(raw, "r1_ratio"),
                        R2Net:       getFloat(raw, "r2_net"),
                        R2Ratio:     getFloat(raw, "r2_ratio"),
                        R3Net:       getFloat(raw, "r3_net"),
                        R3Ratio:     getFloat(raw, "r3_ratio"),
                })
        }

        return items, record, nil
}

// Helpers
func getString(m map[string]interface{}, key string) string {
        if v, ok := m[key]; ok {
                if s, ok := v.(string); ok {
                        return s
                }
                return fmt.Sprintf("%v", v)
        }
        return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
        if v, ok := m[key]; ok {
                switch val := v.(type) {
                case float64:
                        return val
                case string:
                        var f float64
                        fmt.Sscanf(val, "%f", &f)
                        return f
                }
        }
        return 0
}
