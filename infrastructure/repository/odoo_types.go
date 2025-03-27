package repository

import "github.com/skilld-labs/go-odoo"

// Criteria represents search criteria for Odoo API
type Criteria struct {
	Conditions []interface{}
}

// Options represents options for Odoo API
type Options struct {
	Fields []string
	Limit  int
}

// SearchResult represents the result from Odoo API
type SearchResult []map[string]interface{}

// ToOdooOptions converts our Options to odoo.Options
func (o *Options) ToOdooOptions() *odoo.Options {
	opts := odoo.NewOptions()
	opts.FetchFields(o.Fields...)
	if o.Limit > 0 {
		opts.Limit(o.Limit)
	}
	return opts
}
