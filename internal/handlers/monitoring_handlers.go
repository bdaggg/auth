package handlers

import (
	"auth-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

func GetMetrics(monitoringService *service.MonitoringService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		metrics, err := monitoringService.GetMetrics(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(metrics)
	}
}

func GetActiveUsers(monitoringService *service.MonitoringService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		count, err := monitoringService.GetActiveUsersCount(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"active_users": count,
		})
	}
}

func GetBlockedUsers(monitoringService *service.MonitoringService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		count, err := monitoringService.GetBlockedUsersCount(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"blocked_users": count,
		})
	}
}
