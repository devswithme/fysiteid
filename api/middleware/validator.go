package middleware

import (
	"fybe/model/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidatorMiddleware[T any](validator *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse[any]{
				Data: nil,
				Meta: dto.ApiMeta{
					Success: false,
					Message: "Bad body request",
				},
			})
		}

		if err := validator.Struct(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse[any]{
				Data: nil,
				Meta: dto.ApiMeta{
					Success: false,
					Message: "Bad validation request",
				},
			})
		}

		c.Locals("body", body)
		return c.Next()
	}
}
