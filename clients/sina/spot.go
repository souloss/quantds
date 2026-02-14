package sina

import (
        "context"
        "fmt"
        "strings"
        "time"

        "github.com/souloss/quantds/request"
        "golang.org/x/text/encoding/simplifiedchinese"
)

// SpotParams 实时行情查询参数
type SpotParams struct {
        Symbols []string // 股票代码列表，格式: sh600001, sz000001 或 600001.SH, 000001.SZ
}

// SpotResult 实时行情查询结果
type SpotResult struct {
        Data []SpotQuote
}

// SpotQuote 单只股票实时行情
type SpotQuote struct {
        Symbol   string  // 新浪格式股票代码 (sh600001)
        Name     string  // 股票名称
        Open     float64 // 今开盘
        PreClose float64 // 昨收盘
        Latest   float64 // 最新价
        High     float64 // 最高价
        Low      float64 // 最低价
        Volume   float64 // 成交量（手）
        Amount   float64 // 成交额（万元）
        Date     string  // 日期
        Time     string  // 时间
}

// GetSpot 获取实时行情数据
//
// 参数说明:
//   - Symbols: 股票代码列表，支持两种格式:
//     1. 新浪格式: "sh600001", "sz000001"
//     2. 标准格式: "600001.SH", "000001.SZ"
//
// 限制:
//   - 单次请求建议不超过100只股票
//   - 返回数据使用 GBK 编码
//   - 数据有3-5秒延迟
//   - 盘前盘后数据不准确
//   - 高频请求会被限流
func (c *Client) GetSpot(ctx context.Context, params *SpotParams) (*SpotResult, *request.Record, error) {
        if len(params.Symbols) == 0 {
                return nil, nil, fmt.Errorf("symbols required")
        }

        symbols := make([]string, 0, len(params.Symbols))
        for _, s := range params.Symbols {
                sym, err := toSinaSymbol(s)
                if err != nil {
                        continue
                }
                symbols = append(symbols, sym)
        }

        if len(symbols) == 0 {
                return nil, nil, fmt.Errorf("no valid symbols")
        }

        url := fmt.Sprintf("%s/list=_%s&list=%s",
                SpotAPI, fmt.Sprintf("%d", time.Now().UnixMilli()), strings.Join(symbols, ","))

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

        // 新浪返回 GBK 编码，需要转换
        utf8Body, err := simplifiedchinese.GBK.NewDecoder().Bytes(resp.Body)
        if err != nil {
                return nil, record, fmt.Errorf("decode gbk: %w", err)
        }

        result, err := parseSinaSpot(string(utf8Body))
        if err != nil {
                return nil, record, err
        }

        return result, record, nil
}

func parseSinaSpot(body string) (*SpotResult, error) {
        // 格式: var hq_str_sh600001="浦发银行,8.50,8.45,8.60,8.65,8.45,8.60,8.61,12345678,123456789,..."
        lines := strings.Split(body, ";")
        quotes := make([]SpotQuote, 0, len(lines))

        for _, line := range lines {
                line = strings.TrimSpace(line)
                if line == "" {
                        continue
                }

                // 提取 var hq_str_xxx = "..." 中的内容
                startIdx := strings.Index(line, "=\"")
                endIdx := strings.LastIndex(line, "\"")
                if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
                        continue
                }

                // 提取股票代码
                symbolStart := strings.Index(line, "hq_str_")
                if symbolStart == -1 {
                        continue
                }
                symbol := line[symbolStart+7 : startIdx]

                // 解析数据
                data := line[startIdx+2 : endIdx]
                fields := strings.Split(data, ",")
                if len(fields) < 32 {
                        continue
                }

                quote := SpotQuote{
                        Symbol:   symbol,
                        Name:     fields[0],
                        Open:     parseFloat(fields[1]),
                        PreClose: parseFloat(fields[2]),
                        Latest:   parseFloat(fields[3]),
                        High:     parseFloat(fields[4]),
                        Low:      parseFloat(fields[5]),
                        Volume:   parseFloat(fields[8]),
                        Amount:   parseFloat(fields[9]),
                        Date:     fields[30],
                        Time:     fields[31],
                }
                quotes = append(quotes, quote)
        }

        return &SpotResult{Data: quotes}, nil
}
