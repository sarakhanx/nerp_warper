package entity

import "time"

// PeriodType represents the type of period for sales summary
type PeriodType string

const (
	PeriodTypeDay     PeriodType = "1D"      // Last 24 hours
	PeriodTypeWeek    PeriodType = "7D"      // Last 7 days
	PeriodTypeMonth   PeriodType = "30D"     // Last 30 days
	PeriodTypeQuarter PeriodType = "90D"     // Last 90 days
	PeriodTypeMonthly PeriodType = "MONTHLY" // Group by months
	PeriodTypeYearly  PeriodType = "YEARLY"  // Group by years
)

// DateRange represents a date range for filtering
type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

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

// PeriodSalesSummaryResponse represents the response for period-based sales summary
type PeriodSalesSummaryResponse struct {
	Period       string              `json:"period"`
	PeriodType   PeriodType          `json:"period_type"`
	DateRange    DateRange           `json:"date_range"`
	Items        []DailySalesSummary `json:"items"`
	TotalAmount  float64             `json:"total_amount"`
	OrderCount   int                 `json:"order_count"`
	AverageDaily float64             `json:"average_daily"`
}
