package handler

import (
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vsuaiqq/cicd/api-gateway/internal/client"
	"github.com/vsuaiqq/cicd/shared/httputil"
)

const maxRequestBodyBytes = 1 << 20

func requireUserID(w http.ResponseWriter, r *http.Request, authClient *client.AuthClient) string {
	token, ok := bearerToken(r)
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "missing or invalid Authorization header", nil)
		return ""
	}
	resp, err := authClient.ValidateToken(r.Context(), token)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unauthenticated {
			httputil.Error(w, http.StatusUnauthorized, st.Message(), nil)
			return ""
		}
		httputil.Error(w, http.StatusInternalServerError, "token validation failed", err)
		return ""
	}
	if !resp.Valid || resp.UserId == "" {
		httputil.Error(w, http.StatusUnauthorized, "invalid token", nil)
		return ""
	}
	return resp.UserId
}

func bearerToken(r *http.Request) (string, bool) {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	return token, ok && token != ""
}

func grpcToHTTPError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		httputil.Error(w, http.StatusInternalServerError, "internal server error", err)
		return
	}
	switch st.Code() {
	case codes.InvalidArgument:
		httputil.Error(w, http.StatusBadRequest, st.Message(), nil)
	case codes.Unauthenticated:
		httputil.Error(w, http.StatusUnauthorized, st.Message(), nil)
	case codes.PermissionDenied:
		httputil.Error(w, http.StatusForbidden, st.Message(), nil)
	case codes.NotFound:
		httputil.Error(w, http.StatusNotFound, st.Message(), nil)
	case codes.AlreadyExists:
		httputil.Error(w, http.StatusConflict, st.Message(), nil)
	case codes.ResourceExhausted:
		httputil.Error(w, http.StatusTooManyRequests, st.Message(), nil)
	case codes.DeadlineExceeded:
		httputil.Error(w, http.StatusRequestTimeout, "request timeout", nil)
	case codes.Unavailable:
		httputil.Error(w, http.StatusServiceUnavailable, "service unavailable", nil)
	default:
		httputil.Error(w, http.StatusInternalServerError, "internal server error", err)
	}
}
