package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/vsuaiqq/cicd/shared/logger"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.L().Warn().Err(err).Msg("httputil: encode response failed")
	}
}

func Error(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		logger.L().Warn().Err(err).Msg("httputil: response error")
	}
	JSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
