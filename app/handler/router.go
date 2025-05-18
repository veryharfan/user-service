package handler

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App, handler *userHandler) {
	// Setup user routes
	userGroup := app.Group("/user-service")

	userGroup.Post("/users", handler.Register)
	userGroup.Post("/login", handler.Login)
}
