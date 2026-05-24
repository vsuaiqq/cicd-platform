package handler

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vsuaiqq/cicd/auth-service/internal/service"
	"github.com/vsuaiqq/cicd/shared/logger"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/auth"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	username := req.Username
	if username == "" {
		username = req.Email
	}

	user, err := h.authService.Register(ctx, req.Email, username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		logger.L().Error().Err(err).Msg("register error")
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "user registered successfully",
		UserId:  user.ID,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	user, tokens, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		logger.L().Error().Err(err).Msg("login error")
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &pb.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int64(time.Until(tokens.AccessExpiresAt).Seconds()),
		UserId:       user.ID,
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	user, claims, err := h.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrExpiredToken) ||
			errors.Is(err, service.ErrInvalidCredentials) {
			return &pb.ValidateTokenResponse{Valid: false}, nil
		}
		logger.L().Error().Err(err).Msg("validate token error")
		return nil, status.Error(codes.Internal, "failed to validate token")
	}

	var expiresAt int64
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Unix()
	}

	return &pb.ValidateTokenResponse{
		Valid:     true,
		UserId:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		ExpiresAt: expiresAt,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	tokens, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrExpiredToken) ||
			errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired refresh token")
		}
		logger.L().Error().Err(err).Msg("refresh token error")
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int64(time.Until(tokens.AccessExpiresAt).Seconds()),
	}, nil
}

func (h *AuthHandler) Health(_ context.Context, _ *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Healthy: true, Status: "ok"}, nil
}

func (h *AuthHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	user, err := h.authService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		logger.L().Error().Err(err).Msg("get user by email error")
		return nil, status.Error(codes.Internal, "failed to look up user")
	}

	return &pb.GetUserByEmailResponse{
		UserId:      user.ID,
		Email:       user.Email,
		DisplayName: user.Username,
	}, nil
}
