package service

import (
	"nerp_wrapper/domain/entity"
	"nerp_wrapper/infrastructure/odoo"
)

// SaleService handles business logic for sale orders
type SaleService struct {
	saleRepo *odoo.OdooSaleRepository
}

// NewSaleService creates a new instance of SaleService
func NewSaleService(saleRepo *odoo.OdooSaleRepository) *SaleService {
	return &SaleService{saleRepo: saleRepo}
}

// GetAllSaleOrders retrieves all sale orders
func (s *SaleService) GetAllSaleOrders() ([]*entity.SaleOrder, error) {
	return s.saleRepo.GetAllSaleOrders()
}
