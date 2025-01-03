package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RateLimit(redisClient *redis.Client, maxRequests int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		key := "rate_limit:" + ip

		ctx := context.Background()
		pipe := redisClient.Pipeline()

		// İstek sayısını artır ve TTL ayarla
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)

		cmds, err := pipe.Exec(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "rate limit kontrolü yapılamadı",
			})
		}

		// İstek sayısını kontrol et
		count := cmds[0].(*redis.IntCmd).Val()
		if count > int64(maxRequests) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "çok fazla istek gönderildi, lütfen bekleyin",
			})
		}

		return c.Next()
	}
}
