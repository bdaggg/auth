package handlers

import (
	"auth-service/internal/domain/entity"
	"auth-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

func ChangePassword(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz istek formatı",
			})
		}

		userID := c.Locals("claims").(*entity.TokenClaims).UserID
		if err := authService.ChangePassword(c.Context(), userID, input.OldPassword, input.NewPassword); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func Enable2FA(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("claims").(*entity.TokenClaims).UserID
		secret, qrCode, err := authService.Enable2FA(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"secret":  secret,
			"qr_code": qrCode,
		})
	}
}

func Verify2FA(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Code string `json:"code"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz istek formatı",
			})
		}

		userID := c.Locals("claims").(*entity.TokenClaims).UserID
		if err := authService.Verify2FA(c.Context(), userID, input.Code); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func GetAuditLogs(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("claims").(*entity.TokenClaims).UserID
		limit := c.QueryInt("limit", 10)
		offset := c.QueryInt("offset", 0)

		logs, err := authService.GetAuditLogs(c.Context(), userID, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(logs)
	}
}
