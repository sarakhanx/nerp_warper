package entity

import (
	"time"
)

type Invoice struct {
	ID             int64
	Name           string
	Partner        int64
	PartnerName    string
	PartnerVat     string
	PartnerPhone   string
	PartnerMobile  string
	DateInvoice    time.Time
	DateDue        time.Time
	Reference      string
	AmountUntaxed  float64
	AmountTax      float64
	AmountTotal    float64
	AmountResidual float64
	State          string
	Type           string
	JournalID      int64
	JournalName    string
	CurrencyID     int64
	CurrencyName   string
	Note           string
}

type InvoicePagination struct {
	Items      []*Invoice `json:"items"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalItems int        `json:"total_items"`
	TotalPages int        `json:"total_pages"`
}
