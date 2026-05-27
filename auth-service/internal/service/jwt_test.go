package service

import (
	"testing"
	"time"

	"github.com/vsuaiqq/cicd/auth-service/internal/config"
)

func TestJWTService_accessTokenRoundTrip(t *testing.T) {
	svc := NewJWTService(config.JWTConfig{
		AccessTokenSecret:  "access-secret-min-32-characters-long",
		RefreshTokenSecret: "refresh-secret-min-32-characters-long",
		AccessTokenTTL:     3600,
		RefreshTokenTTL:    86400,
	})

	token, expiresAt, err := svc.GenerateAccessToken("user-1", "demo@example.com")
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}
	if expiresAt.Before(time.Now()) {
		t.Fatal("expiresAt in the past")
	}

	claims, err := svc.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken: %v", err)
	}
	if claims.UserID != "user-1" || claims.Email != "demo@example.com" {
		t.Fatalf("claims = %+v", claims)
	}
}

func TestJWTService_refreshTokenRoundTrip(t *testing.T) {
	svc := NewJWTService(config.JWTConfig{
		AccessTokenSecret:  "access-secret-min-32-characters-long",
		RefreshTokenSecret: "refresh-secret-min-32-characters-long",
		AccessTokenTTL:     3600,
		RefreshTokenTTL:    86400,
	})

	token, _, err := svc.GenerateRefreshToken("user-42")
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}

	claims, err := svc.ValidateRefreshToken(token)
	if err != nil {
		t.Fatalf("ValidateRefreshToken: %v", err)
	}
	if claims.UserID != "user-42" {
		t.Fatalf("user_id = %q", claims.UserID)
	}
}

func TestJWTService_rejectsWrongSecret(t *testing.T) {
	svc := NewJWTService(config.JWTConfig{
		AccessTokenSecret:  "access-secret-min-32-characters-long",
		RefreshTokenSecret: "refresh-secret-min-32-characters-long",
		AccessTokenTTL:     3600,
		RefreshTokenTTL:    86400,
	})
	other := NewJWTService(config.JWTConfig{
		AccessTokenSecret:  "other-access-secret-min-32-characters",
		RefreshTokenSecret: "other-refresh-secret-min-32-characters",
		AccessTokenTTL:     3600,
		RefreshTokenTTL:    86400,
	})

	token, _, err := svc.GenerateAccessToken("user-1", "demo@example.com")
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}
	if _, err := other.ValidateAccessToken(token); err == nil {
		t.Fatal("expected validation error with wrong secret")
	}
}
