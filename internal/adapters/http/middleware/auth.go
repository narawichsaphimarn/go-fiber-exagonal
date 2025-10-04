package middleware

import (
	"strings"

	"example.com/practice/fiber/internal/ports"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	TokenProvider ports.TokenProvider
	Repo ports.UserRepository
}

func NewAuthMiddleware(tokenProvider ports.TokenProvider) *AuthMiddleware {
	return &AuthMiddleware{
		TokenProvider: tokenProvider,
	}
}

func (m *AuthMiddleware) Protect(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "missing token",
		})
	}
	tokenTrim := strings.TrimPrefix(token, "Bearer ")
	userId, err := m.TokenProvider.ValidateToken(tokenTrim)
	if err != nil || userId == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token",
		})
	}
	c.Locals("userId", userId)
	return c.Next()
}

