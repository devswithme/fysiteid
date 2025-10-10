package middleware

import (
	"context"
	"fmt"
	"fybe/model/dto"
	"strconv"

	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
)

func LimiterMiddleware(limiter *redis_rate.Limiter, limit redis_rate.Limit) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		userID := c.Locals("user_id")

		var userIDStr string
		if userID == nil {
			userIDStr = "anonymous"
		} else {
			switch v := userID.(type) {
			case uint:
				userIDStr = strconv.FormatUint(uint64(v), 10)
			case int:
				userIDStr = strconv.Itoa(v)
			case string:
				userIDStr = v
			default:
				userIDStr = fmt.Sprintf("%v", v)
			}
		}

		key := "rate_limit:" + c.IP() + userIDStr

		res, err := limiter.Allow(ctx, key, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
				Data: nil,
				Meta: dto.ApiMeta{
					Success:   false,
					Message:   "error handling requests",
					RequestID: c.Locals("request_id").(string),
				},
			})
		}

		if res.Allowed == 0 {
			return c.Status(fiber.StatusTooManyRequests).JSON(dto.APIResponse[any]{
				Data: nil,
				Meta: dto.ApiMeta{
					Success:   false,
					Message:   "too many requests",
					RequestID: c.Locals("request_id").(string),
				},
			})
		}

		return c.Next()
	}
}
