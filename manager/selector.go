package manager

import "sort"

type Selector interface {
	Select(providers []ProviderInfo) []string
}

type PrioritySelector struct{}

func NewPrioritySelector() *PrioritySelector {
	return &PrioritySelector{}
}

func (s *PrioritySelector) Select(providers []ProviderInfo) []string {
	if len(providers) == 0 {
		return nil
	}

	sorted := make([]ProviderInfo, len(providers))
	copy(sorted, providers)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority > sorted[j].Priority
	})

	names := make([]string, len(sorted))
	for i, p := range sorted {
		names[i] = p.Name
	}
	return names
}

type WeightedSelector struct{}

func NewWeightedSelector() *WeightedSelector {
	return &WeightedSelector{}
}

func (s *WeightedSelector) Select(providers []ProviderInfo) []string {
	if len(providers) == 0 {
		return nil
	}

	sorted := make([]ProviderInfo, len(providers))
	copy(sorted, providers)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Weight > sorted[j].Weight
	})

	names := make([]string, len(sorted))
	for i, p := range sorted {
		names[i] = p.Name
	}
	return names
}

var _ Selector = (*PrioritySelector)(nil)
var _ Selector = (*WeightedSelector)(nil)
