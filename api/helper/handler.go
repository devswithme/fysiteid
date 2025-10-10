package helper

import (
	"fybe/model/dto"

	"github.com/gofiber/fiber/v2"
)

func JWTErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse[any]{
		Data: nil,
		Meta: dto.ApiMeta{
			Message:   "Unauthorized",
			Success:   false,
			RequestID: c.Locals("request_id").(string),
		},
	})
}
