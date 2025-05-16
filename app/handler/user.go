package handler

import (
	"log/slog"
	"user-service/app/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userUsecase domain.UserUsecase
	validator   *validator.Validate
}

func NewUserHandler(userUsecase domain.UserUsecase, validator *validator.Validate) *userHandler {
	return &userHandler{userUsecase, validator}
}

func (h *userHandler) Register(c *fiber.Ctx) error {
	var req domain.UserRequest
	if err := c.BodyParser(&req); err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Register", "bodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(Error(domain.ErrBadRequest))
	}

	if err := h.validator.Struct(req); err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Register", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(Error(domain.ErrValidation))
	}

	user, err := h.userUsecase.Register(c.Context(), &req)
	if err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Register", "usecase", err)
		status, response := FromError(err)
		return c.Status(status).JSON(response)
	}

	return c.Status(fiber.StatusCreated).JSON(Success(user))
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	var req domain.UserRequest
	if err := c.BodyParser(&req); err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Login", "bodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(Error(domain.ErrBadRequest))
	}

	if err := h.validator.Struct(req); err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Login", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(Error(domain.ErrValidation))
	}

	login, err := h.userUsecase.Login(c.Context(), &req)
	if err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Login", "usecase", err)
		status, response := FromError(err)
		return c.Status(status).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(Success(login))
}
