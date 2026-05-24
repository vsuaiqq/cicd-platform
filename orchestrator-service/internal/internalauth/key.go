package internalauth

import (
	"crypto/subtle"
	"net/http"
)

const APIKeyHeader = "X-Internal-API-Key"

func APIKeyMiddleware(expectedKey string) func(http.Handler) http.Handler {
	if expectedKey == "" {
		return func(next http.Handler) http.Handler { return next }
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Header.Get(APIKeyHeader)
			if !secureStringEqual(got, expectedKey) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func secureStringEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
