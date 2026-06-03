package orchestrator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vsuaiqq/cicd/orchestrator-service/internal/db"
	"github.com/vsuaiqq/cicd/orchestrator-service/internal/pipeline"
	orchWs "github.com/vsuaiqq/cicd/orchestrator-service/internal/ws"
	sharedEvents "github.com/vsuaiqq/cicd/shared/events"
	"github.com/vsuaiqq/cicd/shared/qualitygate"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/analytics"
)

type PerformanceGateEvaluator interface {
	EvaluatePerformanceGate(ctx context.Context, req *pb.EvaluatePerformanceGateRequest) (*pb.EvaluatePerformanceGateResponse, error)
}

func (o *Orchestrator) runPerformanceGate(
	ctx context.Context,
	runID, projectID, branch string,
	dbJob *db.PipelineJob,
	plJob *pipeline.Job,
	completedJobName string,
	completedMetrics []sharedEvents.PerformanceMetricValue,
) error {
	pg := plJob.PerformanceGate
	if pg == nil {
		return fmt.Errorf("orchestrator: job %q is not a performance gate", dbJob.Name)
	}

	if err := o.repo.UpdateJobStatus(ctx, dbJob.ID, db.StatusRunning); err != nil {
		return err
	}
	started := time.Now()
	o.hub.Broadcast(runID, orchWs.Event{
		Type:           orchWs.EventJobUpdated,
		RunID:          runID,
		Status:         db.StatusRunning,
		JobName:        dbJob.Name,
		JobID:          dbJob.ID,
		JobStartedAtMs: started.UnixMilli(),
	})

	current := make(map[string]float64)
	if completedJobName == pg.SourceJob && len(completedMetrics) > 0 {
		for _, m := range completedMetrics {
			current[m.Name] = m.Value
		}
	}
	if len(current) == 0 {
		return o.finishPerformanceGate(ctx, runID, projectID, dbJob, false, "performance gate failed: no metrics from source job "+pg.SourceJob)
	}

	if o.perfGate == nil {
		return o.finishPerformanceGate(ctx, runID, projectID, dbJob, false, "performance gate failed: analytics service unavailable")
	}

	req := buildEvaluateGateRequest(projectID, runID, branch, dbJob.Name, pg, current)
	resp, err := o.perfGate.EvaluatePerformanceGate(ctx, req)
	if err != nil {
		return o.finishPerformanceGate(ctx, runID, projectID, dbJob, false, "performance gate evaluation error: "+err.Error())
	}

	stepOutput := formatGateOutput(resp)
	status := db.StatusSuccess
	if !resp.Passed {
		status = db.StatusFailed
	}

	finished := time.Now()
	_ = o.repo.UpsertStep(ctx, &db.PipelineStep{
		JobID:      dbJob.ID,
		StepIndex:  0,
		Name:       "Adaptive performance quality gate",
		Status:     status,
		LogOutput:  stepOutput,
		ExitCode:   gateExitCode(resp.Passed),
		StartedAt:  &started,
		FinishedAt: &finished,
	})

	return o.finishPerformanceGateJob(ctx, runID, projectID, dbJob, status, stepOutput, resp.Passed)
}

func buildEvaluateGateRequest(
	projectID, runID, branch, gateJobName string,
	pg *pipeline.PerformanceGateConfig,
	current map[string]float64,
) *pb.EvaluatePerformanceGateRequest {
	cfg := pipelineGateToQualityGate(pg)
	cfg = cfg.WithDefaults()

	req := &pb.EvaluatePerformanceGateRequest{
		ProjectId:   projectID,
		RunId:       runID,
		JobName:     gateJobName,
		SourceJob:   pg.SourceJob,
		Branch:      branch,
		Baseline: &pb.PerformanceBaselineConfig{
			WindowDays: int32(cfg.Baseline.WindowDays),
			MinSamples: int32(cfg.Baseline.MinSamples),
			Branch:     cfg.Baseline.Branch,
		},
	}
	enabled := cfg.Adaptive.IsEnabled()
	req.Adaptive = &pb.PerformanceAdaptiveConfig{
		Enabled:          enabled,
		SigmaFactor:      cfg.Adaptive.SigmaFactor,
		MaxRegressionPct: cfg.Adaptive.MaxRegressionPct,
	}
	for name, value := range current {
		req.CurrentMetrics = append(req.CurrentMetrics, &pb.PerformanceMetricInput{Name: name, Value: value})
	}
	for _, rule := range cfg.Metrics {
		dir := string(rule.Direction)
		if dir == "" {
			dir = string(qualitygate.LowerIsBetter)
		}
		pr := &pb.PerformanceMetricRule{Name: rule.Name, Direction: dir}
		if rule.Max != nil {
			pr.Max = rule.Max
		}
		if rule.Min != nil {
			pr.Min = rule.Min
		}
		req.MetricRules = append(req.MetricRules, pr)
	}
	return req
}

func pipelineGateToQualityGate(pg *pipeline.PerformanceGateConfig) qualitygate.GateConfig {
	cfg := qualitygate.GateConfig{SourceJob: pg.SourceJob}
	for _, m := range pg.Metrics {
		dir := qualitygate.Direction(m.Direction)
		if dir == "" {
			dir = qualitygate.LowerIsBetter
		}
		cfg.Metrics = append(cfg.Metrics, qualitygate.MetricRule{
			Name:      m.Name,
			Direction: dir,
			Max:       m.Max,
			Min:       m.Min,
		})
	}
	cfg.Baseline = qualitygate.BaselineConfig{
		WindowDays: pg.Baseline.WindowDays,
		MinSamples: pg.Baseline.MinSamples,
		Branch:     pg.Baseline.Branch,
	}
	cfg.Adaptive = qualitygate.AdaptiveConfig{
		Enabled:          pg.Adaptive.Enabled,
		SigmaFactor:      pg.Adaptive.SigmaFactor,
		MaxRegressionPct: pg.Adaptive.MaxRegressionPct,
	}
	return cfg
}

func (o *Orchestrator) finishPerformanceGate(ctx context.Context, runID, projectID string, dbJob *db.PipelineJob, passed bool, summary string) error {
	started := time.Now()
	finished := started
	status := db.StatusFailed
	if passed {
		status = db.StatusSuccess
	}
	_ = o.repo.UpsertStep(ctx, &db.PipelineStep{
		JobID:      dbJob.ID,
		StepIndex:  0,
		Name:       "Adaptive performance quality gate",
		Status:     status,
		LogOutput:  summary,
		ExitCode:   gateExitCode(passed),
		StartedAt:  &started,
		FinishedAt: &finished,
	})
	return o.finishPerformanceGateJob(ctx, runID, projectID, dbJob, status, summary, passed)
}

func (o *Orchestrator) finishPerformanceGateJob(
	ctx context.Context,
	runID, projectID string,
	dbJob *db.PipelineJob,
	status, output string,
	passed bool,
) error {
	if err := o.repo.UpdateJobStatus(ctx, dbJob.ID, status); err != nil {
		return err
	}

	o.hub.Broadcast(runID, orchWs.Event{
		Type: orchWs.EventJobUpdated,
		RunID: runID,
		Status: status,
		JobName: dbJob.Name,
		JobID: dbJob.ID,
		Steps: []orchWs.StepEvent{{
			Index: 0, Name: "Adaptive performance quality gate", Status: status,
			ExitCode: gateExitCode(passed), LogOutput: output,
		}},
		JobFinishedAtMs: time.Now().UnixMilli(),
	})

	// Re-use job result handling to schedule dependents or fail the run.
	result := &sharedEvents.PipelineJobResult{
		RunID: runID, JobID: dbJob.ID, JobName: dbJob.Name,
		ProjectID: projectID, Status: status,
		Steps: []sharedEvents.StepResult{{
			Index: 0, Name: "Adaptive performance quality gate",
			Status: status, Output: output, ExitCode: gateExitCode(passed),
		}},
		StartedAt: time.Now().Unix(), FinishedAt: time.Now().Unix(),
	}
	return o.HandleJobResult(ctx, result)
}

func gateExitCode(passed bool) int {
	if passed {
		return 0
	}
	return 1
}

func formatGateOutput(resp *pb.EvaluatePerformanceGateResponse) string {
	var b strings.Builder
	b.WriteString(resp.Summary)
	b.WriteString("\n\n")
	if resp.ColdStart {
		b.WriteString("Mode: cold start (establishing baseline)\n")
	} else {
		b.WriteString(fmt.Sprintf("Baseline samples: %d\n", resp.BaselineSamples))
	}
	b.WriteString("\nMetric comparison:\n")
	for _, m := range resp.Metrics {
		status := "PASS"
		if !m.Passed {
			status = "FAIL"
		}
		b.WriteString(fmt.Sprintf("  [%s] %s: current=%.4f threshold=%.4f baseline_mean=%.4f\n",
			status, m.Name, m.Current, m.Threshold, m.BaselineMean))
		b.WriteString(fmt.Sprintf("         %s\n", m.Reason))
	}
	return b.String()
}
