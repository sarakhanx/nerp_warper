package entity

import "time"

// SaleOrderSummary represents a single sale order in the summary
type SaleOrderSummary struct {
	OrderNumber int       `json:"order_number"`
	OrderID     int64     `json:"order_id"`
	OrderName   string    `json:"order_name"`
	AmountTotal float64   `json:"amount_total"`
	DateOrder   time.Time `json:"date_order"`
}

// DailySalesSummary represents sales summary for a specific date
type DailySalesSummary struct {
	Date        time.Time          `json:"date"`
	TotalAmount float64            `json:"total_amount"`
	OrderCount  int                `json:"order_count"`
	Orders      []SaleOrderSummary `json:"orders"`
}

// SalesSummaryResponse represents the paginated response for sales summary
type SalesSummaryResponse struct {
	Items      []DailySalesSummary `json:"items"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalItems int                 `json:"total_items"`
	TotalPages int                 `json:"total_pages"`
}
