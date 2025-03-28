package router

import (
	"nerp_wrapper/interfaces/http/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRouter sets up all the routes for the application
func SetupRouter(app *fiber.App, authHandler *handler.AuthHandler, saleHandler *handler.SaleHandler) {
	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/user-info/:id", authHandler.GetUserInfo)

	// Sales routes
	sales := app.Group("/sales")
	sales.Get("/", saleHandler.GetAllSaleOrders)
}
