package handler

import (
	"log/slog"
	"strconv"
	"user-service/app/domain"
	"user-service/app/handler/response"

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
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	if err := h.validator.Struct(req); err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Register", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrValidation))
	}

	user, err := h.userUsecase.Register(c.Context(), &req)
	if err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Register", "usecase", err)
		status, resp := response.FromError(err)
		return c.Status(status).JSON(resp)
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(user))
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	var req domain.UserRequest
	if err := c.BodyParser(&req); err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Login", "bodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	if err := h.validator.Struct(req); err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Login", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrValidation))
	}

	login, err := h.userUsecase.Login(c.Context(), &req)
	if err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] Login", "usecase", err)
		status, resp := response.FromError(err)
		return c.Status(status).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(login))
}

func (h *userHandler) AddShopID(c *fiber.Ctx) error {
	var req domain.AddShopIDRequest
	if err := c.BodyParser(&req); err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] AddShopID", "bodyParser", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	if err := h.validator.Struct(req); err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] AddShopID", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrValidation))
	}

	idstr := c.Params("id")
	if idstr == "" {
		slog.ErrorContext(c.Context(), "[userHandler] AddShopID", "params", "user_id not found")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil || id <= 0 {
		slog.ErrorContext(c.Context(), "[productReadHandler] GetByID", "params:"+idstr, err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	err = h.userUsecase.AddShopID(c.Context(), id, req.ShopID)
	if err != nil {
		slog.ErrorContext(c.Context(), "[userHandler] AddShopID", "usecase", err)
		status, resp := response.FromError(err)
		return c.Status(status).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil))
}
