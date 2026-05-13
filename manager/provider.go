package manager

import (
	"context"

	"github.com/souloss/quantds/domain"
)

type Provider[Req, Resp any] interface {
	Name() string
	// SupportedMarkets 返回该 Provider 支持的市场列表
	SupportedMarkets() []domain.Market
	// CanHandle 检查是否支持处理指定的 symbol
	CanHandle(symbol string) bool
	// Fetch 从数据源获取数据
	Fetch(ctx context.Context, req Req) (Resp, *RequestTrace, error)
}

// HealthCheckable is an optional interface that Providers can implement
// to support health checking. When a Provider implements this interface,
// the Manager can verify its availability via HealthCheck.
type HealthCheckable interface {
	HealthCheck(ctx context.Context) error
}

type ProviderInfo struct {
	Name     string
	Priority int
	Weight   int
	Tags     map[string]string
}

type ProviderOption func(*ProviderInfo)

func WithPriority(priority int) ProviderOption {
	return func(info *ProviderInfo) {
		info.Priority = priority
	}
}

func WithWeight(weight int) ProviderOption {
	return func(info *ProviderInfo) {
		info.Weight = weight
	}
}

func WithTags(tags map[string]string) ProviderOption {
	return func(info *ProviderInfo) {
		info.Tags = tags
	}
}
