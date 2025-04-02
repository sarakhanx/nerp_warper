package entity

import "time"

// InvoiceSummary represents a single invoice in the summary
type InvoiceSummary struct {
	InvoiceNumber int       `json:"invoice_number"`
	InvoiceID     int64     `json:"invoice_id"`
	InvoiceName   string    `json:"invoice_name"`
	AmountTotal   float64   `json:"amount_total"`
	DateInvoice   time.Time `json:"date_invoice"`
}

// DailyInvoiceSummary represents invoice summary for a specific date
type DailyInvoiceSummary struct {
	Date         time.Time        `json:"date"`
	TotalAmount  float64          `json:"total_amount"`
	InvoiceCount int              `json:"invoice_count"`
	Invoices     []InvoiceSummary `json:"invoices"`
}

// InvoiceSummaryResponse represents the paginated response for invoice summary
type InvoiceSummaryResponse struct {
	Items      []DailyInvoiceSummary `json:"items"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalItems int                   `json:"total_items"`
	TotalPages int                   `json:"total_pages"`
}

// PeriodInvoiceSummaryResponse represents the response for period-based invoice summary
type PeriodInvoiceSummaryResponse struct {
	Period       string                `json:"period"`
	PeriodType   PeriodType            `json:"period_type"`
	DateRange    DateRange             `json:"date_range"`
	Items        []DailyInvoiceSummary `json:"items"`
	TotalAmount  float64               `json:"total_amount"`
	InvoiceCount int                   `json:"invoice_count"`
	AverageDaily float64               `json:"average_daily"`
}
