package auth

import (
	"fmt"

	odoo "github.com/skilld-labs/go-odoo"
)

// OdooAuthService handles authentication with Odoo
type OdooAuthService struct {
	client *odoo.Client
}

// NewOdooAuthService creates a new instance of OdooAuthService
func NewOdooAuthService(admin, password, database, url string) (*OdooAuthService, error) {
	client, err := odoo.NewClient(&odoo.ClientConfig{
		Admin:    admin,
		Password: password,
		Database: database,
		URL:      url,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Odoo client: %v", err)
	}
	return &OdooAuthService{client: client}, nil
}

// Login authenticates a user with Odoo
func (s *OdooAuthService) Login(username, password string) (bool, error) {
	// Search for user with matching credentials
	criteria := odoo.NewCriteria().
		Add("login", "=", username).
		Add("password", "=", password)

	users, err := s.client.FindResUsers(criteria)
	if err != nil {
		return false, fmt.Errorf("failed to search for user: %v", err)
	}

	return users != nil, nil
}

// Logout handles user logout
func (s *OdooAuthService) Logout() error {
	// In Odoo, the logout is typically handled on the client side
	// We can just clear any session data if needed
	return nil
}

// GetUserInfo retrieves user information
func (s *OdooAuthService) GetUserInfo(userID int64) (*odoo.ResUsers, error) {
	user, err := s.client.GetResUsers(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	return user, nil
}
