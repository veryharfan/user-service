package domain

import (
	"context"
	"time"
)

type UserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type User struct {
	ID        int64
	Email     *string
	Phone     *string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
}

type UserUsecase interface {
	Register(ctx context.Context, req *UserRequest) (*UserResponse, error)
	Login(ctx context.Context, req *UserRequest) (*LoginResponse, error)
}
