package repository

import (
	"auth-service/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "repository.SaveUser"

	// SQL-запрос для PostgreSQL с RETURNING
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query, email, passHash).Scan(&id)
	if err != nil {
		// Проверка на уникальность email (PostgreSQL unique violation)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // unique_violation
			return 0, fmt.Errorf("%s:%w", op, ErrUserExists)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}

func (r *UserRepository) GetUser(ctx context.Context, email string) (model.User, error) {
	const op = "repository.GetUser"

	var user model.User
	// SQL-запрос для PostgreSQL
	query := `SELECT id, email, pass_hash, is_admin, created_at, updated_at
	          FROM users
	          WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PassHash,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserRepository) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "repository.IsAdmin"

	var isAdmin bool
	// SQL-запрос для PostgreSQL
	query := `SELECT is_admin FROM users WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("%s: user not found: %w", op, ErrAppNotFound)
		}
		return false, fmt.Errorf("%s: query error: %w", op, err)
	}

	return isAdmin, nil
}

func (r *UserRepository) App(ctx context.Context, appID int) (model.App, error) {
	const op = "repository.App"

	var app model.App
	query := `SELECT id, name, secret FROM apps WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, appID).Scan(
		&app.ID,
		&app.Name,
		&app.Secret,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.App{}, fmt.Errorf("%s: app not found: %w", op, ErrAppNotFound)
		}
		return model.App{}, fmt.Errorf("%s: query error: %w", op, err)
	}

	return app, nil
}
