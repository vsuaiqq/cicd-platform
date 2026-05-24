package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type RefreshTokenRepository struct {
	client *redis.Client
}

func NewRefreshTokenRepository(client *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{client: client}
}

func (r *RefreshTokenRepository) SaveRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return errors.New("token already expired")
	}
	return r.client.Set(ctx, r.key(userID), token, ttl).Err()
}

func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	token, err := r.client.Get(ctx, r.key(userID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrRefreshTokenNotFound
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
	return r.client.Del(ctx, r.key(userID)).Err()
}

func (r *RefreshTokenRepository) key(userID string) string {
	return "refresh_token:" + userID
}
