package repository

import "nerp_wrapper/domain/entity"

// AuthRepository defines the interface for authentication operations
type AuthRepository interface {
	// Login authenticates a user with the given credentials
	Login(username, password string) (*entity.User, error)

	// Logout handles user logout
	Logout() error

	// GetUserInfo retrieves user information by ID
	GetUserInfo(userID int64) (*entity.User, error)
}
