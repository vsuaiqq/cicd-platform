package handler

import (
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vsuaiqq/cicd/api-gateway/internal/client"
	"github.com/vsuaiqq/cicd/shared/httputil"
)

func requireProjectAccess(ctx context.Context, w http.ResponseWriter, projects *client.ProjectsClient, userID, projectID string) bool {
	if projectID == "" {
		httputil.Error(w, http.StatusBadRequest, "project_id is required", nil)
		return false
	}
	_, err := projects.GetProject(ctx, userID, projectID)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			httputil.Error(w, http.StatusNotFound, "project not found", nil)
			return false
		}
		grpcToHTTPError(w, err)
		return false
	}
	return true
}

func projectIDFromRunMap(m map[string]any) (string, bool) {
	v, ok := m["project_id"]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok && s != ""
}
