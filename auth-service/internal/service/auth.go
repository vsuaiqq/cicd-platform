package service

import (
	"context"
	"crypto/subtle"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/vsuaiqq/cicd/auth-service/internal/models"
	"github.com/vsuaiqq/cicd/auth-service/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	jwtService       *JWTService
}

func NewAuthService(userRepo *repository.UserRepository, refreshTokenRepo *repository.RefreshTokenRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
	}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (*models.User, error) {
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, ErrUserExists
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.NewUser(email, username, string(hashedPassword))
	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, ErrUserExists
		}
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, TokenPair, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, TokenPair{}, ErrInvalidCredentials
		}
		return nil, TokenPair{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, TokenPair{}, ErrInvalidCredentials
	}

	tokens, err := s.generateAndSaveTokens(ctx, user)
	if err != nil {
		return nil, TokenPair{}, err
	}

	return user, tokens, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*models.User, *Claims, error) {
	claims, err := s.jwtService.ValidateAccessToken(token)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}

	return user, claims, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (TokenPair, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return TokenPair{}, err
	}

	storedToken, err := s.refreshTokenRepo.GetRefreshToken(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return TokenPair{}, ErrInvalidCredentials
		}
		return TokenPair{}, err
	}

	if subtle.ConstantTimeCompare([]byte(storedToken), []byte(refreshToken)) != 1 {
		return TokenPair{}, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return TokenPair{}, ErrInvalidCredentials
		}
		return TokenPair{}, err
	}

	tokens, err := s.generateAndSaveTokens(ctx, user)
	if err != nil {
		return TokenPair{}, err
	}

	return tokens, nil
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}
	return user, err
}

func (s *AuthService) generateAndSaveTokens(ctx context.Context, user *models.User) (TokenPair, error) {
	accessToken, accessExpiresAt, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, refreshExpiresAt, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return TokenPair{}, err
	}

	if err := s.refreshTokenRepo.SaveRefreshToken(ctx, user.ID, refreshToken, refreshExpiresAt); err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}
