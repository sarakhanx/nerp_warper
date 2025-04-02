package odoo

import (
	"fmt"
	"nerp_wrapper/domain/entity"
	"sort"
	"time"

	odoo "github.com/skilld-labs/go-odoo"
)

// OdooInvoiceRepository handles invoice operations with Odoo
type OdooInvoiceRepository struct {
	client *odoo.Client
}

// NewOdooInvoiceRepository creates a new instance of OdooInvoiceRepository
func NewOdooInvoiceRepository(client *odoo.Client) *OdooInvoiceRepository {
	return &OdooInvoiceRepository{client: client}
}

// GetAllInvoices retrieves invoices from Odoo with pagination
func (r *OdooInvoiceRepository) GetAllInvoices(page, pageSize int) (*entity.InvoicePagination, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 500 {
		pageSize = 500
	}

	offset := (page - 1) * pageSize

	criteria := odoo.NewCriteria().Add("state", "!=", "cancel")
	totalCount, err := r.client.Count("account.invoice.report", criteria, odoo.NewOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	searchOptions := odoo.NewOptions().
		Limit(pageSize).
		Offset(offset)

	var invoices []odoo.AccountInvoiceReport
	if err := r.client.SearchRead("account.invoice.report", criteria, searchOptions, &invoices); err != nil {
		return nil, fmt.Errorf("failed to search invoices: %v", err)
	}

	if len(invoices) == 0 {
		return &entity.InvoicePagination{
			Items:      []*entity.Invoice{},
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(totalCount),
			TotalPages: totalPages,
		}, nil
	}

	// Collect related record IDs
	partnerIDs := make(map[int64]bool)
	journalIDs := make(map[int64]bool)
	currencyIDs := make(map[int64]bool)

	for _, invoice := range invoices {
		if invoice.PartnerId != nil {
			partnerIDs[invoice.PartnerId.Get()] = true
		}
		if invoice.JournalId != nil {
			journalIDs[invoice.JournalId.Get()] = true
		}
		if invoice.CurrencyId != nil {
			currencyIDs[invoice.CurrencyId.Get()] = true
		}
	}

	// Fetch partners
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

	// Fetch journals
	journals := make(map[int64]*odoo.AccountJournal)
	if len(journalIDs) > 0 {
		journalIDsList := make([]int64, 0, len(journalIDs))
		for id := range journalIDs {
			journalIDsList = append(journalIDsList, id)
		}
		var journalRecords []odoo.AccountJournal
		journalOptions := odoo.NewOptions().FetchFields("id", "name")
		if err := r.client.Read("account.journal", journalIDsList, journalOptions, &journalRecords); err == nil {
			for i := range journalRecords {
				journals[journalRecords[i].Id.Get()] = &journalRecords[i]
			}
		}
	}

	// Fetch currencies
	currencies := make(map[int64]*odoo.ResCurrency)
	if len(currencyIDs) > 0 {
		currencyIDsList := make([]int64, 0, len(currencyIDs))
		for id := range currencyIDs {
			currencyIDsList = append(currencyIDsList, id)
		}
		var currencyRecords []odoo.ResCurrency
		currencyOptions := odoo.NewOptions().FetchFields("id", "name")
		if err := r.client.Read("res.currency", currencyIDsList, currencyOptions, &currencyRecords); err == nil {
			for i := range currencyRecords {
				currencies[currencyRecords[i].Id.Get()] = &currencyRecords[i]
			}
		}
	}

	var result []*entity.Invoice
	for _, invoice := range invoices {
		var partnerName, partnerVat, partnerPhone, partnerMobile string
		var partnerID int64
		if invoice.PartnerId != nil {
			partnerID = invoice.PartnerId.Get()
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

		var journalID int64
		var journalName string
		if invoice.JournalId != nil {
			journalID = invoice.JournalId.Get()
			if journal, exists := journals[journalID]; exists && journal.Name != nil {
				journalName = journal.Name.Get()
			}
		}

		var currencyID int64
		var currencyName string
		if invoice.CurrencyId != nil {
			currencyID = invoice.CurrencyId.Get()
			if currency, exists := currencies[currencyID]; exists && currency.Name != nil {
				currencyName = currency.Name.Get()
			}
		}

		inv := &entity.Invoice{
			ID:             invoice.Id.Get(),
			Name:           invoice.DisplayName.Get(),
			Partner:        partnerID,
			PartnerName:    partnerName,
			PartnerVat:     partnerVat,
			PartnerPhone:   partnerPhone,
			PartnerMobile:  partnerMobile,
			AmountUntaxed:  invoice.PriceTotal.Get() - invoice.UserCurrencyPriceTotal.Get(),
			AmountTax:      invoice.UserCurrencyPriceTotal.Get() - invoice.PriceTotal.Get(),
			AmountTotal:    invoice.PriceTotal.Get(),
			AmountResidual: invoice.Residual.Get(),
			JournalID:      journalID,
			JournalName:    journalName,
			CurrencyID:     currencyID,
			CurrencyName:   currencyName,
		}

		if invoice.State != nil {
			if s, ok := invoice.State.Get().(string); ok {
				inv.State = s
			}
		}

		if invoice.Type != nil {
			if t, ok := invoice.Type.Get().(string); ok {
				inv.Type = t
			}
		}

		if invoice.Date != nil {
			inv.DateInvoice = invoice.Date.Get()
		}

		if invoice.DateDue != nil {
			inv.DateDue = invoice.DateDue.Get()
		}

		result = append(result, inv)
	}

	return &entity.InvoicePagination{
		Items:      result,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: int(totalCount),
		TotalPages: totalPages,
	}, nil
}

// GetDailyInvoiceSummary retrieves daily invoice summary with pagination
func (r *OdooInvoiceRepository) GetDailyInvoiceSummary(page, pageSize int) (*entity.InvoiceSummaryResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 500 {
		pageSize = 500
	}

	offset := (page - 1) * pageSize

	criteria := odoo.NewCriteria().
		Add("state", "=", "posted")

	totalCount, err := r.client.Count("account.invoice.report", criteria, odoo.NewOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	searchOptions := odoo.NewOptions().
		Limit(pageSize).
		Offset(offset)

	var records []odoo.AccountInvoiceReport
	if err := r.client.SearchRead("account.invoice.report", criteria, searchOptions, &records); err != nil {
		return nil, fmt.Errorf("failed to search invoices: %v", err)
	}

	if len(records) == 0 {
		return &entity.InvoiceSummaryResponse{
			Items:      []entity.DailyInvoiceSummary{},
			Page:       page,
			PageSize:   pageSize,
			TotalItems: int(totalCount),
			TotalPages: totalPages,
		}, nil
	}

	dateMap := make(map[string]*entity.DailyInvoiceSummary)
	for _, invoice := range records {
		invoiceDate := invoice.Date.Get()
		dateStr := invoiceDate.Format("2006-01-02")

		if _, exists := dateMap[dateStr]; !exists {
			dateMap[dateStr] = &entity.DailyInvoiceSummary{
				Date:         invoiceDate.Truncate(24 * time.Hour),
				TotalAmount:  0,
				InvoiceCount: 0,
				Invoices:     []entity.InvoiceSummary{},
			}
		}

		summary := dateMap[dateStr]
		summary.TotalAmount += invoice.PriceTotal.Get()
		summary.InvoiceCount++

		invoiceSummary := entity.InvoiceSummary{
			InvoiceNumber: summary.InvoiceCount,
			InvoiceID:     invoice.Id.Get(),
			InvoiceName:   invoice.DisplayName.Get(),
			AmountTotal:   invoice.PriceTotal.Get(),
			DateInvoice:   invoiceDate,
		}
		summary.Invoices = append(summary.Invoices, invoiceSummary)
	}

	var summaries []entity.DailyInvoiceSummary
	for _, summary := range dateMap {
		summaries = append(summaries, *summary)
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Date.After(summaries[j].Date)
	})

	return &entity.InvoiceSummaryResponse{
		Items:      summaries,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: int(totalCount),
		TotalPages: totalPages,
	}, nil
}

// GetPeriodInvoiceSummary retrieves invoice summary for a specific period
func (r *OdooInvoiceRepository) GetPeriodInvoiceSummary(periodType entity.PeriodType, customStartDate, customEndDate *time.Time) (*entity.PeriodInvoiceSummaryResponse, error) {
	now := time.Now()
	var startDate, endDate time.Time

	if customStartDate != nil && customEndDate != nil {
		startDate = *customStartDate
		endDate = *customEndDate
	} else {
		endDate = now
		switch periodType {
		case entity.PeriodTypeDay:
			startDate = now.AddDate(0, 0, -1)
		case entity.PeriodTypeWeek:
			startDate = now.AddDate(0, 0, -7)
		case entity.PeriodTypeMonth:
			startDate = now.AddDate(0, 0, -30)
		case entity.PeriodTypeQuarter:
			startDate = now.AddDate(0, 0, -90)
		case entity.PeriodTypeMonthly:
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
		case entity.PeriodTypeYearly:
			startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
			endDate = startDate.AddDate(1, 0, 0).Add(-time.Second)
		default:
			startDate = now.AddDate(0, 0, -30)
		}
	}

	criteria := odoo.NewCriteria().
		Add("state", "=", "posted").
		Add("invoice_date", ">=", startDate.Format("2006-01-02")).
		Add("invoice_date", "<=", endDate.Format("2006-01-02"))

	searchOptions := odoo.NewOptions().
		FetchFields("id", "display_name", "invoice_date", "price_total")

	var records []odoo.AccountInvoiceReport
	if err := r.client.SearchRead("account.invoice.report", criteria, searchOptions, &records); err != nil {
		return nil, fmt.Errorf("failed to search invoices: %v", err)
	}

	dateMap := make(map[string]*entity.DailyInvoiceSummary)
	var totalAmount float64
	var totalInvoices int

	for _, invoice := range records {
		invoiceDate := invoice.Date.Get()
		dateStr := invoiceDate.Format("2006-01-02")

		if _, exists := dateMap[dateStr]; !exists {
			dateMap[dateStr] = &entity.DailyInvoiceSummary{
				Date:         invoiceDate.Truncate(24 * time.Hour),
				TotalAmount:  0,
				InvoiceCount: 0,
				Invoices:     []entity.InvoiceSummary{},
			}
		}

		summary := dateMap[dateStr]
		amount := invoice.PriceTotal.Get()
		summary.TotalAmount += amount
		summary.InvoiceCount++
		totalAmount += amount
		totalInvoices++

		invoiceSummary := entity.InvoiceSummary{
			InvoiceNumber: summary.InvoiceCount,
			InvoiceID:     invoice.Id.Get(),
			InvoiceName:   invoice.DisplayName.Get(),
			AmountTotal:   amount,
			DateInvoice:   invoiceDate,
		}
		summary.Invoices = append(summary.Invoices, invoiceSummary)
	}

	var summaries []entity.DailyInvoiceSummary
	for _, summary := range dateMap {
		summaries = append(summaries, *summary)
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Date.After(summaries[j].Date)
	})

	days := endDate.Sub(startDate).Hours() / 24
	if days < 1 {
		days = 1
	}

	var periodDesc string
	switch periodType {
	case entity.PeriodTypeDay:
		periodDesc = "Last 24 Hours"
	case entity.PeriodTypeWeek:
		periodDesc = "Last 7 Days"
	case entity.PeriodTypeMonth:
		periodDesc = "Last 30 Days"
	case entity.PeriodTypeQuarter:
		periodDesc = "Last 90 Days"
	case entity.PeriodTypeMonthly:
		periodDesc = endDate.Format("January 2006")
	case entity.PeriodTypeYearly:
		periodDesc = endDate.Format("2006")
	default:
		periodDesc = fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}

	return &entity.PeriodInvoiceSummaryResponse{
		Period:       periodDesc,
		PeriodType:   periodType,
		DateRange:    entity.DateRange{StartDate: startDate, EndDate: endDate},
		Items:        summaries,
		TotalAmount:  totalAmount,
		InvoiceCount: totalInvoices,
		AverageDaily: totalAmount / days,
	}, nil
}
