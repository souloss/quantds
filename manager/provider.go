package manager

import (
	"context"

	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/request"
)

type Provider[Req, Resp any] interface {
	Name() string
	Fetch(ctx context.Context, client request.Client, req Req) (Resp, *RequestTrace, error)

	// SupportedMarkets 返回该 Provider 支持的市场列表
	SupportedMarkets() []domain.Market

	// CanHandle 检查是否支持处理指定的 symbol
	CanHandle(symbol string) bool
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
