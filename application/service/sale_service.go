package service

import (
	"nerp_wrapper/domain/entity"
	"nerp_wrapper/infrastructure/odoo"
)

type SaleService struct {
	saleRepo *odoo.OdooSaleRepository
}

func NewSaleService(saleRepo *odoo.OdooSaleRepository) *SaleService {
	return &SaleService{saleRepo: saleRepo}
}

func (s *SaleService) GetAllSaleOrders(page, pageSize int) (*entity.SaleOrderPagination, error) {
	return s.saleRepo.GetAllSaleOrders(page, pageSize)
}

// GetDailySalesSummary retrieves daily sales summary with pagination
func (s *SaleService) GetDailySalesSummary(page, pageSize int) (*entity.SalesSummaryResponse, error) {
	return s.saleRepo.GetDailySalesSummary(page, pageSize)
}
