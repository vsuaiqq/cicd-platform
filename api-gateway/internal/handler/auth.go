package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vsuaiqq/cicd/api-gateway/internal/client"
	"github.com/vsuaiqq/cicd/shared/httputil"
)

type AuthHandler struct {
	authClient *client.AuthClient
}

func NewAuthHandler(authClient *client.AuthClient) *AuthHandler {
	return &AuthHandler{authClient: authClient}
}

func (h *AuthHandler) Healthcheck(ctx context.Context) error {
	_, err := h.authClient.Health(ctx)
	return err
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/validate", h.ValidateToken)
		r.Post("/refresh", h.RefreshToken)
	})
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if req.Email == "" || req.Password == "" {
		httputil.Error(w, http.StatusBadRequest, "email and password are required", nil)
		return
	}

	resp, err := h.authClient.Register(r.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, RegisterResponse{
		Success: resp.Success,
		Message: resp.Message,
		UserID:  resp.UserId,
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	UserID       string `json:"user_id"`
	TokenType    string `json:"token_type"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if req.Email == "" || req.Password == "" {
		httputil.Error(w, http.StatusBadRequest, "email and password are required", nil)
		return
	}

	resp, err := h.authClient.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		UserID:       resp.UserId,
		TokenType:    "Bearer",
	})
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid     bool   `json:"valid"`
	UserID    string `json:"user_id,omitempty"`
	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req ValidateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if req.Token == "" {
		httputil.Error(w, http.StatusBadRequest, "token is required", nil)
		return
	}

	resp, err := h.authClient.ValidateToken(r.Context(), req.Token)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, ValidateTokenResponse{
		Valid:     resp.Valid,
		UserID:    resp.UserId,
		Email:     resp.Email,
		Username:  resp.Username,
		ExpiresAt: resp.ExpiresAt,
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if req.RefreshToken == "" {
		httputil.Error(w, http.StatusBadRequest, "refresh_token is required", nil)
		return
	}

	resp, err := h.authClient.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		grpcToHTTPError(w, err)
		return
	}

	httputil.JSON(w, http.StatusOK, RefreshTokenResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		TokenType:    "Bearer",
	})
}
