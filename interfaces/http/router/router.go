package router

import (
	"nerp_wrapper/interfaces/http/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRouter sets up all the routes for the application
func SetupRouter(app *fiber.App, authHandler *handler.AuthHandler, saleHandler *handler.SaleHandler, invoiceHandler *handler.InvoiceHandler) {
	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/user-info/:id", authHandler.GetUserInfo)

	// Sales routes
	sales := app.Group("/sales")
	sales.Get("/", saleHandler.GetAllSaleOrders)
	sales.Get("/daily-summary", saleHandler.GetDailySalesSummary)
	sales.Get("/period-summary", saleHandler.GetPeriodSalesSummary)

	// Invoice routes
	invoices := app.Group("/invoices")
	invoices.Get("/", invoiceHandler.GetAllInvoices)
	invoices.Get("/daily-summary", invoiceHandler.GetDailyInvoiceSummary)
	invoices.Get("/period-summary", invoiceHandler.GetPeriodInvoiceSummary)
}
