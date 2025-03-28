package handler

import (
	"nerp_wrapper/application/service"
	"nerp_wrapper/domain/entity"

	"time"

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

// GetDailySalesSummary handles GET request to retrieve daily sales summary
func (h *SaleHandler) GetDailySalesSummary(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 500)

	summary, err := h.saleService.GetDailySalesSummary(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}

// GetPeriodSalesSummary handles GET request to retrieve period-based sales summary
func (h *SaleHandler) GetPeriodSalesSummary(c *fiber.Ctx) error {
	// Get period type from query param, default to 30D
	periodType := entity.PeriodType(c.Query("period_type", string(entity.PeriodTypeMonth)))

	// Parse custom date range if provided
	var startDate, endDate *time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		if date, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = &date
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if date, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = &date
		}
	}

	// Both start and end dates must be provided if using custom range
	if (startDate != nil && endDate == nil) || (startDate == nil && endDate != nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both start_date and end_date must be provided for custom date range",
		})
	}

	// Validate period type
	switch periodType {
	case entity.PeriodTypeDay, entity.PeriodTypeWeek, entity.PeriodTypeMonth,
		entity.PeriodTypeQuarter, entity.PeriodTypeMonthly, entity.PeriodTypeYearly:
		// Valid period types
	default:
		if startDate == nil || endDate == nil {
			// Invalid period type and no custom range
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid period_type. Must be one of: 1D, 7D, 30D, 90D, MONTHLY, YEARLY",
			})
		}
	}

	summary, err := h.saleService.GetPeriodSalesSummary(periodType, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}
