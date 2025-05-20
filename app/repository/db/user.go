package db

import (
	"context"
	"database/sql"
	"log/slog"
	"user-service/app/domain"
)

type userRepository struct {
	conn *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, phone, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.conn.QueryRowContext(ctx, query, user.Email, user.Phone, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "[userRepository] Create", "exec", err)
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, phone, password, shop_id, created_at, updated_at FROM users WHERE email = $1`
	row := r.conn.QueryRowContext(ctx, query, email)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Phone, &user.Password, &user.ShopID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "[userRepository] GetByEmail", "scan", err)
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	query := `SELECT id, email, phone, password, shop_id, created_at, updated_at FROM users WHERE phone = $1`
	row := r.conn.QueryRowContext(ctx, query, phone)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Phone, &user.Password, &user.ShopID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "[userRepository] GetByPhone", "scan", err)
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}
