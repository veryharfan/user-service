package handler

import (
	"user-service/app/middleware"
	"user-service/config"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App, handler *userHandler, cfg *config.Config) {
	// Setup user routes
	userGroup := app.Group("/user-service")

	userGroup.Post("/users", handler.Register)
	userGroup.Post("/login", handler.Login)

	authInternal := app.Group("/internal/user-service").Use(middleware.AuthInternal(cfg))
	authInternal.Patch("/users/:id/shop", handler.AddShopID)
}
