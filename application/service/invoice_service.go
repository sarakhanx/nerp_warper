package service

import (
	"nerp_wrapper/domain/entity"
	"nerp_wrapper/infrastructure/odoo"
	"time"
)

type InvoiceService struct {
	invoiceRepo *odoo.OdooInvoiceRepository
}

func NewInvoiceService(invoiceRepo *odoo.OdooInvoiceRepository) *InvoiceService {
	return &InvoiceService{invoiceRepo: invoiceRepo}
}

func (s *InvoiceService) GetAllInvoices(page, pageSize int) (*entity.InvoicePagination, error) {
	return s.invoiceRepo.GetAllInvoices(page, pageSize)
}

// GetDailyInvoiceSummary retrieves daily invoice summary with pagination
func (s *InvoiceService) GetDailyInvoiceSummary(page, pageSize int) (*entity.InvoiceSummaryResponse, error) {
	return s.invoiceRepo.GetDailyInvoiceSummary(page, pageSize)
}

// GetPeriodInvoiceSummary retrieves invoice summary for a specific period
func (s *InvoiceService) GetPeriodInvoiceSummary(periodType entity.PeriodType, startDate, endDate *time.Time) (*entity.PeriodInvoiceSummaryResponse, error) {
	return s.invoiceRepo.GetPeriodInvoiceSummary(periodType, startDate, endDate)
}
