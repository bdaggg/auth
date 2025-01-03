package middleware

import (
	"time"

	"auth-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

func RequestMetrics(monitoringService *service.MonitoringService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		monitoringService.RecordRequestDuration(duration)
		return err
	}
}
