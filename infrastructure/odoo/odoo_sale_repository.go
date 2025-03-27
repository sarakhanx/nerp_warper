package odoo

import (
	"fmt"
	"log"
	"nerp_wrapper/domain/entity"
	"time"

	odoo "github.com/skilld-labs/go-odoo"
)

// OdooSaleRepository handles sale order operations with Odoo
type OdooSaleRepository struct {
	client *odoo.Client
}

// NewOdooSaleRepository creates a new instance of OdooSaleRepository
func NewOdooSaleRepository(client *odoo.Client) *OdooSaleRepository {
	return &OdooSaleRepository{client: client}
}

// GetAllSaleOrders retrieves all sale orders from Odoo
func (r *OdooSaleRepository) GetAllSaleOrders() ([]*entity.SaleOrder, error) {
	criteria := odoo.NewCriteria().Add("state", "!=", "cancel")
	ids, err := r.client.Search("sale.order", criteria, odoo.NewOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to search sale orders: %v", err)
	}
	log.Printf("Found %d sale orders\n", len(ids))
	if len(ids) == 0 {
		return nil, fmt.Errorf("no sale order records found")
	}
	log.Printf("All's good, IDs: %v", ids)

	// Read complete sale order data
	orders := []odoo.SaleOrder{}
	options := odoo.NewOptions().FetchFields("id", "name", "partner_id", "date_order", "amount_total", "state")
	if err := r.client.Read("sale.order", ids, options, &orders); err != nil {
		return nil, fmt.Errorf("failed to read sale orders: %v", err)
	}

	// Convert Odoo sale orders to domain entities
	var result []*entity.SaleOrder
	for _, order := range orders {
		// Get order ID
		var orderID int64
		if order.Id != nil {
			orderID = order.Id.Get()
		}

		// Get name
		var name string
		if order.Name != nil {
			name = order.Name.Get()
		}

		// Get partner ID
		var partnerID int64
		if order.PartnerId != nil {
			partnerID = order.PartnerId.Get()
		}

		// Get date order
		var dateOrder time.Time
		if order.DateOrder != nil {
			dateOrder = order.DateOrder.Get()
		}

		// Get amount total
		var amountTotal float64
		if order.AmountTotal != nil {
			amountTotal = order.AmountTotal.Get()
		}

		// Get state
		var state string
		if order.State != nil {
			if s, ok := order.State.Get().(string); ok {
				state = s
			}
		}

		// Create domain entity
		result = append(result, &entity.SaleOrder{
			ID:          orderID,
			Name:        name,
			Partner:     partnerID,
			DateOrder:   dateOrder,
			AmountTotal: amountTotal,
			State:       state,
		})
	}

	return result, nil
}
