package qualitygate

import (
	"fmt"
	"strings"
)

// CheckMode describes which threshold types apply to a metric.
type CheckMode string

const (
	CheckModeAdaptive  CheckMode = "adaptive"
	CheckModeConstant  CheckMode = "constant"
	CheckModeBoth      CheckMode = "both"
	CheckModeColdStart CheckMode = "cold_start"
)

// MetricVerdict is the pass/fail result for a single metric comparison.
type MetricVerdict struct {
	Name              string    `json:"name"`
	Direction         Direction `json:"direction"`
	Current           float64   `json:"current"`
	BaselineMean      float64   `json:"baseline_mean"`
	BaselineStdDev    float64   `json:"baseline_stddev"`
	Threshold         float64   `json:"threshold"`
	ConstantThreshold *float64  `json:"constant_threshold,omitempty"`
	AdaptiveThreshold *float64  `json:"adaptive_threshold,omitempty"`
	ConstantPassed    bool      `json:"constant_passed"`
	AdaptivePassed    bool      `json:"adaptive_passed"`
	AdaptiveSkipped   bool      `json:"adaptive_skipped,omitempty"`
	CheckMode         CheckMode `json:"check_mode"`
	Passed            bool      `json:"passed"`
	Reason            string    `json:"reason"`
	ColdStart         bool      `json:"cold_start,omitempty"`
}

// Verdict is the overall performance gate evaluation result.
type Verdict struct {
	Passed          bool            `json:"passed"`
	ColdStart       bool            `json:"cold_start"`
	BaselineSamples int             `json:"baseline_samples"`
	Summary         string          `json:"summary"`
	Metrics         []MetricVerdict `json:"metrics"`
}

// Evaluate compares current metrics against constant and/or adaptive thresholds.
//
// Constant thresholds (per-metric max/min) always apply when configured.
// Adaptive thresholds apply when adaptive is enabled and baseline has min_samples.
// During cold start adaptive checks are skipped; constant checks still apply.
func Evaluate(cfg GateConfig, current map[string]float64, history map[string][]float64) Verdict {
	cfg = cfg.WithDefaults()

	coldStart := isColdStart(cfg, history)
	metrics := make([]MetricVerdict, 0, len(cfg.Metrics))
	allPassed := true

	for _, rule := range cfg.Metrics {
		mv := evaluateMetric(rule, current, history, coldStart, cfg)
		metrics = append(metrics, mv)
		if !mv.Passed {
			allPassed = false
		}
	}

	maxSamples := 0
	for _, rule := range cfg.Metrics {
		if n := len(history[rule.Name]); n > maxSamples {
			maxSamples = n
		}
	}

	return Verdict{
		Passed:          allPassed,
		ColdStart:       coldStart,
		BaselineSamples: maxSamples,
		Summary:         buildSummary(allPassed, coldStart, metrics, cfg.Adaptive.IsEnabled()),
		Metrics:         metrics,
	}
}

func evaluateMetric(
	rule MetricRule,
	current map[string]float64,
	history map[string][]float64,
	coldStart bool,
	cfg GateConfig,
) MetricVerdict {
	value, ok := current[rule.Name]
	if !ok {
		return MetricVerdict{
			Name:      rule.Name,
			Direction: rule.Direction,
			Passed:    false,
			Reason:    "metric missing in current report",
		}
	}

	hasConstant := rule.HasConstant()
	adaptiveEnabled := cfg.Adaptive.IsEnabled()
	samples := history[rule.Name]
	stats := ComputeBaselineStats(samples)

	constantPassed, constThresh := evaluateConstant(rule, value)
	adaptivePassed, adaptThresh, adaptiveSkipped := evaluateAdaptive(rule, value, stats, coldStart, adaptiveEnabled, cfg.Adaptive)

	passed := constantPassed && adaptivePassed
	checkMode := resolveCheckMode(hasConstant, adaptiveEnabled, adaptiveSkipped, coldStart)

	threshold := displayThreshold(constThresh, adaptThresh, rule.Direction)

	return MetricVerdict{
		Name:              rule.Name,
		Direction:         rule.Direction,
		Current:           value,
		BaselineMean:      stats.Mean,
		BaselineStdDev:    stats.StdDev,
		Threshold:         threshold,
		ConstantThreshold: constThresh,
		AdaptiveThreshold: adaptThresh,
		ConstantPassed:    constantPassed,
		AdaptivePassed:    adaptivePassed,
		AdaptiveSkipped:   adaptiveSkipped,
		CheckMode:         checkMode,
		Passed:            passed,
		ColdStart:         coldStart && adaptiveSkipped && adaptiveEnabled,
		Reason:            buildMetricReason(rule, value, stats, constThresh, adaptThresh, constantPassed, adaptivePassed, adaptiveSkipped, cfg.Adaptive),
	}
}

func evaluateConstant(rule MetricRule, value float64) (bool, *float64) {
	th := rule.ConstantThreshold()
	if th == nil {
		return true, nil
	}
	return compare(value, *th, rule.Direction), th
}

func evaluateAdaptive(
	rule MetricRule,
	value float64,
	stats BaselineStats,
	coldStart, adaptiveEnabled bool,
	adaptive AdaptiveConfig,
) (passed bool, threshold *float64, skipped bool) {
	if !adaptiveEnabled || coldStart {
		return true, nil, true
	}
	th := adaptiveThreshold(stats.Mean, stats.StdDev, rule.Direction, adaptive)
	return compare(value, th, rule.Direction), &th, false
}

func resolveCheckMode(hasConstant, adaptiveEnabled, adaptiveSkipped, coldStart bool) CheckMode {
	if coldStart && adaptiveEnabled && adaptiveSkipped {
		if hasConstant {
			return CheckModeBoth
		}
		return CheckModeColdStart
	}
	if hasConstant && adaptiveEnabled && !adaptiveSkipped {
		return CheckModeBoth
	}
	if hasConstant {
		return CheckModeConstant
	}
	return CheckModeAdaptive
}

func displayThreshold(constThresh, adaptThresh *float64, dir Direction) float64 {
	if dir == HigherIsBetter {
		if constThresh != nil && adaptThresh != nil {
			return max(*constThresh, *adaptThresh)
		}
		if adaptThresh != nil {
			return *adaptThresh
		}
		if constThresh != nil {
			return *constThresh
		}
		return 0
	}
	if constThresh != nil && adaptThresh != nil {
		return min(*constThresh, *adaptThresh)
	}
	if adaptThresh != nil {
		return *adaptThresh
	}
	if constThresh != nil {
		return *constThresh
	}
	return 0
}

func isColdStart(cfg GateConfig, history map[string][]float64) bool {
	if !cfg.Adaptive.IsEnabled() {
		return false
	}
	minCount := -1
	for _, rule := range cfg.Metrics {
		n := len(history[rule.Name])
		if minCount < 0 || n < minCount {
			minCount = n
		}
	}
	return minCount < cfg.Baseline.MinSamples
}

func adaptiveThreshold(mean, stddev float64, dir Direction, adaptive AdaptiveConfig) float64 {
	sigma := adaptive.SigmaFactor
	if sigma <= 0 {
		sigma = 2.0
	}
	pct := adaptive.MaxRegressionPct
	if pct <= 0 {
		pct = 15
	}

	switch dir {
	case HigherIsBetter:
		statistical := mean - sigma*stddev
		percentage := mean * (1 - pct/100)
		return max(statistical, percentage)
	default:
		statistical := mean + sigma*stddev
		percentage := mean * (1 + pct/100)
		return min(statistical, percentage)
	}
}

func compare(current, threshold float64, dir Direction) bool {
	switch dir {
	case HigherIsBetter:
		return current >= threshold
	default:
		return current <= threshold
	}
}

func buildMetricReason(
	rule MetricRule,
	current float64,
	stats BaselineStats,
	constThresh, adaptThresh *float64,
	constantPassed, adaptivePassed, adaptiveSkipped bool,
	adaptive AdaptiveConfig,
) string {
	var parts []string
	status := "PASS"
	if !(constantPassed && adaptivePassed) {
		status = "FAIL"
	}

	if constThresh != nil {
		label := "max"
		if rule.Direction == HigherIsBetter {
			label = "min"
		}
		cStatus := "ok"
		if !constantPassed {
			cStatus = "FAIL"
		}
		parts = append(parts, fmt.Sprintf("constant[%s]=%.4f: %s (current=%.4f)", label, *constThresh, cStatus, current))
	}

	if adaptThresh != nil {
		aStatus := "ok"
		if !adaptivePassed {
			aStatus = "FAIL"
		}
		parts = append(parts, fmt.Sprintf(
			"adaptive threshold=%.4f baseline_mean=%.4f baseline_std=%.4f (σ=%.1f, max_regression=%.0f%%): %s (current=%.4f)",
			*adaptThresh, stats.Mean, stats.StdDev, adaptive.SigmaFactor, adaptive.MaxRegressionPct, aStatus, current,
		))
	}

	if adaptiveSkipped && adaptive.IsEnabled() {
		parts = append(parts, "adaptive: skipped (cold start — establishing baseline)")
	} else if adaptiveSkipped {
		parts = append(parts, "adaptive: disabled")
	}

	if len(parts) == 0 {
		return fmt.Sprintf("%s: no thresholds configured", status)
	}
	return fmt.Sprintf("%s: %s", status, strings.Join(parts, "; "))
}

func buildSummary(passed, coldStart bool, metrics []MetricVerdict, adaptiveEnabled bool) string {
	hasConstants := false
	for _, m := range metrics {
		if m.ConstantThreshold != nil {
			hasConstants = true
			break
		}
	}

	if coldStart && adaptiveEnabled {
		if hasConstants {
			if passed {
				return "Performance gate passed: constant thresholds ok; adaptive baseline establishing"
			}
		} else if passed {
			return "Performance gate passed (cold start): adaptive baseline establishing"
		}
	}

	if passed {
		switch {
		case hasConstants && adaptiveEnabled:
			return "Performance gate passed: all metrics within constant and adaptive thresholds"
		case hasConstants:
			return "Performance gate passed: all metrics within constant thresholds"
		default:
			return "Performance gate passed: all metrics within adaptive thresholds"
		}
	}

	var failed []string
	for _, m := range metrics {
		if !m.Passed {
			failed = append(failed, m.Name)
		}
	}
	return fmt.Sprintf("Performance gate failed: threshold violation in %s", strings.Join(failed, ", "))
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
