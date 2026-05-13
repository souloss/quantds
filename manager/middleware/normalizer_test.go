package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/manager"
)

var errTest = errors.New("test error")

func TestNormalizer_AppliesFunction(t *testing.T) {
	transform := func(s string) string {
		return s + "_normalized"
	}

	n := Normalizer[string, string](transform)

	next := &mockProvider[string, string]{
		name:  "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "data", nil, nil },
	}

	wrapped := n(next)

	resp, _, err := wrapped.Fetch(context.Background(), "req")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != "data_normalized" {
		t.Errorf("got %q, want %q", resp, "data_normalized")
	}
}

func TestNormalizer_PassesThroughErrors(t *testing.T) {
	errFail := errTest
	n := Normalizer[string, string](func(s string) string { return s + "_normalized" })

	next := &mockProvider[string, string]{
		name:  "test",
		fetch: func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "", nil, errFail },
	}

	wrapped := n(next)

	_, _, err := wrapped.Fetch(context.Background(), "req")
	if err != errFail {
		t.Errorf("got error %v, want %v", err, errFail)
	}
}

func TestNormalizer_PreservesProviderInfo(t *testing.T) {
	n := Normalizer[string, string](func(s string) string { return s })

	next := &mockProvider[string, string]{
		name:    "myprovider",
		fetch:   func(ctx context.Context, req string) (string, *manager.RequestTrace, error) { return "ok", nil, nil },
		markets: []domain.Market{domain.MarketCN},
	}

	wrapped := n(next)

	if wrapped.Name() != "myprovider" {
		t.Errorf("Name() = %q, want %q", wrapped.Name(), "myprovider")
	}
	if !wrapped.CanHandle("anything") {
		t.Error("CanHandle() = false, want true")
	}
}

func TestNormalizer_WithIntType(t *testing.T) {
	// Test with a different type to verify generic works
	double := func(n int) int { return n * 2 }

	n := Normalizer[string, int](double)

	next := &mockProvider[string, int]{
		name:  "test",
		fetch: func(ctx context.Context, req string) (int, *manager.RequestTrace, error) { return 21, nil, nil },
	}

	wrapped := n(next)

	resp, _, err := wrapped.Fetch(context.Background(), "req")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != 42 {
		t.Errorf("got %d, want %d", resp, 42)
	}
}
