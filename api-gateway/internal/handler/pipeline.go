package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vsuaiqq/cicd/api-gateway/internal/client"
	"github.com/vsuaiqq/cicd/shared/httputil"
)

type PipelineHandler struct {
	authClient     *client.AuthClient
	pipelineClient *client.PipelineClient
	projectsClient *client.ProjectsClient
}

func NewPipelineHandler(authClient *client.AuthClient, pipelineClient *client.PipelineClient, projectsClient *client.ProjectsClient) *PipelineHandler {
	return &PipelineHandler{
		authClient:     authClient,
		pipelineClient: pipelineClient,
		projectsClient: projectsClient,
	}
}

func (h *PipelineHandler) RegisterRoutes(r chi.Router) {
	r.Route("/pipeline", func(r chi.Router) {
		r.Get("/runs", h.ListRuns)
		r.Get("/runs/{runID}", h.GetRun)
		r.Post("/runs/{runID}/cancel", h.CancelRun)
		r.Post("/runs/{runID}/jobs/{jobID}/approve", h.ApproveJob)
		r.Post("/runs/{runID}/jobs/{jobID}/reject", h.RejectJob)
	})
}

func (h *PipelineHandler) ListRuns(w http.ResponseWriter, r *http.Request) {
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
	runs, err := h.pipelineClient.ListRuns(r.Context(), projectID)
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, "orchestrator unavailable", err)
		return
	}
	httputil.JSON(w, http.StatusOK, runs)
}

func (h *PipelineHandler) GetRun(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}
	run, err := h.pipelineClient.GetRun(r.Context(), chi.URLParam(r, "runID"))
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "run not found", nil)
			return
		}
		httputil.Error(w, http.StatusBadGateway, "orchestrator unavailable", err)
		return
	}
	projID, ok := projectIDFromRunMap(run)
	if !ok {
		httputil.Error(w, http.StatusBadGateway, "invalid run payload", nil)
		return
	}
	if !requireProjectAccess(r.Context(), w, h.projectsClient, userID, projID) {
		return
	}
	httputil.JSON(w, http.StatusOK, run)
}

func (h *PipelineHandler) CancelRun(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}
	run, err := h.pipelineClient.GetRun(r.Context(), chi.URLParam(r, "runID"))
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "run not found", nil)
			return
		}
		httputil.Error(w, http.StatusBadGateway, "orchestrator unavailable", err)
		return
	}
	projID, ok := projectIDFromRunMap(run)
	if !ok {
		httputil.Error(w, http.StatusBadGateway, "invalid run payload", nil)
		return
	}
	if !requireProjectAccess(r.Context(), w, h.projectsClient, userID, projID) {
		return
	}
	if err := h.pipelineClient.CancelRun(r.Context(), chi.URLParam(r, "runID")); err != nil {
		if errors.Is(err, client.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "run not found", nil)
			return
		}
		httputil.Error(w, http.StatusBadGateway, "cancel failed", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PipelineHandler) ApproveJob(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}
	run, err := h.pipelineClient.GetRun(r.Context(), chi.URLParam(r, "runID"))
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "run not found", nil)
			return
		}
		httputil.Error(w, http.StatusBadGateway, "orchestrator unavailable", err)
		return
	}
	projID, ok := projectIDFromRunMap(run)
	if !ok {
		httputil.Error(w, http.StatusBadGateway, "invalid run payload", nil)
		return
	}
	if !requireProjectAccess(r.Context(), w, h.projectsClient, userID, projID) {
		return
	}
	err = h.pipelineClient.ApproveJob(r.Context(), chi.URLParam(r, "runID"), chi.URLParam(r, "jobID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "approve failed", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PipelineHandler) RejectJob(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}
	run, err := h.pipelineClient.GetRun(r.Context(), chi.URLParam(r, "runID"))
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			httputil.Error(w, http.StatusNotFound, "run not found", nil)
			return
		}
		httputil.Error(w, http.StatusBadGateway, "orchestrator unavailable", err)
		return
	}
	projID, ok := projectIDFromRunMap(run)
	if !ok {
		httputil.Error(w, http.StatusBadGateway, "invalid run payload", nil)
		return
	}
	if !requireProjectAccess(r.Context(), w, h.projectsClient, userID, projID) {
		return
	}
	err = h.pipelineClient.RejectJob(r.Context(), chi.URLParam(r, "runID"), chi.URLParam(r, "jobID"))
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "reject failed", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
