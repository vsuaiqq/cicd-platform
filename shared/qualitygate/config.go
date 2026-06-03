package qualitygate

// Direction defines whether lower or higher metric values are better.
type Direction string

const (
	LowerIsBetter  Direction = "lower_is_better"
	HigherIsBetter Direction = "higher_is_better"
)

// MetricRule configures evaluation for a single metric.
// Use max for lower-is-better metrics, min for higher-is-better metrics (constant ceiling/floor).
type MetricRule struct {
	Name      string    `yaml:"name" json:"name"`
	Direction Direction `yaml:"direction" json:"direction"`
	Max       *float64  `yaml:"max,omitempty" json:"max,omitempty"`
	Min       *float64  `yaml:"min,omitempty" json:"min,omitempty"`
}

func (r MetricRule) HasConstant() bool {
	return r.ConstantThreshold() != nil
}

func (r MetricRule) ConstantThreshold() *float64 {
	switch r.Direction {
	case HigherIsBetter:
		return r.Min
	default:
		return r.Max
	}
}

// BaselineConfig controls historical baseline sampling.
type BaselineConfig struct {
	WindowDays int    `yaml:"window_days" json:"window_days"`
	MinSamples int    `yaml:"min_samples" json:"min_samples"`
	Branch     string `yaml:"branch,omitempty" json:"branch,omitempty"`
}

// AdaptiveConfig controls adaptive threshold calculation.
type AdaptiveConfig struct {
	Enabled          *bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	SigmaFactor      float64 `yaml:"sigma_factor" json:"sigma_factor"`
	MaxRegressionPct float64 `yaml:"max_regression_pct" json:"max_regression_pct"`
}

func (a AdaptiveConfig) IsEnabled() bool {
	if a.Enabled == nil {
		return true
	}
	return *a.Enabled
}

// GateConfig is the full performance gate configuration from pipeline YAML.
type GateConfig struct {
	SourceJob string         `yaml:"source_job" json:"source_job"`
	Metrics   []MetricRule   `yaml:"metrics,omitempty" json:"metrics,omitempty"`
	Baseline  BaselineConfig `yaml:"baseline" json:"baseline"`
	Adaptive  AdaptiveConfig `yaml:"adaptive" json:"adaptive"`
}

func DefaultMetricRules() []MetricRule {
	return []MetricRule{
		{Name: "http_req_duration_p95", Direction: LowerIsBetter},
		{Name: "http_req_failed_rate", Direction: LowerIsBetter},
		{Name: "http_reqs", Direction: HigherIsBetter},
	}
}

func (c *GateConfig) WithDefaults() GateConfig {
	out := *c
	if out.SourceJob == "" {
		out.SourceJob = ""
	}
	if len(out.Metrics) == 0 {
		out.Metrics = DefaultMetricRules()
	}
	if out.Baseline.WindowDays <= 0 {
		out.Baseline.WindowDays = 30
	}
	if out.Baseline.MinSamples <= 0 {
		out.Baseline.MinSamples = 3
	}
	if out.Adaptive.SigmaFactor <= 0 {
		out.Adaptive.SigmaFactor = 2.0
	}
	if out.Adaptive.MaxRegressionPct <= 0 {
		out.Adaptive.MaxRegressionPct = 15
	}
	return out
}
