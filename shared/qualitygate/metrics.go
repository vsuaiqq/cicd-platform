package qualitygate

import (
	"encoding/json"
	"fmt"
	"os"
)

const DefaultMetricsPath = ".flow/perf-metrics.json"

// Report is the standard performance metrics payload produced by load tests.
type Report struct {
	Version int                `json:"version"`
	Tool    string             `json:"tool,omitempty"`
	Metrics map[string]float64 `json:"metrics"`
}

func ParseReport(data []byte) (*Report, error) {
	var r Report
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("qualitygate: parse report: %w", err)
	}
	if len(r.Metrics) == 0 {
		return nil, fmt.Errorf("qualitygate: report has no metrics")
	}
	if r.Version == 0 {
		r.Version = 1
	}
	return &r, nil
}

func LoadReportFromFile(path string) (*Report, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("qualitygate: read %s: %w", path, err)
	}
	return ParseReport(data)
}
