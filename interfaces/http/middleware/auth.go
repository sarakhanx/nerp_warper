package middleware

import (
	"nerp_wrapper/application/dto"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// RequireAuth middleware checks for valid authentication token
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Success: false,
				Message: "Authorization header is required",
			})
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Success: false,
				Message: "Invalid authorization header format",
			})
		}

		// TODO: Validate token and get user information
		// For now, we'll just pass the token to the next handler
		c.Locals("token", parts[1])

		return c.Next()
	}
}
