// Package spot provides real-time market quote domain types.
//
// This package defines the request/response types for real-time spot/quote
// data retrieval across multiple markets.
package spot

import (
	"context"
	"time"

	"github.com/souloss/quantds/domain"
)

// Exchange represents the trading exchange.
type Exchange = domain.Exchange

const (
	ExchangeSH = domain.ExchangeSH // 上海证券交易所
	ExchangeSZ = domain.ExchangeSZ // 深圳证券交易所
	ExchangeBJ = domain.ExchangeBJ // 北京证券交易所
)

// Request represents a real-time quote request.
type Request struct {
	Symbols []string // 标的代码列表
}

// CacheKey returns the cache key for the request.
func (r Request) CacheKey() string {
	if len(r.Symbols) == 0 {
		return "spot:all"
	}
	return "spot:" + r.Symbols[0]
}

// Response represents a real-time quote response.
type Response struct {
	Quotes []Quote // 行情列表
	Total  int     // 总数
	Source string  // 数据源名称
}

// Quote represents a single real-time market quote.
type Quote struct {
	Symbol       string    // 标的代码
	Name         string    // 证券名称
	Latest       float64   // 最新价
	Open         float64   // 开盘价
	High         float64   // 最高价
	Low          float64   // 最低价
	PreClose     float64   // 昨收价
	Change       float64   // 涨跌额
	ChangeRate   float64   // 涨跌幅 (%)
	Volume       float64   // 成交量
	Turnover     float64   // 成交额
	Amplitude    float64   // 振幅 (%)
	TurnoverRate float64   // 换手率 (%)
	Timestamp    time.Time // 行情时间戳
	BidPrice     float64   // 买一价
	BidVolume    float64   // 买一量
	AskPrice     float64   // 卖一价
	AskVolume    float64   // 卖一量
}

// Source defines the interface for spot data providers.
type Source interface {
	Name() string
	Fetch(ctx context.Context, req Request) (Response, error)
	HealthCheck(ctx context.Context) error
}

// ParseSymbol parses a symbol string into code and exchange.
func ParseSymbol(symbol string) (code string, exchange Exchange, ok bool) {
	return domain.ParseSymbol(symbol)
}

// FormatSymbol formats code and exchange into a symbol string.
func FormatSymbol(code string, exchange Exchange) string {
	return domain.FormatSymbol(code, exchange)
}
