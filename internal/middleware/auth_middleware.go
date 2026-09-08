package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/l0ng7h0r/ecommerce/pkg/config"
	"github.com/l0ng7h0r/ecommerce/pkg/security"
)

type AuthMiddleware struct {
	cfg *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

func (m *AuthMiddleware) Auth() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header required"})
		}

		tokenStr := strings.TrimSpace(authHeader)
		if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		} else if strings.HasPrefix(tokenStr, "bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "bearer ")
		}

		claims, err := security.ValidateToken(tokenStr, m.cfg.JWTSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_roles", claims.Roles)

		return c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(role string) fiber.Handler {
	return func(c fiber.Ctx) error {
		rolesVal := c.Locals("user_roles")
		if rolesVal == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}

		roles, ok := rolesVal.([]string)
		if !ok {
			// In case deserialized as interface{} array
			if rawRoles, ok := rolesVal.([]interface{}); ok {
				for _, r := range rawRoles {
					if strRole, ok := r.(string); ok && strRole == role {
						return c.Next()
					}
				}
			}
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}

		for _, r := range roles {
			if r == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied: required role " + role})
	}
}
