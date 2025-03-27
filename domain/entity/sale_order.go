package entity

import (
	"time"
	// odoo "github.com/skilld-labs/go-odoo"
)

// SaleOrder represents a sale order entity
type SaleOrder struct {
	ID          int64
	Name        string
	Partner     int64
	DateOrder   time.Time
	AmountTotal float64
	State       string
}

type SaleOrderRepository interface {
	GetAllSaleOrders() ([]*SaleOrder, error)
}
