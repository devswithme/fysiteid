package controller

import (
	"context"
	"fybe/model/dto"
	"fybe/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TicketController interface {
	Get(c *fiber.Ctx) error
	GetPublicByUsername(c *fiber.Ctx) error
	GetPublicByID(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type ticketController struct {
	ticketService service.TicketService
	logger        *zap.Logger
}

func NewTicketController(ticketService service.TicketService, logger *zap.Logger) TicketController {
	return &ticketController{
		ticketService: ticketService,
		logger:        logger,
	}
}

// @Summary		Get
// @Description	Get tickets
// @Tags			ticket
// @Produce		json
// @Success		200	{object}	dto.APIResponse[[]dto.TicketGet]	"Successfully get tickets"
// @Failure		500	{object}	dto.APIResponse[any]				"Failed to get tickets"
// @Router			/ticket [get]
func (t *ticketController) Get(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	userID := c.Locals("userID").(uint)

	tickets, err := t.ticketService.Get(ctx, userID)

	if err != nil {
		t.logger.Error("failed to get ticket", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to get tickets",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[[]dto.TicketGet]{
		Data: tickets,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully get tickets",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Get public by username
// @Description	Get public tickets by username
// @Tags			ticket
// @Produce		json
// @Param			username	path		string									true	"username"
// @Success		200			{object}	dto.APIResponse[[]dto.TicketPublicGet]	"Successfully get public tickets by username"
// @Failure		500			{object}	dto.APIResponse[any]					"Failed to get public tickets by username"
// @Router			/public/ticket/{username} [get]
func (t *ticketController) GetPublicByUsername(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	username := c.Params("username")

	tickets, err := t.ticketService.GetPublicByUsername(ctx, username)

	if err != nil {
		t.logger.Error("failed to get tickets by username", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to get public tickets by username",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[[]dto.TicketPublicGet]{
		Data: tickets,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully get public tickets by username",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Get public by ID
// @Description	Get public ticket by ID
// @Tags			ticket
// @Produce		json
// @Param			id	path		string							true	"id"
// @Success		200	{object}	dto.APIResponse[dto.TicketGet]	"Successfully get public ticket by id"
// @Failure		500	{object}	dto.APIResponse[any]			"Failed to get public ticket by id"
// @Router			/public/ticket/id/{id} [get]
func (t *ticketController) GetPublicByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	ticketID := c.Params("id")

	ticket, err := t.ticketService.GetPublicByID(ctx, ticketID)

	if err != nil {
		t.logger.Error("failed to get public ticket by id", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to get public ticket by id",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[*dto.TicketGet]{
		Data: ticket,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully get public ticket by id",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Get by ID
// @Description	Get ticket by ID
// @Tags			ticket
// @Produce		json
// @Param			id	path		string							true	"id"
// @Success		200	{object}	dto.APIResponse[dto.TicketGet]	"Successfully get ticket by id"
// @Failure		500	{object}	dto.APIResponse[any]			"Failed to get ticket by id"
// @Router			/ticket/{id} [get]
func (t *ticketController) GetByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	ticketID := c.Params("id")
	userID := c.Locals("userID").(uint)

	ticket, err := t.ticketService.GetByID(ctx, ticketID, userID)

	if err != nil {
		t.logger.Error("failed to get ticket by id", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to get ticket by id",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[*dto.TicketGet]{
		Data: ticket,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully get ticket by id",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Update
// @Description	Update ticket by ID
// @Tags			ticket
// @Produce		json
// @Param			id			path		string					true	"id"
// @Param			title		formData	string					true	"title"
// @Param			description	formData	string					true	"description"
// @Param			mode		formData	bool					true	"mode"
// @Param			quota		formData	uint					true	"quota"
// @Param			image		formData	file					false	"image"
// @Success		200			{object}	dto.APIResponse[any]	"Successfully update ticket"
// @Failure		500			{object}	dto.APIResponse[any]	"Failed to update ticket"
// @Router			/ticket/{id} [patch]
func (t *ticketController) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	ticketID := c.Params("id")
	userID := c.Locals("userID").(uint)
	body := c.Locals("body").(dto.TicketCreate)

	if err := t.ticketService.Update(ctx, c, body, ticketID, userID); err != nil {
		t.logger.Error("failed to update ticket", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to update ticket",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[any]{
		Data: nil,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully update ticket",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Create
// @Description	Create new ticket
// @Tags			ticket
// @Produce		json
// @Param			title		formData	string					true	"title"
// @Param			description	formData	string					true	"description"
// @Param			mode		formData	bool					true	"mode"
// @Param			quota		formData	uint					true	"quota"
// @Param			image		formData	file					false	"image"
// @Success		201			{object}	dto.APIResponse[any]	"Successfully create ticket"
// @Failure		500			{object}	dto.APIResponse[any]	"Failed to create ticket"
// @Router			/ticket [post]
func (t *ticketController) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	userID := c.Locals("userID").(uint)
	body := c.Locals("body").(dto.TicketCreate)

	if err := t.ticketService.Create(ctx, c, body, userID); err != nil {
		t.logger.Error("failed to create ticket", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to create ticket",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse[any]{
		Data: nil,
		Meta: dto.ApiMeta{
			Success:   false,
			Message:   "Successfully create ticket",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Delete
// @Description	Delete ticket by ID
// @Tags			ticket
// @Produce		json
// @Param			id	path		string					true	"id"
// @Success		200	{object}	dto.APIResponse[any]	"Successfully delete ticket"
// @Failure		500	{object}	dto.APIResponse[any]	"Failed to delete ticket"
// @Router			/ticket/{id} [delete]
func (t *ticketController) Delete(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	ticketID := c.Params("id")
	userID := c.Locals("userID").(uint)

	if err := t.ticketService.Delete(ctx, ticketID, userID); err != nil {
		t.logger.Error("failed to delete ticket", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to delete ticket",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[any]{
		Data: nil,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully delete ticket",
			RequestID: c.Locals("request_id").(string),
		},
	})
}
