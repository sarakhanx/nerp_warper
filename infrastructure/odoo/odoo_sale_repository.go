package odoo

import (
	"fmt"
	"nerp_wrapper/domain/entity"
	"sort"
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

// GetAllSaleOrders retrieves sale orders from Odoo with pagination
func (r *OdooSaleRepository) GetAllSaleOrders(page, pageSize int) (*entity.SaleOrderPagination, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 500 {
		pageSize = 500
	}

	offset := (page - 1) * pageSize

	criteria := odoo.NewCriteria().Add("state", "!=", "cancel")
	totalCount, err := r.client.Count("sale.order", criteria, odoo.NewOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	searchOptions := odoo.NewOptions().
		Limit(pageSize).
		Offset(offset)

	ids, err := r.client.Search("sale.order", criteria, searchOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to search sale orders: %v", err)
	}

	if len(ids) == 0 {
		return &entity.SaleOrderPagination{
			Items:      []*entity.SaleOrder{},
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(totalCount),
			TotalPages: totalPages,
		}, nil
	}

	orders := []odoo.SaleOrder{}
	readOptions := odoo.NewOptions().FetchFields(
		"id", "name", "partner_id", "date_order", "amount_total", "state",
		"partner_invoice_id", "partner_shipping_id", "validity_date",
		"client_order_ref", "user_id", "note",
	)
	if err := r.client.Read("sale.order", ids, readOptions, &orders); err != nil {
		return nil, fmt.Errorf("failed to read sale orders: %v", err)
	}

	partnerIDs := make(map[int64]bool)
	userIDs := make(map[int64]bool)
	for _, order := range orders {
		if order.PartnerId != nil {
			partnerIDs[order.PartnerId.Get()] = true
		}
		if order.UserId != nil {
			userIDs[order.UserId.Get()] = true
		}
	}

	partners := make(map[int64]*odoo.ResPartner)
	if len(partnerIDs) > 0 {
		partnerIDsList := make([]int64, 0, len(partnerIDs))
		for id := range partnerIDs {
			partnerIDsList = append(partnerIDsList, id)
		}
		var partnerRecords []odoo.ResPartner
		partnerOptions := odoo.NewOptions().FetchFields("id", "name", "vat", "phone", "mobile")
		if err := r.client.Read("res.partner", partnerIDsList, partnerOptions, &partnerRecords); err == nil {
			for i := range partnerRecords {
				partners[partnerRecords[i].Id.Get()] = &partnerRecords[i]
			}
		}
	}

	users := make(map[int64]*odoo.ResUsers)
	if len(userIDs) > 0 {
		userIDsList := make([]int64, 0, len(userIDs))
		for id := range userIDs {
			userIDsList = append(userIDsList, id)
		}
		var userRecords []odoo.ResUsers
		userOptions := odoo.NewOptions().FetchFields("id", "name")
		if err := r.client.Read("res.users", userIDsList, userOptions, &userRecords); err == nil {
			for i := range userRecords {
				users[userRecords[i].Id.Get()] = &userRecords[i]
			}
		}
	}

	var result []*entity.SaleOrder
	for _, order := range orders {
		var partnerName, partnerVat, partnerPhone, partnerMobile string
		var partnerID int64
		if order.PartnerId != nil {
			partnerID = order.PartnerId.Get()
			if partner, exists := partners[partnerID]; exists {
				if partner.Name != nil {
					partnerName = partner.Name.Get()
				}
				if partner.Vat != nil {
					partnerVat = partner.Vat.Get()
				}
				if partner.Phone != nil {
					partnerPhone = partner.Phone.Get()
				}
				if partner.Mobile != nil {
					partnerMobile = partner.Mobile.Get()
				}
			}
		}

		var salespersonID int64
		var salespersonName string
		if order.UserId != nil {
			salespersonID = order.UserId.Get()
			if user, exists := users[salespersonID]; exists && user.Name != nil {
				salespersonName = user.Name.Get()
			}
		}

		saleOrder := &entity.SaleOrder{
			ID:              order.Id.Get(),
			Name:            order.Name.Get(),
			Partner:         partnerID,
			PartnerName:     partnerName,
			PartnerVat:      partnerVat,
			PartnerPhone:    partnerPhone,
			PartnerMobile:   partnerMobile,
			DateOrder:       order.DateOrder.Get(),
			AmountTotal:     order.AmountTotal.Get(),
			SalespersonID:   salespersonID,
			SalespersonName: salespersonName,
			Note:            order.Note.Get(),
		}

		if order.State != nil {
			if s, ok := order.State.Get().(string); ok {
				saleOrder.State = s
			}
		}

		if order.ValidityDate != nil {
			saleOrder.ValidityDate = order.ValidityDate.Get()
		}

		if order.ClientOrderRef != nil {
			saleOrder.ClientOrderRef = order.ClientOrderRef.Get()
		}

		if order.PartnerInvoiceId != nil {
			saleOrder.PartnerInvoiceID = order.PartnerInvoiceId.Get()
		}

		if order.PartnerShippingId != nil {
			saleOrder.PartnerShippingID = order.PartnerShippingId.Get()
		}

		result = append(result, saleOrder)
	}

	return &entity.SaleOrderPagination{
		Items:      result,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: int(totalCount),
		TotalPages: totalPages,
	}, nil
}

// GetDailySalesSummary retrieves daily sales summary with pagination
func (r *OdooSaleRepository) GetDailySalesSummary(page, pageSize int) (*entity.SalesSummaryResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 500 {
		pageSize = 500
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Search for confirmed sales orders
	criteria := odoo.NewCriteria().Add("state", "=", "sale")

	// Get total count
	totalCount, err := r.client.Count("sale.order", criteria, odoo.NewOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	// Search with pagination
	searchOptions := odoo.NewOptions().
		Limit(pageSize).
		Offset(offset).
		FetchFields("id", "name", "date_order", "amount_total")

	// Execute search and read in one call
	var records []odoo.SaleOrder
	if err := r.client.SearchRead("sale.order", criteria, searchOptions, &records); err != nil {
		return nil, fmt.Errorf("failed to search sale orders: %v", err)
	}

	if len(records) == 0 {
		return &entity.SalesSummaryResponse{
			Items:      []entity.DailySalesSummary{},
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(totalCount),
			TotalPages: totalPages,
		}, nil
	}

	// Group orders by date
	dateMap := make(map[string]*entity.DailySalesSummary)
	for _, order := range records {
		orderDate := order.DateOrder.Get()
		dateStr := orderDate.Format("2006-01-02")

		if _, exists := dateMap[dateStr]; !exists {
			dateMap[dateStr] = &entity.DailySalesSummary{
				Date:        orderDate.Truncate(24 * time.Hour),
				TotalAmount: 0,
				OrderCount:  0,
				Orders:      []entity.SaleOrderSummary{},
			}
		}

		summary := dateMap[dateStr]
		summary.TotalAmount += order.AmountTotal.Get()
		summary.OrderCount++

		orderSummary := entity.SaleOrderSummary{
			OrderNumber: summary.OrderCount,
			OrderID:     order.Id.Get(),
			OrderName:   order.Name.Get(),
			AmountTotal: order.AmountTotal.Get(),
			DateOrder:   orderDate,
		}
		summary.Orders = append(summary.Orders, orderSummary)
	}

	// Convert map to sorted slice
	var summaries []entity.DailySalesSummary
	for _, summary := range dateMap {
		summaries = append(summaries, *summary)
	}

	// Sort by date descending
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Date.After(summaries[j].Date)
	})

	return &entity.SalesSummaryResponse{
		Items:      summaries,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: int(totalCount),
		TotalPages: totalPages,
	}, nil
}
