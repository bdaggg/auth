package middleware

import (
	"auth-service/internal/domain/entity"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(roles ...entity.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*entity.TokenClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "yetkilendirme başarısız",
			})
		}

		for _, role := range roles {
			if claims.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "bu işlem için yetkiniz yok",
		})
	}
}
