package entity

import (
	"time"
)

type SaleOrder struct {
	ID                int64
	Name              string
	Partner           int64
	PartnerName       string
	PartnerVat        string
	PartnerPhone      string
	PartnerMobile     string
	PartnerInvoiceID  int64
	PartnerShippingID int64
	DateOrder         time.Time
	ValidityDate      time.Time
	ClientOrderRef    string
	SalespersonID     int64
	SalespersonName   string
	AmountTotal       float64
	State             string
	Note              string
}

type SaleOrderPagination struct {
	Items      []*SaleOrder `json:"items"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalItems int          `json:"total_items"`
	TotalPages int          `json:"total_pages"`
}
