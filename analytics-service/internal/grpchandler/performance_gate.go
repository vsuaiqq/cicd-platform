package grpchandler

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vsuaiqq/cicd/analytics-service/internal/db"
	"github.com/vsuaiqq/cicd/shared/logger"
	"github.com/vsuaiqq/cicd/shared/qualitygate"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/analytics"
)

func (s *AnalyticsServer) EvaluatePerformanceGate(ctx context.Context, req *pb.EvaluatePerformanceGateRequest) (*pb.EvaluatePerformanceGateResponse, error) {
	if req.ProjectId == "" || req.RunId == "" || req.JobName == "" || req.SourceJob == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id, run_id, job_name and source_job are required")
	}
	if len(req.CurrentMetrics) == 0 {
		return nil, status.Error(codes.InvalidArgument, "current_metrics are required")
	}

	cfg := protoToGateConfig(req)
	current := make(map[string]float64, len(req.CurrentMetrics))
	for _, m := range req.CurrentMetrics {
		current[m.Name] = m.Value
	}

	since := time.Now().UTC().AddDate(0, 0, -cfg.Baseline.WindowDays)
	metricNames := make([]string, len(cfg.Metrics))
	for i, rule := range cfg.Metrics {
		metricNames[i] = rule.Name
	}

	branch := req.Branch
	if cfg.Baseline.Branch != "" {
		branch = cfg.Baseline.Branch
	}

	history, err := s.repo.GetBaselineMetricValues(ctx, req.ProjectId, req.SourceJob, branch, req.RunId, metricNames, since, 50)
	if err != nil {
		logger.L().Error().Err(err).Msg("GetBaselineMetricValues error")
		return nil, status.Error(codes.Internal, "failed to load baseline")
	}

	verdict := qualitygate.Evaluate(cfg, current, history)

	detailsJSON, err := db.EncodeGateDetails(verdict.Metrics)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to encode gate details")
	}

	now := time.Now()
	if err := s.repo.InsertPerformanceGateResult(ctx, &db.PerformanceGateResult{
		ProjectID:       req.ProjectId,
		RunID:           req.RunId,
		JobName:         req.JobName,
		Passed:          verdict.Passed,
		ColdStart:       verdict.ColdStart,
		BaselineSamples: uint32(verdict.BaselineSamples),
		Summary:         verdict.Summary,
		DetailsJSON:     detailsJSON,
		CreatedAt:       now,
	}); err != nil {
		logger.L().Error().Err(err).Msg("InsertPerformanceGateResult error")
		return nil, status.Error(codes.Internal, "failed to store gate result")
	}

	if err := s.repo.MarkRunMetricsGatePassed(ctx, req.ProjectId, req.RunId, req.SourceJob, verdict.Passed); err != nil {
		logger.L().Error().Err(err).Msg("MarkRunMetricsGatePassed error")
		return nil, status.Error(codes.Internal, "failed to update baseline eligibility")
	}

	return verdictToProto(verdict), nil
}

func (s *AnalyticsServer) GetPerformanceGateResult(ctx context.Context, req *pb.GetPerformanceGateResultRequest) (*pb.GetPerformanceGateResultResponse, error) {
	if req.RunId == "" || req.JobName == "" {
		return nil, status.Error(codes.InvalidArgument, "run_id and job_name are required")
	}

	result, err := s.repo.GetPerformanceGateResult(ctx, req.RunId, req.JobName)
	if err != nil {
		logger.L().Error().Err(err).Msg("GetPerformanceGateResult error")
		return nil, status.Error(codes.Internal, "failed to load gate result")
	}
	if result == nil {
		return &pb.GetPerformanceGateResultResponse{Found: false}, nil
	}

	resp := &pb.GetPerformanceGateResultResponse{
		Found:            true,
		Passed:           result.Passed,
		ColdStart:        result.ColdStart,
		BaselineSamples:  int32(result.BaselineSamples),
		Summary:          result.Summary,
		EvaluatedAt:      result.CreatedAt.UTC().Format(time.RFC3339),
	}
	if result.DetailsJSON != "" {
		var metrics []qualitygate.MetricVerdict
		if err := json.Unmarshal([]byte(result.DetailsJSON), &metrics); err == nil {
			resp.Metrics = metricVerdictsToProto(metrics)
		}
	}
	return resp, nil
}

func protoToGateConfig(req *pb.EvaluatePerformanceGateRequest) qualitygate.GateConfig {
	cfg := qualitygate.GateConfig{
		SourceJob: req.SourceJob,
	}
	for _, rule := range req.MetricRules {
		mr := qualitygate.MetricRule{
			Name:      rule.Name,
			Direction: qualitygate.Direction(rule.Direction),
		}
		if rule.Max != nil {
			v := rule.GetMax()
			mr.Max = &v
		}
		if rule.Min != nil {
			v := rule.GetMin()
			mr.Min = &v
		}
		cfg.Metrics = append(cfg.Metrics, mr)
	}
	if req.Baseline != nil {
		cfg.Baseline = qualitygate.BaselineConfig{
			WindowDays: int(req.Baseline.WindowDays),
			MinSamples: int(req.Baseline.MinSamples),
			Branch:     req.Baseline.Branch,
		}
	}
	if req.Adaptive != nil {
		enabled := req.Adaptive.Enabled
		cfg.Adaptive = qualitygate.AdaptiveConfig{
			Enabled:          &enabled,
			SigmaFactor:      req.Adaptive.SigmaFactor,
			MaxRegressionPct: req.Adaptive.MaxRegressionPct,
		}
	}
	return cfg.WithDefaults()
}

func verdictToProto(v qualitygate.Verdict) *pb.EvaluatePerformanceGateResponse {
	return &pb.EvaluatePerformanceGateResponse{
		Passed:           v.Passed,
		ColdStart:        v.ColdStart,
		BaselineSamples:  int32(v.BaselineSamples),
		Summary:          v.Summary,
		Metrics:          metricVerdictsToProto(v.Metrics),
	}
}

func metricVerdictsToProto(metrics []qualitygate.MetricVerdict) []*pb.PerformanceMetricVerdict {
	out := make([]*pb.PerformanceMetricVerdict, 0, len(metrics))
	for _, m := range metrics {
		pv := &pb.PerformanceMetricVerdict{
			Name:            m.Name,
			Direction:       string(m.Direction),
			Current:         m.Current,
			BaselineMean:    m.BaselineMean,
			BaselineStddev:  m.BaselineStdDev,
			Threshold:       m.Threshold,
			Passed:          m.Passed,
			Reason:          m.Reason,
			ColdStart:       m.ColdStart,
			ConstantPassed:  m.ConstantPassed,
			AdaptivePassed:  m.AdaptivePassed,
			AdaptiveSkipped: m.AdaptiveSkipped,
			CheckMode:       string(m.CheckMode),
		}
		if m.ConstantThreshold != nil {
			v := *m.ConstantThreshold
			pv.ConstantThreshold = &v
		}
		if m.AdaptiveThreshold != nil {
			v := *m.AdaptiveThreshold
			pv.AdaptiveThreshold = &v
		}
		out = append(out, pv)
	}
	return out
}
