package middleware

import (
	"strings"

	"auth-service/internal/domain/entity"
	"auth-service/pkg/security"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth(jwtManager *security.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Authorization header'ı al
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "authorization header eksik",
			})
		}

		// Bearer token'ı ayır
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "geçersiz authorization format",
			})
		}

		// Token'ı doğrula
		claims, err := jwtManager.ValidateToken(parts[1], entity.AccessToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "geçersiz token",
			})
		}

		// Claims'i context'e ekle
		c.Locals("claims", claims)
		return c.Next()
	}
}
