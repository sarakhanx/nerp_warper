package service

import (
	"fmt"
	"nerp_wrapper/domain/entity"
	"nerp_wrapper/domain/repository"
)

// AuthService handles authentication business logic
type AuthService struct {
	authRepo repository.AuthRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(authRepo repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

// Login handles user authentication
func (s *AuthService) Login(username, password string) (*entity.User, error) {
	user, err := s.authRepo.Login(username, password)
	if err != nil {
		return nil, fmt.Errorf("login failed: %v", err)
	}
	return user, nil
}

// Logout handles user logout
func (s *AuthService) Logout() error {
	return s.authRepo.Logout()
}

// GetUserInfo retrieves user information
func (s *AuthService) GetUserInfo(userID int64) (*entity.User, error) {
	user, err := s.authRepo.GetUserInfo(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	return user, nil
}
