package controller

import (
	"context"
	"fybe/helper"
	"fybe/model/dto"
	"fybe/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserController interface {
	Me(c *fiber.Ctx) error
	GetByUsername(c *fiber.Ctx) error
	GetByTicketID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type userController struct {
	userService service.UserService
	logger      *zap.Logger
}

func NewUserController(userService service.UserService, logger *zap.Logger) UserController {
	return &userController{
		userService: userService,
		logger:      logger,
	}
}

// @Summary		Get profile
// @Description	Get current logged-in user
// @Tags			user
// @Produce		json
// @Success		200	{object}	dto.APIResponse[dto.UserGet]	"Successfully get current user"
// @Failure		500	{object}	dto.APIResponse[any]			"Failed to get current user"
// @Router			/user/me [get]
func (t *userController) Me(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	userID := c.Locals("userID").(uint)
	user, err := t.userService.GetById(ctx, userID)

	if err != nil {
		t.logger.Error("failed get current user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to get current user",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[dto.UserGet]{
		Data: *user,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully get current user",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Get public by username
// @Description	Get public user by username
// @Tags			user
// @Produce		json
// @Param			username	path		string							true	"username"
// @Success		200			{object}	dto.APIResponse[dto.UserGet]	"Successfully get public user by username"
// @Failure		500			{object}	dto.APIResponse[any]			"Failed to get public user by username"
// @Router			/public/user/{username} [get]
func (t *userController) GetByUsername(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	username := c.Params("username")
	user, err := t.userService.GetByUsername(ctx, username)

	if err != nil {
		t.logger.Error("failed get user by username", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to get public user by username",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[*dto.UserGet]{
		Data: user,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully get public user by username",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Get public by ticket ID
// @Description	Get public user by ticket ID
// @Tags			user
// @Produce		json
// @Param			ticketID	path		string							true	"ticket ID"
// @Success		200			{object}	dto.APIResponse[dto.UserGet]	"Successfully get public user by ticket id"
// @Failure		500			{object}	dto.APIResponse[any]			"Failed to get public user by ticket id"
// @Router			/public/user/id/{ticketID} [get]
func (t *userController) GetByTicketID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	ticketID := c.Params("ticketID")
	user, err := t.userService.GetPublicByTicketID(ctx, ticketID)

	if err != nil {
		t.logger.Error("failed get user by ticket id", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to get public user by ticket id",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[*dto.UserGet]{
		Data: user,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully get public user by ticket id",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Update user
// @Description	Update current logged-in user
// @Tags			user
// @Produce		json
// @Param			name		formData	string					true	"name"
// @Param			username	formData	string					true	"username"
// @Param			avatar		formData	file					false	"avatar"
// @Success		200			{object}	dto.APIResponse[any]	"Successfully update current user"
// @Failure		400			{object}	dto.APIResponse[any]	"Username already exists"
// @Failure		500			{object}	dto.APIResponse[any]	"Failed to update current user"
// @Router			/user/me [patch]
func (t *userController) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	userID := c.Locals("userID").(uint)
	body := c.Locals("body").(dto.UserUpdate)

	if err := t.userService.Update(ctx, c, body, userID); err != nil {
		t.logger.Error("failed to update current user", zap.Error(err))

		if helper.IsUniqueConstraintError(err) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse[any]{
				Data: nil,
				Meta: dto.ApiMeta{
					Success:   false,
					Message:   "Username already exists",
					RequestID: c.Locals("request_id").(string),
				},
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to update current user",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[any]{
		Data: nil,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully update current user",
			RequestID: c.Locals("request_id").(string),
		},
	})
}
