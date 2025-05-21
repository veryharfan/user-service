package usecase

import (
	"context"
	"log/slog"
	"strings"
	"user-service/app/domain"
	"user-service/config"
	"user-service/pkg"
)

type userUsecase struct {
	userRepo domain.UserRepository
	cfg      *config.Config
}

func NewUserUsecase(userRepo domain.UserRepository, cfg *config.Config) domain.UserUsecase {
	return &userUsecase{userRepo, cfg}
}

func (u *userUsecase) Register(ctx context.Context, req *domain.UserRequest) (*domain.UserResponse, error) {
	user := &domain.User{}
	switch pkg.GetUsernameType(req.Username) {
	case "email":
		user.Email = &req.Username
	case "phone":
		req.Username = strings.Trim(req.Username, "+")
		user.Phone = &req.Username
	default:
		slog.ErrorContext(ctx, "[userUsecase] Register", "invalidUsernameType", req.Username)
		return nil, domain.ErrValidation
	}

	var err error
	user.Password, err = pkg.HashPassword(req.Password)
	if err != nil {
		slog.ErrorContext(ctx, "[userUsecase] Register", "hashPassword", err)
		return nil, err
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		slog.ErrorContext(ctx, "[userUsecase] Register", "createUser", err)
		return nil, err
	}

	return &domain.UserResponse{
		Username: req.Username,
	}, nil
}

func (u *userUsecase) Login(ctx context.Context, req *domain.UserRequest) (*domain.LoginResponse, error) {
	var user *domain.User
	var err error

	switch pkg.GetUsernameType(req.Username) {
	case "email":
		user, err = u.userRepo.GetByEmail(ctx, req.Username)
	case "phone":
		user, err = u.userRepo.GetByPhone(ctx, strings.Trim(req.Username, "+"))
	default:
		slog.ErrorContext(ctx, "[userUsecase] Login", "invalidUsernameType", req.Username)
		return nil, domain.ErrValidation
	}

	if err != nil {
		slog.ErrorContext(ctx, "[userUsecase] Login", "getUser", err)
		return nil, err
	}

	if !pkg.CheckPasswordHash(req.Password, user.Password) {
		return nil, domain.ErrUnauthorized
	}

	token, err := pkg.CreateJwtToken(user.ID, user.ShopID, u.cfg.Jwt.SecretKey, u.cfg.Jwt.Expire)
	if err != nil {
		slog.ErrorContext(ctx, "[userUsecase] Login", "createJwtToken", err)
		return nil, err
	}

	slog.InfoContext(ctx, "[userUsecase] Login", req.Username, "User logged in successfully")

	return &domain.LoginResponse{
		Token: token,
	}, nil
}

func (u *userUsecase) AddShopID(ctx context.Context, userID int64, shopID int64) error {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "[userUsecase] AddShopID", "getUser", err)
		return err
	}

	err = u.userRepo.AddShopID(ctx, user.ID, shopID)
	if err != nil {
		slog.ErrorContext(ctx, "[userUsecase] AddShopID", "addShopID", err)
		return err
	}

	slog.InfoContext(ctx, "[userUsecase] AddShopID", "success", "Shop ID added successfully")
	return nil
}
