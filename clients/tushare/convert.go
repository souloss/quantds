package tushare

import (
	"fmt"
	"strings"
)

// ParseTushareSymbol 解析 symbol 字符串为 code 和 exchange 部分。
// e.g., "000001.SZ" → ("000001", "SZ", true)
func ParseTushareSymbol(symbol string) (code, exchange string, ok bool) {
	if len(symbol) < 6 {
		return "", "", false
	}
	idx := strings.IndexByte(symbol, '.')
	if idx > 0 && idx < len(symbol)-1 {
		return symbol[:idx], symbol[idx+1:], true
	}
	return "", "", false
}

// ToTushareSymbol 将标准 symbol 格式转换为 Tushare ts_code 格式。
// Tushare 使用 "CODE.EXCHANGE" 格式，与本系统的 "CODE.EXCHANGE" 一致。
// e.g., "000001.SZ" → "000001.SZ"
func ToTushareSymbol(symbol string) (string, error) {
	_, _, ok := ParseTushareSymbol(symbol)
	if !ok {
		if strings.Contains(symbol, ".") {
			return symbol, nil
		}
		return "", fmt.Errorf("invalid symbol: %s", symbol)
	}
	return symbol, nil
}

// ToTushareExchange 将系统交易所代码转换为 Tushare 交易所代码。
// SH → SSE, SZ → SZSE, BJ → BSE
func ToTushareExchange(exchange string) string {
	switch exchange {
	case "SH":
		return "SSE"
	case "SZ":
		return "SZSE"
	case "BJ":
		return "BSE"
	default:
		return exchange
	}
}
