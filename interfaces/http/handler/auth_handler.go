package handler

import (
	"fmt"
	"nerp_wrapper/application/dto"
	"nerp_wrapper/application/service"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Success: false,
			Message: "Invalid credentials",
		})
	}

	return c.JSON(dto.LoginResponse{
		Success: true,
		Message: "Login successful",
		User: &dto.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Active:   user.Active,
		},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	if err := h.authService.Logout(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Success: false,
			Message: "Failed to logout",
		})
	}

	return c.JSON(dto.ErrorResponse{
		Success: true,
		Message: "Logout successful",
	})
}

// GetUserInfo retrieves user information
func (h *AuthHandler) GetUserInfo(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Message: "User ID is required",
		})
	}

	// Convert string to int64
	var id int64
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
	}

	user, err := h.authService.GetUserInfo(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Success: false,
			Message: "User not found",
		})
	}

	return c.JSON(dto.LoginResponse{
		Success: true,
		Message: "User info retrieved successfully",
		User: &dto.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Active:   user.Active,
		},
	})
}
