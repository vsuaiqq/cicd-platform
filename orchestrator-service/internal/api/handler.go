package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/vsuaiqq/cicd/orchestrator-service/internal/db"
	"github.com/vsuaiqq/cicd/orchestrator-service/internal/pipeline"
)

type RunController interface {
	CancelRun(ctx context.Context, runID string) error
	ApproveJob(ctx context.Context, runID, jobID string) error
	RejectJob(ctx context.Context, runID, jobID string) error
}

type Handler struct {
	repo       *db.Repository
	controller RunController
}

func NewHandler(repo *db.Repository, controller RunController) *Handler {
	return &Handler{repo: repo, controller: controller}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		r.Get("/runs", h.ListRuns)
		r.Get("/runs/{runID}", h.GetRun)
		r.Post("/runs/{runID}/cancel", h.CancelRun)
		r.Post("/runs/{runID}/jobs/{jobID}/approve", h.ApproveJob)
		r.Post("/runs/{runID}/jobs/{jobID}/reject", h.RejectJob)
	})
}

type runSummary struct {
	ID          string  `json:"id"`
	ProjectID   string  `json:"project_id"`
	CommitSHA   string  `json:"commit_sha"`
	Branch      string  `json:"branch"`
	TriggerType string  `json:"trigger_type"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	StartedAt   *string `json:"started_at"`
	FinishedAt  *string `json:"finished_at"`
}

type stepResponse struct {
	Index      int     `json:"index"`
	Name       string  `json:"name"`
	Status     string  `json:"status"`
	LogOutput  string  `json:"log_output"`
	ExitCode   int     `json:"exit_code"`
	StartedAt  *string `json:"started_at"`
	FinishedAt *string `json:"finished_at"`
}

type jobResponse struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	DisplayName string         `json:"display_name"`
	Status      string         `json:"status"`
	Attempt     int            `json:"attempt"`
	StartedAt   *string        `json:"started_at"`
	FinishedAt  *string        `json:"finished_at"`
	Steps       []stepResponse `json:"steps"`
}

type runDetailResponse struct {
	runSummary
	PipelineYAML string        `json:"pipeline_yaml,omitempty"`
	Jobs         []jobResponse `json:"jobs"`
}

func (h *Handler) ListRuns(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		jsonError(w, http.StatusBadRequest, "project_id is required")
		return
	}
	runs, err := h.repo.ListRunsByProject(r.Context(), projectID, 20)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list runs")
		return
	}
	out := make([]runSummary, 0, len(runs))
	for _, r := range runs {
		out = append(out, toRunSummary(r))
	}
	jsonOK(w, out)
}

func (h *Handler) GetRun(w http.ResponseWriter, r *http.Request) {
	run, err := h.repo.GetRun(r.Context(), chi.URLParam(r, "runID"))
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to get run")
		return
	}
	if run == nil {
		jsonError(w, http.StatusNotFound, "run not found")
		return
	}

	jobs, err := h.repo.GetJobsByRun(r.Context(), run.ID)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to get jobs")
		return
	}

	jobsOut := make([]jobResponse, 0, len(jobs))
	for _, j := range jobs {
		steps, _ := h.repo.GetStepsByJob(r.Context(), j.ID)
		jobsOut = append(jobsOut, toJobResponse(j, steps))
	}
	sortJobsByExecutionOrder(jobsOut, run.PipelineYAML)

	jsonOK(w, runDetailResponse{
		runSummary:   toRunSummary(run),
		PipelineYAML: run.PipelineYAML,
		Jobs:         jobsOut,
	})
}

func (h *Handler) CancelRun(w http.ResponseWriter, r *http.Request) {
	runID := chi.URLParam(r, "runID")
	if err := h.controller.CancelRun(r.Context(), runID); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ApproveJob(w http.ResponseWriter, r *http.Request) {
	runID := chi.URLParam(r, "runID")
	jobID := chi.URLParam(r, "jobID")
	if err := h.controller.ApproveJob(r.Context(), runID, jobID); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RejectJob(w http.ResponseWriter, r *http.Request) {
	runID := chi.URLParam(r, "runID")
	jobID := chi.URLParam(r, "jobID")
	if err := h.controller.RejectJob(r.Context(), runID, jobID); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toRunSummary(r *db.PipelineRun) runSummary {
	return runSummary{
		ID:          r.ID,
		ProjectID:   r.ProjectID,
		CommitSHA:   r.CommitSHA,
		Branch:      r.Branch,
		TriggerType: r.TriggerType,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		StartedAt:   fmtTime(r.StartedAt),
		FinishedAt:  fmtTime(r.FinishedAt),
	}
}

func toJobResponse(j *db.PipelineJob, steps []*db.PipelineStep) jobResponse {
	sr := make([]stepResponse, 0, len(steps))
	for _, s := range steps {
		sr = append(sr, stepResponse{
			Index:      s.StepIndex,
			Name:       s.Name,
			Status:     s.Status,
			LogOutput:  s.LogOutput,
			ExitCode:   s.ExitCode,
			StartedAt:  fmtTime(s.StartedAt),
			FinishedAt: fmtTime(s.FinishedAt),
		})
	}
	return jobResponse{
		ID:          j.ID,
		Name:        j.Name,
		DisplayName: j.DisplayName,
		Status:      j.Status,
		Attempt:     j.Attempt,
		StartedAt:   fmtTime(j.StartedAt),
		FinishedAt:  fmtTime(j.FinishedAt),
		Steps:       sr,
	}
}

func fmtTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

func sortJobsByExecutionOrder(jobs []jobResponse, pipelineYAML string) {
	if pipelineYAML == "" || len(jobs) <= 1 {
		return
	}
	order, err := pipeline.ExecutionOrder([]byte(pipelineYAML))
	if err != nil || len(order) == 0 {
		return
	}
	rank := make(map[string]int, len(order))
	for i, name := range order {
		rank[name] = i
	}
	sort.Slice(jobs, func(i, j int) bool {
		return rank[jobs[i].Name] < rank[jobs[j].Name]
	})
}

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
