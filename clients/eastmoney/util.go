package eastmoney

import (
	"fmt"
	"strings"
)

// Helper functions for parsing responses

func getString(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case string:
			return val
		case float64:
			return fmt.Sprintf("%.0f", val)
		}
	}
	return ""
}

func getFloat(data map[string]interface{}, key string) float64 {
	if v, ok := data[key]; ok {
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

func getInt(data map[string]interface{}, key string) int {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case string:
			var i int
			fmt.Sscanf(val, "%d", &i)
			return i
		}
	}
	return 0
}

// toEastMoneySecid converts symbol to EastMoney secid format
func toEastMoneySecid(symbol string) (string, error) {
	code, exchange, ok := parseSymbol(symbol)
	if !ok {
		return "", fmt.Errorf("invalid symbol: %s", symbol)
	}
	switch exchange {
	case "SH":
		return "1." + code, nil
	case "SZ", "BJ":
		return "0." + code, nil
	default:
		return "1." + code, nil
	}
}

// parseSymbol parses a symbol string into code and exchange
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
