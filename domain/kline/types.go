// Package kline provides K-line (candlestick) data domain types.
//
// This package defines the request/response types for K-line data retrieval
// across multiple markets and asset types.
package kline

import (
	"context"
	"time"

	"github.com/souloss/quantds/domain"
)

// Timeframe represents the time interval for K-line data.
type Timeframe string

const (
	Timeframe1m  Timeframe = "1m"  // 1分钟
	Timeframe5m  Timeframe = "5m"  // 5分钟
	Timeframe15m Timeframe = "15m" // 15分钟
	Timeframe30m Timeframe = "30m" // 30分钟
	Timeframe60m Timeframe = "60m" // 60分钟
	Timeframe1d  Timeframe = "1d"  // 日线
	Timeframe1w  Timeframe = "1w"  // 周线
	Timeframe1M  Timeframe = "1M"  // 月线
)

// AdjustType represents the price adjustment type.
type AdjustType string

const (
	AdjustNone    AdjustType = ""    // 不复权
	AdjustForward AdjustType = "qfq" // 前复权
	AdjustBack    AdjustType = "hfq" // 后复权
)

// Exchange represents the trading exchange.
type Exchange = domain.Exchange

const (
	ExchangeSH = domain.ExchangeSH // 上海证券交易所
	ExchangeSZ = domain.ExchangeSZ // 深圳证券交易所
	ExchangeBJ = domain.ExchangeBJ // 北京证券交易所
)

// Request represents a K-line data request.
type Request struct {
	Symbol    string     // 标的代码 (e.g., "000001.SZ", "600519.SH")
	Timeframe Timeframe  // K线周期
	StartTime time.Time  // 起始时间
	EndTime   time.Time  // 结束时间
	Adjust    AdjustType // 复权类型
}

// CacheKey returns the cache key for the request.
func (r Request) CacheKey() string {
	return "kline:" + r.Symbol + ":" + string(r.Timeframe) + ":" +
		r.StartTime.Format("20060102") + ":" + r.EndTime.Format("20060102")
}

// Response represents a K-line data response.
type Response struct {
	Symbol string // 标的代码
	Bars   []Bar  // K线数据
	Source string // 数据源名称
}

// Bar represents a single K-line (OHLCV) data point.
type Bar struct {
	Timestamp    time.Time // 时间戳
	Open         float64   // 开盘价
	High         float64   // 最高价
	Low          float64   // 最低价
	Close        float64   // 收盘价
	Volume       float64   // 成交量
	Turnover     float64   // 成交额
	Change       float64   // 涨跌额
	ChangeRate   float64   // 涨跌幅 (%)
	TurnoverRate float64   // 换手率 (%)
}

// Source defines the interface for K-line data providers.
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
