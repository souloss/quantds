package manager

import (
	"math"
	"testing"
)

func TestConsistencyChecker_CheckFloats(t *testing.T) {
	cc := NewConsistencyChecker(1.0) // 1% tolerance

	tests := []struct {
		name    string
		field   string
		values  map[string]float64
		wantDev float64
		wantOk  bool
	}{
		{
			name:    "consistent within tolerance",
			field:   "Close",
			values:  map[string]float64{"p1": 10.0, "p2": 10.05},
			wantDev: 0.5,
			wantOk:  true,
		},
		{
			name:    "inconsistent beyond tolerance",
			field:   "Close",
			values:  map[string]float64{"p1": 10.0, "p2": 10.5},
			wantDev: 5.0,
			wantOk:  false,
		},
		{
			name:    "single provider always consistent",
			field:   "Close",
			values:  map[string]float64{"p1": 10.0},
			wantDev: 0,
			wantOk:  true,
		},
		{
			name:    "identical values",
			field:   "Close",
			values:  map[string]float64{"p1": 10.0, "p2": 10.0, "p3": 10.0},
			wantDev: 0,
			wantOk:  true,
		},
		{
			name:    "three providers with spread",
			field:   "Close",
			values:  map[string]float64{"p1": 10.0, "p2": 10.03, "p3": 10.06},
			wantDev: 0.6,
			wantOk:  true,
		},
		{
			name:    "zero values",
			field:   "Volume",
			values:  map[string]float64{"p1": 0, "p2": 0},
			wantDev: 0,
			wantOk:  true,
		},
		{
			name:    "one zero one nonzero",
			field:   "Volume",
			values:  map[string]float64{"p1": 0, "p2": 100},
			wantDev: 100,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cc.CheckFloats(tt.field, tt.values)
			if math.Abs(result.MaxDeviation-tt.wantDev) > 0.01 {
				t.Errorf("MaxDeviation = %v, want %v", result.MaxDeviation, tt.wantDev)
			}
			if result.Consistent != tt.wantOk {
				t.Errorf("Consistent = %v, want %v", result.Consistent, tt.wantOk)
			}
		})
	}
}

func TestConsistencyChecker_ToleranceBoundary(t *testing.T) {
	// Exactly at tolerance boundary
	cc := NewConsistencyChecker(1.0)

	result := cc.CheckFloats("Close", map[string]float64{"p1": 100.0, "p2": 101.0})
	if !result.Consistent {
		t.Errorf("1%% deviation should be within 1%% tolerance, got Consistent=%v", result.Consistent)
	}

	result = cc.CheckFloats("Close", map[string]float64{"p1": 100.0, "p2": 101.01})
	if result.Consistent {
		t.Errorf(">1%% deviation should be outside 1%% tolerance, got Consistent=%v", result.Consistent)
	}
}

func TestDeviation(t *testing.T) {
	tests := []struct {
		a, b float64
		want float64
	}{
		{0, 0, 0},
		{10, 10, 0},
		{10, 11, 10},   // 10% deviation
		{100, 101, 1},  // 1% deviation
		{0, 100, 100},  // one zero
		{100, 0, 100},  // one zero
		{-10, -11, 10}, // negative values
	}

	for _, tt := range tests {
		got := deviation(tt.a, tt.b)
		if math.Abs(got-tt.want) > 0.01 {
			t.Errorf("deviation(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestFormatConsistencyReport(t *testing.T) {
	report := &ConsistencyReport{
		Symbol:     "AAPL",
		Providers:  []string{"yahoo", "alphavantage"},
		Consistent: false,
		Fields: []ConsistencyFieldCheck{
			{
				Field:        "Close",
				Values:       map[string]float64{"yahoo": 150.0, "alphavantage": 151.5},
				MaxDeviation: 1.0,
				Consistent:   true,
			},
			{
				Field:        "Volume",
				Values:       map[string]float64{"yahoo": 1000000, "alphavantage": 500000},
				MaxDeviation: 100,
				Consistent:   false,
			},
		},
	}

	s := FormatConsistencyReport(report)
	if len(s) == 0 {
		t.Error("FormatConsistencyReport returned empty string")
	}
}
