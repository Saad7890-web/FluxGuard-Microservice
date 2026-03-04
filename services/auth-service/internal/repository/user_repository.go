package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (string, error) {
	id := uuid.New().String()

	query := `
	INSERT INTO users (id, email, password_hash, role)
	VALUES ($1, $2, $3, 'user')`

	_, err := r.db.Exec(ctx, query, id, email, passwordHash)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (string, string, string, error) {
	var id, passwordHash, role string

	query := `
	SELECT id, password_hash, role
	FROM users WHERE email=$1`

	err := r.db.QueryRow(ctx, query, email).Scan(&id, &passwordHash, &role)
	if err != nil {
		return "", "", "", err
	}

	return id, passwordHash, role, nil
}

func (r *UserRepository) StoreRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt int64) error {
	query := `
	INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
	VALUES ($1, $2, to_timestamp($3))`

	_, err := r.db.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}