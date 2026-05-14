// Package kline provides K-line (candlestick) data domain types.
//
// This package defines the request/response types for K-line data retrieval
// across multiple markets and asset types.
package kline

import (
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
	Timeframe60m Timeframe = "60m" // 60分钟 (兼容旧值)
	Timeframe1H  Timeframe = "1H"  // 1小时
	Timeframe4H  Timeframe = "4H"  // 4小时
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
	Symbol     string     // 标的代码 (e.g., "000001.SZ", "600519.SH")
	Timeframe  Timeframe  // K线周期
	StartTime  time.Time  // 起始时间
	EndTime    time.Time  // 结束时间
	Limit      int        // 返回K线条数上限 (0表示不限)
	BeforeTime time.Time  // 分页游标：返回该时间之前的K线
	Adjust     AdjustType // 复权类型
}

// CacheKey returns the cache key for the request.
func (r Request) CacheKey() string {
	return "kline:" + r.Symbol + ":" + string(r.Timeframe) + ":" +
		r.StartTime.Format("20060102") + ":" + r.EndTime.Format("20060102")
}

// Response represents a K-line data response.
type Response struct {
	Symbol      string // 标的代码
	Bars        []Bar  // K线数据
	Source      string // 数据源名称
	DataVersion int    // 数据格式版本（当前为1）
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

// ParseSymbol parses a symbol string into code and exchange.
func ParseSymbol(symbol string) (code string, exchange Exchange, ok bool) {
	return domain.ParseSymbol(symbol)
}

// FormatSymbol formats code and exchange into a symbol string.
func FormatSymbol(code string, exchange Exchange) string {
	return domain.FormatSymbol(code, exchange)
}

// Normalize fills in computed fields for each Bar that are zero-valued.
// It computes Change, ChangeRate from Open/Close, and Turnover from Close*Volume
// when those fields are missing (zero).
// When Open is zero, the data is considered incomplete and no derived fields are computed.
func (b *Bar) Normalize() {
	if b.Open == 0 {
		return
	}
	if b.Change == 0 {
		b.Change = b.Close - b.Open
	}
	if b.ChangeRate == 0 {
		b.ChangeRate = (b.Close - b.Open) / b.Open * 100
	}
	if b.Turnover == 0 && b.Close != 0 && b.Volume != 0 {
		b.Turnover = b.Close * b.Volume
	}
}

// Normalize fills in computed fields for all Bars in the response.
func (r *Response) Normalize() {
	for i := range r.Bars {
		r.Bars[i].Normalize()
	}
	if r.DataVersion == 0 {
		r.DataVersion = 1
	}
}
