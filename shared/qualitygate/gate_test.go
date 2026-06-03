package qualitygate

import (
	"testing"
)

func ptr(v float64) *float64 { return &v }
func boolPtr(v bool) *bool    { return &v }

func TestEvaluate_coldStart(t *testing.T) {
	cfg := GateConfig{
		Metrics: DefaultMetricRules(),
		Baseline: BaselineConfig{
			MinSamples: 3,
		},
		Adaptive: AdaptiveConfig{
			SigmaFactor:      2,
			MaxRegressionPct: 15,
		},
	}

	current := map[string]float64{
		"http_req_duration_p95": 100,
		"http_req_failed_rate":  0.01,
		"http_reqs":             500,
	}

	v := Evaluate(cfg, current, nil)
	if !v.Passed {
		t.Fatalf("expected cold start pass, got: %s", v.Summary)
	}
	if !v.ColdStart {
		t.Fatal("expected cold_start=true")
	}
}

func TestEvaluate_constantFailsOnColdStart(t *testing.T) {
	cfg := GateConfig{
		Metrics: []MetricRule{
			{Name: "http_req_duration_p95", Direction: LowerIsBetter, Max: ptr(50)},
		},
		Baseline: BaselineConfig{MinSamples: 3},
		Adaptive: AdaptiveConfig{SigmaFactor: 2, MaxRegressionPct: 15},
	}

	v := Evaluate(cfg, map[string]float64{"http_req_duration_p95": 100}, nil)
	if v.Passed {
		t.Fatal("expected constant threshold to fail even during cold start")
	}
	if !v.Metrics[0].AdaptiveSkipped {
		t.Fatal("expected adaptive skipped on cold start")
	}
}

func TestEvaluate_constantOnly(t *testing.T) {
	cfg := GateConfig{
		Metrics: []MetricRule{
			{Name: "http_req_duration_p95", Direction: LowerIsBetter, Max: ptr(150)},
		},
		Adaptive: AdaptiveConfig{Enabled: boolPtr(false)},
	}

	v := Evaluate(cfg, map[string]float64{"http_req_duration_p95": 120}, nil)
	if !v.Passed {
		t.Fatalf("expected pass within constant max: %s", v.Metrics[0].Reason)
	}

	v = Evaluate(cfg, map[string]float64{"http_req_duration_p95": 160}, nil)
	if v.Passed {
		t.Fatal("expected fail above constant max")
	}
}

func TestEvaluate_bothConstantAndAdaptive(t *testing.T) {
	cfg := GateConfig{
		Metrics: []MetricRule{
			{Name: "http_req_duration_p95", Direction: LowerIsBetter, Max: ptr(200)},
		},
		Baseline: BaselineConfig{MinSamples: 3},
		Adaptive: AdaptiveConfig{SigmaFactor: 2, MaxRegressionPct: 15},
	}
	history := map[string][]float64{
		"http_req_duration_p95": {100, 102, 98, 101, 99},
	}

	// Passes both constant (≤200) and adaptive (~103)
	v := Evaluate(cfg, map[string]float64{"http_req_duration_p95": 103}, history)
	if !v.Passed {
		t.Fatalf("expected pass: %s", v.Metrics[0].Reason)
	}

	// Fails adaptive but within constant 200
	v = Evaluate(cfg, map[string]float64{"http_req_duration_p95": 120}, history)
	if v.Passed {
		t.Fatal("expected adaptive failure")
	}

	// Passes adaptive but fails constant
	cfg.Metrics[0].Max = ptr(100)
	v = Evaluate(cfg, map[string]float64{"http_req_duration_p95": 102}, history)
	if v.Passed {
		t.Fatal("expected constant failure despite adaptive pass")
	}
}

func TestEvaluate_regressionOnLatency(t *testing.T) {
	cfg := GateConfig{
		Metrics: []MetricRule{
			{Name: "http_req_duration_p95", Direction: LowerIsBetter},
		},
		Baseline: BaselineConfig{MinSamples: 3},
		Adaptive: AdaptiveConfig{SigmaFactor: 2, MaxRegressionPct: 15},
	}

	history := map[string][]float64{
		"http_req_duration_p95": {100, 102, 98, 101, 99},
	}

	v := Evaluate(cfg, map[string]float64{"http_req_duration_p95": 120}, history)
	if v.Passed {
		t.Fatal("expected failure for latency regression")
	}

	v = Evaluate(cfg, map[string]float64{"http_req_duration_p95": 103}, history)
	if !v.Passed {
		t.Fatalf("expected pass within threshold: %s", v.Metrics[0].Reason)
	}
}

func TestEvaluate_regressionOnThroughput(t *testing.T) {
	cfg := GateConfig{
		Metrics: []MetricRule{
			{Name: "http_reqs", Direction: HigherIsBetter},
		},
		Baseline: BaselineConfig{MinSamples: 3},
		Adaptive: AdaptiveConfig{SigmaFactor: 2, MaxRegressionPct: 15},
	}

	history := map[string][]float64{
		"http_reqs": {1000, 1020, 980, 1010, 990},
	}

	v := Evaluate(cfg, map[string]float64{"http_reqs": 800}, history)
	if v.Passed {
		t.Fatal("expected failure for throughput drop")
	}

	v = Evaluate(cfg, map[string]float64{"http_reqs": 970}, history)
	if !v.Passed {
		t.Fatalf("expected pass: %s", v.Metrics[0].Reason)
	}
}

func TestAdaptiveThreshold_lowerIsBetter(t *testing.T) {
	th := adaptiveThreshold(100, 5, LowerIsBetter, AdaptiveConfig{SigmaFactor: 2, MaxRegressionPct: 15})
	if th != 110 {
		t.Fatalf("threshold = %v, want 110", th)
	}
}

func TestAdaptiveThreshold_higherIsBetter(t *testing.T) {
	th := adaptiveThreshold(100, 5, HigherIsBetter, AdaptiveConfig{SigmaFactor: 2, MaxRegressionPct: 15})
	if th != 90 {
		t.Fatalf("threshold = %v, want 90", th)
	}
}

func TestParseReport(t *testing.T) {
	data := []byte(`{"version":1,"tool":"k6","metrics":{"http_req_duration_p95":42.5}}`)
	r, err := ParseReport(data)
	if err != nil {
		t.Fatal(err)
	}
	if r.Metrics["http_req_duration_p95"] != 42.5 {
		t.Fatalf("metric = %v", r.Metrics)
	}
}
