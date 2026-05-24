package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vsuaiqq/cicd/api-gateway/internal/client"
	"github.com/vsuaiqq/cicd/shared/httputil"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/analytics"
)

type dashboardJSON struct {
	TotalRuns      int64            `json:"total_runs"`
	SuccessCount   int64            `json:"success_count"`
	FailedCount    int64            `json:"failed_count"`
	CancelledCount int64            `json:"cancelled_count"`
	SuccessRate    float64          `json:"success_rate"`
	AvgDurationSec float64          `json:"avg_duration_sec"`
	P50DurationSec float64          `json:"p50_duration_sec"`
	P95DurationSec float64          `json:"p95_duration_sec"`
	Trend          []dailyPointJSON `json:"trend"`
	TopFailingJobs []jobStatJSON    `json:"top_failing_jobs"`
	TopSlowJobs    []jobStatJSON    `json:"top_slow_jobs"`
	FlakyJobs      []jobStatJSON    `json:"flaky_jobs"`
}

type dailyPointJSON struct {
	Date    string `json:"date"`
	Total   int32  `json:"total"`
	Success int32  `json:"success"`
	Failed  int32  `json:"failed"`
}

type jobStatJSON struct {
	JobName        string  `json:"job_name"`
	TotalRuns      int32   `json:"total_runs"`
	FailureRate    float64 `json:"failure_rate"`
	AvgDurationSec float64 `json:"avg_duration_sec"`
	AvgAttempts    float64 `json:"avg_attempts"`
}

func toDashboardJSON(r *pb.DashboardResponse) dashboardJSON {
	trend := make([]dailyPointJSON, 0, len(r.Trend))
	for _, p := range r.Trend {
		trend = append(trend, dailyPointJSON{Date: p.Date, Total: p.Total, Success: p.Success, Failed: p.Failed})
	}
	return dashboardJSON{
		TotalRuns:      r.TotalRuns,
		SuccessCount:   r.SuccessCount,
		FailedCount:    r.FailedCount,
		CancelledCount: r.CancelledCount,
		SuccessRate:    r.SuccessRate,
		AvgDurationSec: r.AvgDurationSec,
		P50DurationSec: r.P50DurationSec,
		P95DurationSec: r.P95DurationSec,
		Trend:          trend,
		TopFailingJobs: toJobStatJSON(r.TopFailingJobs),
		TopSlowJobs:    toJobStatJSON(r.TopSlowJobs),
		FlakyJobs:      toJobStatJSON(r.FlakyJobs),
	}
}

func toJobStatJSON(stats []*pb.JobStat) []jobStatJSON {
	out := make([]jobStatJSON, 0, len(stats))
	for _, s := range stats {
		out = append(out, jobStatJSON{
			JobName:        s.JobName,
			TotalRuns:      s.TotalRuns,
			FailureRate:    s.FailureRate,
			AvgDurationSec: s.AvgDurationSec,
			AvgAttempts:    s.AvgAttempts,
		})
	}
	return out
}

type AnalyticsHandler struct {
	authClient      *client.AuthClient
	analyticsClient *client.AnalyticsClient
	projectsClient  *client.ProjectsClient
}

func NewAnalyticsHandler(authClient *client.AuthClient, analyticsClient *client.AnalyticsClient, projectsClient *client.ProjectsClient) *AnalyticsHandler {
	return &AnalyticsHandler{
		authClient:      authClient,
		analyticsClient: analyticsClient,
		projectsClient:  projectsClient,
	}
}

func (h *AnalyticsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/analytics/dashboard", h.GetDashboard)
}

func (h *AnalyticsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		httputil.Error(w, http.StatusBadRequest, "project_id is required", nil)
		return
	}
	if !requireProjectAccess(r.Context(), w, h.projectsClient, userID, projectID) {
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	resp, err := h.analyticsClient.GetDashboard(r.Context(), projectID, period)
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, "analytics service unavailable", err)
		return
	}

	httputil.JSON(w, http.StatusOK, toDashboardJSON(resp))
}
