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
	Username string `json:"username"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AddShopIDRequest struct {
	ShopID int64 `json:"shop_id" validate:"required"`
}

type User struct {
	ID        int64
	Email     *string
	Phone     *string
	Password  string
	ShopID    *int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	AddShopID(ctx context.Context, userID int64, shopID int64) error
	GetByID(ctx context.Context, id int64) (*User, error)
}

type UserUsecase interface {
	Register(ctx context.Context, req *UserRequest) (*UserResponse, error)
	Login(ctx context.Context, req *UserRequest) (*LoginResponse, error)
	AddShopID(ctx context.Context, userID int64, shopID int64) error
}
