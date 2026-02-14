package tencent

import (
        "context"
        "fmt"
        "strings"

        "github.com/souloss/quantds/request"
        "golang.org/x/text/encoding/simplifiedchinese"
)

// MoneyFlowParams represents parameters for money flow request
type MoneyFlowParams struct {
        Symbol string
}

// MoneyFlowData represents real-time money flow data
type MoneyFlowData struct {
        Code           string  // Stock code
        Name           string  // Stock name
        MainIn         float64 // Main inflow
        MainOut        float64 // Main outflow
        MainNet        float64 // Main net inflow
        MainRatio      float64 // Main net inflow ratio
        SuperIn        float64
        SuperOut       float64
        SuperNet       float64
        SuperRatio     float64
        LargeIn        float64
        LargeOut       float64
        LargeNet       float64
        LargeRatio     float64
        MediumIn       float64
        MediumOut      float64
        MediumNet      float64
        MediumRatio    float64
        SmallIn        float64
        SmallOut       float64
        SmallNet       float64
        SmallRatio     float64
        Date           string
}

// GetMoneyFlow retrieves real-time money flow
func (c *Client) GetMoneyFlow(ctx context.Context, params *MoneyFlowParams) (*MoneyFlowData, *request.Record, error) {
        if params.Symbol == "" {
                return nil, nil, fmt.Errorf("symbol required")
        }

        sym, err := toTencentSymbol(params.Symbol)
        if err != nil {
                return nil, nil, err
        }

        url := MoneyFlowAPI + sym

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

        // GBK decoding
        utf8Body, err := simplifiedchinese.GBK.NewDecoder().Bytes(resp.Body)
        if err != nil {
                return nil, record, fmt.Errorf("decode gbk: %w", err)
        }

        data, err := parseMoneyFlow(string(utf8Body))
        if err != nil {
                return nil, record, err
        }
        return data, record, nil
}

func parseMoneyFlow(body string) (*MoneyFlowData, error) {
        // Format: v_ff_sz000001="code~name~..."
        parts := strings.Split(body, "\"")
        if len(parts) < 2 {
                return nil, fmt.Errorf("invalid response format")
        }

        data := parts[1]
        fields := strings.Split(data, "~")
        if len(fields) < 23 {
                return nil, fmt.Errorf("insufficient fields: %d", len(fields))
        }

        return &MoneyFlowData{
                Code:        fields[0],
                Name:        fields[1],
                MainIn:      parseFloat(fields[2]),
                MainOut:     parseFloat(fields[3]),
                MainNet:     parseFloat(fields[4]),
                MainRatio:   parseFloat(fields[5]),
                SuperIn:     parseFloat(fields[6]),
                SuperOut:    parseFloat(fields[7]),
                SuperNet:    parseFloat(fields[8]),
                SuperRatio:  parseFloat(fields[9]),
                LargeIn:     parseFloat(fields[10]),
                LargeOut:    parseFloat(fields[11]),
                LargeNet:    parseFloat(fields[12]),
                LargeRatio:  parseFloat(fields[13]),
                MediumIn:    parseFloat(fields[14]),
                MediumOut:   parseFloat(fields[15]),
                MediumNet:   parseFloat(fields[16]),
                MediumRatio: parseFloat(fields[17]),
                SmallIn:     parseFloat(fields[18]),
                SmallOut:    parseFloat(fields[19]),
                SmallNet:    parseFloat(fields[20]),
                SmallRatio:  parseFloat(fields[21]),
                Date:        fields[22],
        }, nil
}
