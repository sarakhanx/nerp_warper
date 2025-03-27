package handler

import (
	"nerp_wrapper/application/service"

	"github.com/gofiber/fiber/v2"
)

// SaleHandler handles HTTP requests for sale orders
type SaleHandler struct {
	saleService *service.SaleService
}

// NewSaleHandler creates a new instance of SaleHandler
func NewSaleHandler(saleService *service.SaleService) *SaleHandler {
	return &SaleHandler{saleService: saleService}
}

// GetAllSaleOrders handles GET request to retrieve all sale orders
func (h *SaleHandler) GetAllSaleOrders(c *fiber.Ctx) error {
	orders, err := h.saleService.GetAllSaleOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": orders,
	})
}
