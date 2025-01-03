package handlers

import (
	"auth-service/internal/domain/entity"
	"auth-service/internal/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func BlockUser(securityService *service.SecurityService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")
		var input struct {
			Reason string `json:"reason"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz istek formatı",
			})
		}

		adminID := c.Locals("claims").(*entity.TokenClaims).UserID
		if err := securityService.BlockUser(c.Context(), userID, adminID, input.Reason); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func UnblockUser(securityService *service.SecurityService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")
		adminID := c.Locals("claims").(*entity.TokenClaims).UserID

		if err := securityService.UnblockUser(c.Context(), userID, adminID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func GetSecurityAlerts(securityService *service.SecurityService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		from := c.QueryTime("from", time.Now().AddDate(0, 0, -7))
		to := c.QueryTime("to", time.Now())

		alerts, err := securityService.GetSecurityAlerts(c.Context(), from, to)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(alerts)
	}
}

func GetSuspiciousActivities(securityService *service.SecurityService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		threshold := c.QueryInt("threshold", 5)
		activities, err := securityService.GetSuspiciousActivities(c.Context(), threshold)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(activities)
	}
}
