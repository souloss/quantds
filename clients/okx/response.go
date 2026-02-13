package okx

import (
	"encoding/json"
)

// Response represents the standard OKX API response wrapper
type Response struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}
