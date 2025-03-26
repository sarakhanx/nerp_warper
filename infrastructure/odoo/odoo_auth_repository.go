package odoo

import (
	"fmt"
	"nerp_wrapper/domain/entity"
	"time"

	odoo "github.com/skilld-labs/go-odoo"
)

// OdooAuthRepository implements AuthRepository interface using Odoo
type OdooAuthRepository struct {
	client *odoo.Client
}

// NewOdooAuthRepository creates a new instance of OdooAuthRepository
func NewOdooAuthRepository(adminUsername, adminPassword, database, url string) (*OdooAuthRepository, error) {
	// Create Odoo client with admin credentials for system connection
	client, err := odoo.NewClient(&odoo.ClientConfig{
		Admin:    adminUsername,
		Password: adminPassword,
		Database: database,
		URL:      url,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Odoo client: %v", err)
	}
	return &OdooAuthRepository{client: client}, nil
}

// Login authenticates a user with Odoo
func (r *OdooAuthRepository) Login(username, password string) (*entity.User, error) {
	// Use admin client to verify user credentials
	criteria := odoo.NewCriteria().
		Add("login", "=", username).
		Add("password", "=", password)

	users, err := r.client.FindResUsers(criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search for user: %v", err)
	}

	if users == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login in Odoo
	now := time.Now()
	users.LoginDate = odoo.NewTime(now)
	err = r.client.UpdateResUsers(users)
	if err != nil {
		return nil, fmt.Errorf("failed to update last login: %v", err)
	}

	// Convert Odoo user to domain entity
	return entity.NewUser(
		users.Id.Get(),
		users.Login.Get(),
		users.Email.Get(),
		users.Active.Get(),
		now,
	), nil
}

// Logout handles user logout
func (r *OdooAuthRepository) Logout() error {
	// In Odoo, logout is typically handled on the client side
	return nil
}

// GetUserInfo retrieves user information
func (r *OdooAuthRepository) GetUserInfo(userID int64) (*entity.User, error) {
	// Use admin client to get user information
	odooUser, err := r.client.GetResUsers(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	// Convert last login time
	var lastLogin time.Time
	if odooUser.LoginDate != nil {
		lastLogin = odooUser.LoginDate.Get()
	}

	return entity.NewUser(
		odooUser.Id.Get(),
		odooUser.Login.Get(),
		odooUser.Email.Get(),
		odooUser.Active.Get(),
		lastLogin,
	), nil
}
