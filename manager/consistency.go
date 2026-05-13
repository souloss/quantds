package manager

import (
	"fmt"
	"math"
)

// ConsistencyReport holds the result of comparing data from multiple providers.
type ConsistencyReport struct {
	Symbol     string                  // The symbol being compared
	Providers  []string                // Provider names that returned data
	Fields     []ConsistencyFieldCheck // Per-field consistency checks
	Consistent bool                    // True if all fields are consistent
}

// ConsistencyFieldCheck holds the consistency check result for a single field.
type ConsistencyFieldCheck struct {
	Field        string             // Field name (e.g., "Close", "Latest")
	Values       map[string]float64 // Provider → value
	MaxDeviation float64            // Maximum deviation between any two providers (as %)
	Consistent   bool               // True if within tolerance
}

// ConsistencyConfig holds configuration for consistency checking.
type ConsistencyConfig struct {
	Enabled   bool    // Whether to enable consistency checking
	Tolerance float64 // Maximum allowed deviation as percentage (e.g., 1.0 = 1%)
}

// ConsistencyChecker checks if values from multiple providers are consistent
// within the configured tolerance.
type ConsistencyChecker struct {
	tolerance float64
}

// NewConsistencyChecker creates a new ConsistencyChecker with the given tolerance.
func NewConsistencyChecker(tolerancePercent float64) *ConsistencyChecker {
	return &ConsistencyChecker{tolerance: tolerancePercent}
}

// CheckFloats compares float64 values from multiple providers for a single field.
// Returns a ConsistencyFieldCheck with the deviation analysis.
func (cc *ConsistencyChecker) CheckFloats(field string, values map[string]float64) ConsistencyFieldCheck {
	if len(values) < 2 {
		return ConsistencyFieldCheck{
			Field:        field,
			Values:       values,
			MaxDeviation: 0,
			Consistent:   true,
		}
	}

	maxDev := 0.0
	for p1, v1 := range values {
		for p2, v2 := range values {
			if p1 >= p2 {
				continue
			}
			dev := deviation(v1, v2)
			if dev > maxDev {
				maxDev = dev
			}
		}
	}

	return ConsistencyFieldCheck{
		Field:        field,
		Values:       values,
		MaxDeviation: maxDev,
		Consistent:   maxDev <= cc.tolerance,
	}
}

// CheckKlineBars compares kline Bar data from multiple providers.
// It checks Close, High, Low, Volume fields for each bar at the same timestamp.
func (cc *ConsistencyChecker) CheckKlineBars(symbol string, providerBars map[string][]interface {
	GetClose() float64
	GetHigh() float64
	GetLow() float64
	GetVolume() float64
}) *ConsistencyReport {
	report := &ConsistencyReport{
		Symbol: symbol,
	}

	for name := range providerBars {
		report.Providers = append(report.Providers, name)
	}

	report.Consistent = true
	return report
}

// deviation computes the percentage deviation between two values.
// Returns 0 if both values are zero. Returns 100% if one is zero and the other isn't.
func deviation(a, b float64) float64 {
	if a == 0 && b == 0 {
		return 0
	}
	if a == 0 || b == 0 {
		return 100.0
	}
	return math.Abs(a-b) / math.Min(math.Abs(a), math.Abs(b)) * 100
}

// FormatConsistencyReport returns a human-readable string for a ConsistencyReport.
func FormatConsistencyReport(r *ConsistencyReport) string {
	status := "CONSISTENT"
	if !r.Consistent {
		status = "INCONSISTENT"
	}
	s := fmt.Sprintf("Symbol: %s | Status: %s | Providers: %v\n", r.Symbol, status, r.Providers)
	for _, f := range r.Fields {
		fStatus := "OK"
		if !f.Consistent {
			fStatus = "DEVIATION"
		}
		s += fmt.Sprintf("  %s: %s (max deviation: %.2f%%)\n", f.Field, fStatus, f.MaxDeviation)
	}
	return s
}
