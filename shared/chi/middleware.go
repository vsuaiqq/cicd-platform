package chi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func DefaultStack(timeout time.Duration) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(timeout),
	}
}
