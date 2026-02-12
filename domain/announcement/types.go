// Package announcement provides announcements and news domain types.
//
// This package defines the request/response types for company announcements,
// news, and regulatory filings across multiple markets.
// It consolidates the former news domain type.
package announcement

import (
	"context"
	"time"
)

// AnnouncementType represents the type of announcement.
type AnnouncementType string

const (
	TypeReport     AnnouncementType = "REPORT"      // 定期报告
	TypeEarnings   AnnouncementType = "EARNINGS"    // 业绩公告
	TypeDividend   AnnouncementType = "DIVIDEND"    // 分红派息
	TypeRights     AnnouncementType = "RIGHTS"      // 股权变动
	TypeMajorEvent AnnouncementType = "MAJOR_EVENT" // 重大事项
	TypeRegulatory AnnouncementType = "REGULATORY"  // 监管公告
	TypeNews       AnnouncementType = "NEWS"        // 新闻报道
	TypeResearch   AnnouncementType = "RESEARCH"    // 研究报告
)

// Category represents the announcement category.
type Category string

const (
	CategoryCompany    Category = "COMPANY"    // 公司公告
	CategoryExchange   Category = "EXCHANGE"   // 交易所公告
	CategoryRegulatory Category = "REGULATORY" // 监管公告
	CategoryMedia      Category = "MEDIA"      // 媒体资讯
)

// Request represents an announcement/news request.
type Request struct {
	Symbol    string             // 标的代码
	Types     []AnnouncementType // 按类型筛选
	Category  Category           // 按分类筛选
	PageSize  int                // 分页大小
	PageIndex int                // 页码（从1开始）
	StartTime *time.Time         // 起始时间
	EndTime   *time.Time         // 结束时间
}

// CacheKey returns the cache key for the request.
func (r Request) CacheKey() string {
	key := "announcement:"
	if r.Symbol != "" {
		key += r.Symbol
	}
	if r.Category != "" {
		key += ":" + string(r.Category)
	}
	return key
}

// Response represents an announcement/news response.
type Response struct {
	Symbol     string         // 标的代码
	Data       []Announcement // 公告列表
	Source     string         // 数据源名称
	HasMore    bool           // 是否有更多数据
	TotalCount int            // 总数
	PageIndex  int            // 当前页码
	PageSize   int            // 分页大小
}

// Announcement represents a single announcement/news item.
type Announcement struct {
	ID          string           // 唯一标识
	Title       string           // 标题
	Content     string           // 正文内容
	Summary     string           // 摘要
	PublishTime string           // 发布时间
	Source      string           // 来源
	URL         string           // 链接
	Type        AnnouncementType // 公告类型
	Category    Category         // 公告分类
	Code        string           // 股票代码
	Name        string           // 股票名称
	AttachFiles []Attachment     // 附件列表
}

// Attachment represents an attached file.
type Attachment struct {
	Name string // 文件名
	URL  string // 文件链接
	Size int64  // 文件大小（字节）
}

// Item is an alias for Announcement for backward compatibility.
type Item = Announcement

// Source defines the interface for announcement data providers.
type Source interface {
	Name() string
	Fetch(ctx context.Context, req Request) (Response, error)
	HealthCheck(ctx context.Context) error
}
