package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vsuaiqq/cicd/api-gateway/internal/client"
	"github.com/vsuaiqq/cicd/shared/httputil"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/projects"
)

type ProjectsHandler struct {
	authClient     *client.AuthClient
	projectsClient *client.ProjectsClient
}

func NewProjectsHandler(authClient *client.AuthClient, projectsClient *client.ProjectsClient) *ProjectsHandler {
	return &ProjectsHandler{
		authClient:     authClient,
		projectsClient: projectsClient,
	}
}

func (h *ProjectsHandler) RegisterRoutes(r chi.Router) {
	r.Route("/projects", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Patch("/{id}", h.Update)
		r.Post("/{id}/verify", h.Verify)
		r.Delete("/{id}", h.Delete)
		r.Get("/{id}/env", h.GetEnvVars)
		r.Put("/{id}/env", h.UpdateEnvVars)
		r.Get("/{id}/pipeline-yaml", h.GetPipelineYAML)
		r.Put("/{id}/pipeline-yaml", h.SetPipelineYAML)

		r.Get("/{id}/secrets", h.ListSecrets)
		r.Put("/{id}/secrets/{key}", h.SetSecret)
		r.Delete("/{id}/secrets/{key}", h.DeleteSecret)

		r.Get("/{id}/members", h.ListMembers)
		r.Post("/{id}/members", h.InviteMember)
		r.Patch("/{id}/members/{userId}", h.UpdateMemberRole)
		r.Delete("/{id}/members/{userId}", h.RemoveMember)
	})
}

type projectResponse struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	RepoURL       string `json:"repo_url"`
	DefaultBranch string `json:"default_branch"`
	PublicKey     string `json:"public_key"`
	WebhookSecret string `json:"webhook_secret"`
	WebhookURL    string `json:"webhook_url"`
	Status        string `json:"status"`
}

type projectSummary struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	RepoURL       string `json:"repo_url"`
	DefaultBranch string `json:"default_branch"`
	Status        string `json:"status"`
}

type envVarJSON struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type createProjectRequest struct {
	Name          string `json:"name"`
	RepositoryURL string `json:"repository_url"`
	DefaultBranch string `json:"default_branch"`
}

type updateProjectRequest struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

type updateEnvVarsRequest struct {
	Vars []envVarJSON `json:"vars"`
}

type verifyResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req createProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if req.Name == "" || req.RepositoryURL == "" {
		httputil.Error(w, http.StatusBadRequest, "name and repository_url are required", nil)
		return
	}

	resp, err := h.projectsClient.CreateProject(r.Context(), userID, req.Name, req.RepositoryURL, req.DefaultBranch)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, projectResponse{
		ID:            resp.Id,
		UserID:        resp.OwnerUserId,
		Name:          resp.Name,
		RepoURL:       resp.RepoUrl,
		DefaultBranch: resp.DefaultBranch,
		PublicKey:     resp.PublicKey,
		WebhookSecret: resp.WebhookSecret,
		WebhookURL:    resp.WebhookUrl,
		Status:        resp.Status,
	})
}

func (h *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	resp, err := h.projectsClient.ListProjects(r.Context(), userID)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	list := make([]projectSummary, 0, len(resp.Projects))
	for _, p := range resp.Projects {
		list = append(list, projectSummary{
			ID:            p.Id,
			Name:          p.Name,
			RepoURL:       p.RepoUrl,
			DefaultBranch: p.DefaultBranch,
			Status:        p.Status,
		})
	}
	httputil.JSON(w, http.StatusOK, list)
}

func (h *ProjectsHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	resp, err := h.projectsClient.GetProject(r.Context(), userID, chi.URLParam(r, "id"))
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, projectResponse{
		ID:            resp.Id,
		UserID:        resp.OwnerUserId,
		Name:          resp.Name,
		RepoURL:       resp.RepoUrl,
		DefaultBranch: resp.DefaultBranch,
		PublicKey:     resp.PublicKey,
		WebhookSecret: resp.WebhookSecret,
		WebhookURL:    resp.WebhookUrl,
		Status:        resp.Status,
	})
}

func (h *ProjectsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req updateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if req.Name == "" {
		httputil.Error(w, http.StatusBadRequest, "name is required", nil)
		return
	}

	resp, err := h.projectsClient.UpdateProject(r.Context(), userID, chi.URLParam(r, "id"), req.Name, req.DefaultBranch)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, projectResponse{
		ID:            resp.Id,
		UserID:        resp.OwnerUserId,
		Name:          resp.Name,
		RepoURL:       resp.RepoUrl,
		DefaultBranch: resp.DefaultBranch,
		PublicKey:     resp.PublicKey,
		WebhookSecret: resp.WebhookSecret,
		WebhookURL:    resp.WebhookUrl,
		Status:        resp.Status,
	})
}

func (h *ProjectsHandler) Verify(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	resp, err := h.projectsClient.VerifyConnection(r.Context(), userID, chi.URLParam(r, "id"))
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, verifyResponse{
		Success: resp.Success,
		Status:  resp.Status,
		Error:   resp.Error,
	})
}

func (h *ProjectsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	_, err := h.projectsClient.DeleteProject(r.Context(), userID, chi.URLParam(r, "id"))
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectsHandler) GetEnvVars(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	resp, err := h.projectsClient.GetEnvVars(r.Context(), userID, chi.URLParam(r, "id"))
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	vars := make([]envVarJSON, 0, len(resp.Vars))
	for _, v := range resp.Vars {
		vars = append(vars, envVarJSON{Key: v.Key, Value: v.Value})
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"vars": vars})
}

func (h *ProjectsHandler) UpdateEnvVars(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req updateEnvVarsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	pbVars := make([]*pb.EnvVar, 0, len(req.Vars))
	seen := make(map[string]struct{}, len(req.Vars))
	for _, v := range req.Vars {
		if v.Key == "" {
			continue
		}
		if _, dup := seen[v.Key]; dup {
			httputil.Error(w, http.StatusBadRequest, "duplicate env var key: "+v.Key, nil)
			return
		}
		seen[v.Key] = struct{}{}
		pbVars = append(pbVars, &pb.EnvVar{Key: v.Key, Value: v.Value})
	}

	_, err := h.projectsClient.UpdateEnvVars(r.Context(), userID, chi.URLParam(r, "id"), pbVars)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectsHandler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	resp, err := h.projectsClient.ListSecrets(r.Context(), userID, chi.URLParam(r, "id"))
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	type secretItem struct {
		Key       string `json:"key"`
		UpdatedAt string `json:"updated_at"`
	}
	items := make([]secretItem, len(resp.Secrets))
	for i, s := range resp.Secrets {
		items[i] = secretItem{Key: s.Key, UpdatedAt: s.UpdatedAt}
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"secrets": items})
}

func (h *ProjectsHandler) SetSecret(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	key := chi.URLParam(r, "key")
	if key == "" {
		httputil.Error(w, http.StatusBadRequest, "key required", nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if req.Value == "" {
		httputil.Error(w, http.StatusBadRequest, "value must not be empty", nil)
		return
	}

	_, err := h.projectsClient.SetSecret(r.Context(), userID, chi.URLParam(r, "id"), key, req.Value)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectsHandler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	key := chi.URLParam(r, "key")
	if key == "" {
		httputil.Error(w, http.StatusBadRequest, "key required", nil)
		return
	}

	_, err := h.projectsClient.DeleteSecret(r.Context(), userID, chi.URLParam(r, "id"), key)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectsHandler) GetPipelineYAML(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	resp, err := h.projectsClient.GetPipelineYAML(r.Context(), userID, chi.URLParam(r, "id"))
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]string{"yaml": resp.Yaml})
}

func (h *ProjectsHandler) SetPipelineYAML(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req struct {
		YAML string `json:"yaml"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	_, err := h.projectsClient.SetPipelineYAML(r.Context(), userID, chi.URLParam(r, "id"), req.YAML)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type memberResponse struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	InvitedBy   string `json:"invited_by"`
	CreatedAt   string `json:"created_at"`
}

func (h *ProjectsHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	resp, err := h.projectsClient.ListMembers(r.Context(), userID, chi.URLParam(r, "id"))
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	members := make([]memberResponse, 0, len(resp.Members))
	for _, m := range resp.Members {
		members = append(members, memberResponse{
			UserID:      m.UserId,
			Email:       m.Email,
			DisplayName: m.DisplayName,
			Role:        m.Role,
			InvitedBy:   m.InvitedBy,
			CreatedAt:   m.CreatedAt,
		})
	}
	httputil.JSON(w, http.StatusOK, map[string]any{
		"members":        members,
		"requester_role": resp.RequesterRole,
		"owner_user_id":  resp.OwnerUserId,
	})
}

func (h *ProjectsHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if req.Email == "" {
		httputil.Error(w, http.StatusBadRequest, "email is required", nil)
		return
	}
	if req.Role != "editor" && req.Role != "viewer" {
		req.Role = "viewer"
	}

	authResp, err := h.authClient.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	resp, err := h.projectsClient.InviteMember(
		r.Context(),
		userID,
		chi.URLParam(r, "id"),
		authResp.UserId,
		authResp.Email,
		authResp.DisplayName,
		req.Role,
	)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	m := resp.Member
	httputil.JSON(w, http.StatusCreated, memberResponse{
		UserID:      m.UserId,
		Email:       m.Email,
		DisplayName: m.DisplayName,
		Role:        m.Role,
		InvitedBy:   m.InvitedBy,
		CreatedAt:   m.CreatedAt,
	})
}

func (h *ProjectsHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	resp, err := h.projectsClient.UpdateMemberRole(
		r.Context(),
		userID,
		chi.URLParam(r, "id"),
		chi.URLParam(r, "userId"),
		req.Role,
	)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	m := resp.Member
	httputil.JSON(w, http.StatusOK, memberResponse{
		UserID:      m.UserId,
		Email:       m.Email,
		DisplayName: m.DisplayName,
		Role:        m.Role,
		InvitedBy:   m.InvitedBy,
		CreatedAt:   m.CreatedAt,
	})
}

func (h *ProjectsHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID := requireUserID(w, r, h.authClient)
	if userID == "" {
		return
	}

	_, err := h.projectsClient.RemoveMember(
		r.Context(),
		userID,
		chi.URLParam(r, "id"),
		chi.URLParam(r, "userId"),
	)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
