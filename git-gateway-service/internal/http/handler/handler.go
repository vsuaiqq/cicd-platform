package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vsuaiqq/cicd/git-gateway-service/internal/kafka"
	"github.com/vsuaiqq/cicd/git-gateway-service/internal/projects"
	"github.com/vsuaiqq/cicd/shared/events"
	"github.com/vsuaiqq/cicd/shared/httputil"
	"github.com/vsuaiqq/cicd/shared/logger"
)

type Processor interface {
	ValidateAndParse(r *http.Request, secretKey []byte) (*events.GitEvent, error)
}

const maxWebhookBodyBytes = 1 << 20

func WebhookHandler(
	projectsClient *projects.Client,
	publisher *kafka.GitEventPublisher,
	proc Processor,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxWebhookBodyBytes)

		projectID := chi.URLParam(r, "projectID")

		proj, err := projectsClient.GetProjectInternal(r.Context(), projectID)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				httputil.Error(w, http.StatusNotFound, "project not found", nil)
				return
			}
			logger.L().Error().Err(err).Str("project_id", projectID).Msg("GetProjectInternal failed")
			httputil.Error(w, http.StatusInternalServerError, "failed to look up project", nil)
			return
		}

		if proj.WebhookSecret == "" {
			logger.L().Warn().Str("project_id", projectID).Msg("project has no webhook secret configured")
			httputil.Error(w, http.StatusForbidden, "webhook not configured for this project", nil)
			return
		}

		gitEvent, err := proc.ValidateAndParse(r, []byte(proj.WebhookSecret))
		if err != nil {
			logger.L().Warn().Err(err).Str("project_id", projectID).Msg("webhook processing failed")
			httputil.Error(w, http.StatusBadRequest, "invalid webhook request", nil)
			return
		}

		if gitEvent == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		gitEvent.ProjectID = projectID

		if err := publisher.Publish(r.Context(), gitEvent); err != nil {
			logger.L().Error().Err(err).Str("project_id", projectID).Msg("failed to publish event")
			httputil.Error(w, http.StatusInternalServerError, "failed to enqueue event", nil)
			return
		}

		httputil.JSON(w, http.StatusAccepted, nil)
	}
}
