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

// GetAllSaleOrders handles GET request to retrieve sale orders with pagination
func (h *SaleHandler) GetAllSaleOrders(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 500)

	orders, err := h.saleService.GetAllSaleOrders(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(orders)
}
