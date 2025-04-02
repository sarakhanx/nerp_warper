package handler

import (
	"nerp_wrapper/application/service"
	"nerp_wrapper/domain/entity"
	"time"

	"github.com/gofiber/fiber/v2"
)

// InvoiceHandler handles HTTP requests for invoices
type InvoiceHandler struct {
	invoiceService *service.InvoiceService
}

// NewInvoiceHandler creates a new instance of InvoiceHandler
func NewInvoiceHandler(invoiceService *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{invoiceService: invoiceService}
}

// GetAllInvoices handles GET request to retrieve invoices with pagination
func (h *InvoiceHandler) GetAllInvoices(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 500)

	invoices, err := h.invoiceService.GetAllInvoices(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(invoices)
}

// GetDailyInvoiceSummary handles GET request to retrieve daily invoice summary
func (h *InvoiceHandler) GetDailyInvoiceSummary(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 500)

	summary, err := h.invoiceService.GetDailyInvoiceSummary(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}

// GetPeriodInvoiceSummary handles GET request to retrieve period-based invoice summary
func (h *InvoiceHandler) GetPeriodInvoiceSummary(c *fiber.Ctx) error {
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

	summary, err := h.invoiceService.GetPeriodInvoiceSummary(periodType, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}
