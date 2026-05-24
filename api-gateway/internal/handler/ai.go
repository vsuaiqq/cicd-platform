package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vsuaiqq/cicd/api-gateway/internal/client"
	"github.com/vsuaiqq/cicd/shared/httputil"
)

type AIHandler struct {
	authClient     *client.AuthClient
	pipelineClient *client.PipelineClient
	projectsClient *client.ProjectsClient
	aiClient       *client.AIClient
}

func NewAIHandler(
	authClient *client.AuthClient,
	pipelineClient *client.PipelineClient,
	projectsClient *client.ProjectsClient,
	aiClient *client.AIClient,
) *AIHandler {
	return &AIHandler{
		authClient:     authClient,
		pipelineClient: pipelineClient,
		projectsClient: projectsClient,
		aiClient:       aiClient,
	}
}

func (h *AIHandler) RegisterRoutes(r chi.Router) {
	r.Post("/analyze-failure", h.AnalyzeFailure)
	r.Post("/generate-pipeline", h.GeneratePipeline)
}

type analyzeFailureRequest struct {
	RunID string `json:"run_id"`
	JobID string `json:"job_id"`
	Lang  string `json:"lang"`
}

func (h *AIHandler) AnalyzeFailure(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req analyzeFailureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if req.RunID == "" || req.JobID == "" {
		httputil.Error(w, http.StatusBadRequest, "run_id and job_id are required", nil)
		return
	}

	run, err := h.pipelineClient.GetRunDetails(r.Context(), req.RunID)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "run not found", nil)
			return
		}
		httputil.Error(w, http.StatusBadGateway, "orchestrator unavailable", err)
		return
	}

	if run.ProjectID == "" {
		httputil.Error(w, http.StatusBadGateway, "invalid run payload", nil)
		return
	}
	if !requireProjectAccess(r.Context(), w, h.projectsClient, userID, run.ProjectID) {
		return
	}

	var job *client.JobDetail
	for i := range run.Jobs {
		if run.Jobs[i].ID == req.JobID {
			job = &run.Jobs[i]
			break
		}
	}
	if job == nil {
		httputil.Error(w, http.StatusNotFound, "job not found in this run", nil)
		return
	}

	steps := make([]client.AnalyzeStep, len(job.Steps))
	for i, s := range job.Steps {
		steps[i] = client.AnalyzeStep{
			Name:      s.Name,
			Status:    s.Status,
			ExitCode:  s.ExitCode,
			LogOutput: s.LogOutput,
		}
	}

	jobName := job.DisplayName
	if jobName == "" {
		jobName = job.Name
	}

	lang := req.Lang
	if lang != "en" && lang != "ru" {
		lang = "en"
	}

	analysis, err := h.aiClient.Analyze(r.Context(), &client.AnalyzeRequest{
		JobName:      jobName,
		JobStatus:    job.Status,
		PipelineYAML: run.PipelineYAML,
		Steps:        steps,
		Lang:         lang,
	})
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, analysis)
}

func (h *AIHandler) GeneratePipeline(w http.ResponseWriter, r *http.Request) {
	if requireUserID(w, r, h.authClient) == "" {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if req.Description == "" {
		httputil.Error(w, http.StatusBadRequest, "description is required", nil)
		return
	}

	result, err := h.aiClient.GeneratePipeline(r.Context(), req.Description)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, result)
}
