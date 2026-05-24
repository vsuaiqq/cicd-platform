package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vsuaiqq/cicd/auth-service/internal/models"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

const userColumns = `id, email, username, password_hash, created_at, updated_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanUser(s rowScanner) (*models.User, error) {
	u := &models.User{}
	err := s.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	const query = `
		INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (email) DO NOTHING
		RETURNING created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		user.ID, user.Email, user.Username, user.PasswordHash,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserAlreadyExists
	}
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = `SELECT ` + userColumns + ` FROM users WHERE email = $1`
	u, err := scanUser(r.db.QueryRowContext(ctx, query, email))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return u, err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	const query = `SELECT ` + userColumns + ` FROM users WHERE id = $1`
	u, err := scanUser(r.db.QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return u, err
}
