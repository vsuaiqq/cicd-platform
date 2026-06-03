package runner

import (
	"path/filepath"

	"github.com/vsuaiqq/cicd/shared/events"
	"github.com/vsuaiqq/cicd/shared/qualitygate"
)

func extractPerformanceMetrics(repoDir string) []events.PerformanceMetricValue {
	reportPath := filepath.Join(repoDir, qualitygate.DefaultMetricsPath)
	report, err := qualitygate.LoadReportFromFile(reportPath)
	if err != nil {
		return nil
	}
	out := make([]events.PerformanceMetricValue, 0, len(report.Metrics))
	for name, value := range report.Metrics {
		out = append(out, events.PerformanceMetricValue{Name: name, Value: value})
	}
	return out
}
