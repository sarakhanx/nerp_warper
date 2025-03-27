package repository

import (
	"fmt"

	"github.com/skilld-labs/go-odoo"
)

// OdooClient wraps the go-odoo client to use our custom types
type OdooClient struct {
	client *odoo.Client
}

// NewOdooClient creates a new instance of OdooClient
func NewOdooClient(client *odoo.Client) *OdooClient {
	return &OdooClient{
		client: client,
	}
}

// SearchRead performs a search and read operation on the Odoo model
func (c *OdooClient) SearchRead(model string, criteria *Criteria, options *Options, context interface{}) (SearchResult, error) {
	// Create domain array for search criteria
	domain := criteria.Conditions

	// Create options map
	opts := map[string]interface{}{
		"fields": options.Fields,
	}
	if options.Limit > 0 {
		opts["limit"] = options.Limit
	}

	// Execute search_read with proper parameters
	result, err := c.client.ExecuteKw(model, "search_read", []interface{}{
		domain,
		opts,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search_read: %v", err)
	}

	// Convert result to SearchResult
	if records, ok := result.([]interface{}); ok {
		searchResult := make(SearchResult, len(records))
		for i, record := range records {
			if recordMap, ok := record.(map[string]interface{}); ok {
				searchResult[i] = recordMap
			}
		}
		return searchResult, nil
	}

	return nil, fmt.Errorf("unexpected result type: %T", result)
}

// GetClient returns the underlying odoo client
func (c *OdooClient) GetClient() *odoo.Client {
	return c.client
}
