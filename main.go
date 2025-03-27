package main

import (
	"log"
	"nerp_wrapper/application/service"
	"nerp_wrapper/infrastructure/odoo"
	"nerp_wrapper/interfaces/http/handler"
	"nerp_wrapper/interfaces/http/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize Odoo repository
	authRepo, err := odoo.NewOdooAuthRepository(
		"sarawut.kh@needshopping.co",  // admin username
		"$===(Amplication30)",         // admin password
		"N_ERP",                       // database name
		"https://erp.needshopping.co", // Odoo URL
	)
	if err != nil {
		log.Fatalf("Failed to initialize Odoo repository: %v", err)
	}

	// Initialize repositories
	saleRepo := odoo.NewOdooSaleRepository(authRepo.GetClient())

	// Initialize services
	authService := service.NewAuthService(authRepo)
	saleService := service.NewSaleService(saleRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	saleHandler := handler.NewSaleHandler(saleService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "NERP Wrapper API",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup routes
	router.SetupRouter(app, authHandler, saleHandler)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
