package grpchandler

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vsuaiqq/cicd/analytics-service/internal/db"
	"github.com/vsuaiqq/cicd/shared/logger"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/analytics"
)

type AnalyticsServer struct {
	pb.UnimplementedAnalyticsServiceServer
	repo *db.Repository
}

func NewAnalyticsServer(repo *db.Repository) *AnalyticsServer {
	return &AnalyticsServer{repo: repo}
}

func (s *AnalyticsServer) GetDashboard(ctx context.Context, req *pb.DashboardRequest) (*pb.DashboardResponse, error) {
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id is required")
	}

	since, err := parsePeriod(req.Period)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid period %q: use 7d, 30d or 90d", req.Period)
	}

	resp := &pb.DashboardResponse{}

	runStats, err := s.repo.GetRunStats(ctx, req.ProjectId, since)
	if err != nil {
		logger.L().Error().Err(err).Msg("GetRunStats error")
		return nil, status.Error(codes.Internal, "failed to query run stats")
	}

	resp.TotalRuns = runStats.Total
	resp.SuccessCount = runStats.Success
	resp.FailedCount = runStats.Failed
	resp.CancelledCount = runStats.Cancelled
	resp.AvgDurationSec = runStats.AvgDur
	resp.P50DurationSec = runStats.P50Dur
	resp.P95DurationSec = runStats.P95Dur

	if runStats.Total > 0 {
		resp.SuccessRate = float64(runStats.Success) / float64(runStats.Total)
	}

	trend, err := s.repo.GetDailyTrend(ctx, req.ProjectId, since)
	if err != nil {
		logger.L().Error().Err(err).Msg("GetDailyTrend error")
		return nil, status.Error(codes.Internal, "failed to query trend")
	}
	for _, p := range trend {
		resp.Trend = append(resp.Trend, &pb.DailyPoint{
			Date:    p.Date,
			Total:   p.Total,
			Success: p.Success,
			Failed:  p.Failed,
		})
	}

	const top = 5
	const minRuns = 3

	failing, err := s.repo.GetTopFailingJobs(ctx, req.ProjectId, since, top, minRuns)
	if err != nil {
		logger.L().Error().Err(err).Msg("GetTopFailingJobs error")
		return nil, status.Error(codes.Internal, "failed to query failing jobs")
	}
	resp.TopFailingJobs = toProtoJobStats(failing)

	slow, err := s.repo.GetTopSlowJobs(ctx, req.ProjectId, since, top, minRuns)
	if err != nil {
		logger.L().Error().Err(err).Msg("GetTopSlowJobs error")
		return nil, status.Error(codes.Internal, "failed to query slow jobs")
	}
	resp.TopSlowJobs = toProtoJobStats(slow)

	flaky, err := s.repo.GetFlakyJobs(ctx, req.ProjectId, since, top, minRuns)
	if err != nil {
		logger.L().Error().Err(err).Msg("GetFlakyJobs error")
		return nil, status.Error(codes.Internal, "failed to query flaky jobs")
	}
	resp.FlakyJobs = toProtoJobStats(flaky)

	return resp, nil
}

func parsePeriod(period string) (time.Time, error) {
	now := time.Now().UTC()
	switch period {
	case "", "7d":
		return now.AddDate(0, 0, -7), nil
	case "30d":
		return now.AddDate(0, 0, -30), nil
	case "90d":
		return now.AddDate(0, 0, -90), nil
	default:
		return time.Time{}, status.Errorf(codes.InvalidArgument, "period must be 7d, 30d, or 90d")
	}
}

func toProtoJobStats(stats []*db.JobStat) []*pb.JobStat {
	out := make([]*pb.JobStat, 0, len(stats))
	for _, s := range stats {
		out = append(out, &pb.JobStat{
			JobName:        s.JobName,
			TotalRuns:      s.TotalRuns,
			FailureRate:    s.FailureRate,
			AvgDurationSec: s.AvgDurationSec,
			AvgAttempts:    s.AvgAttempts,
		})
	}
	return out
}
