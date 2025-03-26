package router

import (
	"nerp_wrapper/interfaces/http/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRouter configures the application routes
func SetupRouter(app *fiber.App, authHandler *handler.AuthHandler) {
	// API group
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/user/:id", authHandler.GetUserInfo)
}
