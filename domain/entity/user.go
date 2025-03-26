package entity

import "time"

// User represents the user entity in our domain
type User struct {
	ID        int64
	Username  string
	Email     string
	Active    bool
	LastLogin time.Time
}

// NewUser creates a new User instance
func NewUser(id int64, username, email string, active bool, lastLogin time.Time) *User {
	return &User{
		ID:        id,
		Username:  username,
		Email:     email,
		Active:    active,
		LastLogin: lastLogin,
	}
}
