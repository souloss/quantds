// Package instrument provides securities/instruments list domain types.
//
// This package defines the request/response types for securities list retrieval
// across multiple markets and asset types (stocks, funds, bonds, etc.).
// It consolidates the former stocklist domain type.
package instrument

import (
	"context"

	"github.com/souloss/quantds/domain"
)

// AssetType represents the type of financial instrument.
type AssetType string

const (
	AssetTypeStock   AssetType = "STOCK"   // 股票
	AssetTypeFund    AssetType = "FUND"    // 基金
	AssetTypeBond    AssetType = "BOND"    // 债券
	AssetTypeIndex   AssetType = "INDEX"   // 指数
	AssetTypeETF     AssetType = "ETF"     // ETF
	AssetTypeOption  AssetType = "OPTION"  // 期权
	AssetTypeFutures AssetType = "FUTURES" // 期货
)

// Exchange represents the trading exchange.
type Exchange = domain.Exchange

const (
	ExchangeSH = domain.ExchangeSH // 上海证券交易所
	ExchangeSZ = domain.ExchangeSZ // 深圳证券交易所
	ExchangeBJ = domain.ExchangeBJ // 北京证券交易所
)

// Status represents the trading status of an instrument.
type Status string

const (
	StatusNormal    Status = "NORMAL"    // 正常交易
	StatusSuspended Status = "SUSPENDED" // 停牌
	StatusST        Status = "ST"        // ST股票
	StatusStarST    Status = "*ST"       // *ST股票
	StatusDelisted  Status = "DELISTED"  // 退市
)

// Request represents a securities list request.
type Request struct {
	Exchange   Exchange  // 按交易所筛选
	AssetType  AssetType // 按资产类型筛选
	Market     string    // 市场板块
	PageSize   int       // 分页大小
	PageNumber int       // 页码
}

// CacheKey returns the cache key for the request.
func (r Request) CacheKey() string {
	key := "instrument:"
	if r.Exchange != "" {
		key += string(r.Exchange)
	}
	if r.AssetType != "" {
		key += ":" + string(r.AssetType)
	}
	return key
}

// Response represents a securities list response.
type Response struct {
	Data       []Instrument // 证券列表
	Total      int          // 总数
	Source     string       // 数据源名称
	PageNumber int          // 当前页码
	PageSize   int          // 分页大小
}

// Instrument represents a single financial instrument/security.
type Instrument struct {
	Symbol     string    // 标准代码 (e.g., "000001.SZ")
	Code       string    // 原始代码 (e.g., "000001")
	Name       string    // 证券名称
	FullName   string    // 证券全称
	Exchange   Exchange  // 交易所
	Market     string    // 市场板块
	Industry   string    // 行业分类
	Sector     string    // 板块分类
	ListDate   string    // 上市日期 (YYYY-MM-DD)
	DelistDate string    // 退市日期
	Status     Status    // 交易状态
	AssetType  AssetType // 资产类型
	Currency   string    // 币种 (CNY, USD, etc.)
}

// Source defines the interface for instrument list providers.
type Source interface {
	Name() string
	Fetch(ctx context.Context, req Request) (Response, error)
	HealthCheck(ctx context.Context) error
}

// FormatSymbol formats code and exchange into a symbol string.
func FormatSymbol(code string, exchange Exchange) string {
	return domain.FormatSymbol(code, exchange)
}

// GuessStatus guesses the trading status from the stock name.
func GuessStatus(name string) Status {
	if len(name) >= 3 && name[:3] == "*ST" {
		return StatusStarST
	}
	if len(name) >= 2 && name[:2] == "ST" {
		return StatusST
	}
	return StatusNormal
}
