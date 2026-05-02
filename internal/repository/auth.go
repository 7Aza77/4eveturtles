package repository

import (
	"context"
	"goevent/internal/entity"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user entity.User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByID(ctx context.Context, id int64) (entity.User, error)
}

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user entity.User) (int64, error) {
	var id int64
	query := "INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Role).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	query := "SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1"
	err := r.db.GetContext(ctx, &user, query, email)
	return user, err
}

func (r *AuthPostgres) GetUserByID(ctx context.Context, id int64) (entity.User, error) {
	var user entity.User
	query := "SELECT id, email, role, created_at, updated_at FROM users WHERE id = $1"
	err := r.db.GetContext(ctx, &user, query, id)
	return user, err
}
