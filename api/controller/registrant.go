package controller

import (
	"context"
	"fybe/model/dto"
	"fybe/service"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type RegistrantController interface {
	GetByUserID(c *fiber.Ctx) error
	GetByTicketID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Generate(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
}

type registrantController struct {
	registrantService service.RegistrantService
	logger            *zap.Logger
}

func NewRegistrantService(registrantService service.RegistrantService, logger *zap.Logger) RegistrantController {
	return &registrantController{
		registrantService: registrantService,
		logger:            logger,
	}
}

// @Summary		Get by userID
// @Description	Get registrants by userID
// @Tags			registrant
// @Produce		json
// @Success		200	{object}	dto.APIResponse[[]dto.RegistrantGetByUserID]	"Failed to get registrants by user id"
// @Failure		500	{object}	dto.APIResponse[any]							"Failed to get registrants by user id"
// @Router			/registrant [get]
func (t *registrantController) GetByUserID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	userID := c.Locals("userID").(uint)

	registrants, err := t.registrantService.GetByUserID(ctx, userID)

	if err != nil {
		t.logger.Error("failed to get registrants by user id", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to get registrants by user id",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[[]dto.RegistrantGetByUserID]{
		Data: registrants,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully get registrants by user id",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Get by ticketID
// @Description	Get registrants by ticketID
// @Tags			registrant
// @Produce		json
// @Param			ticketID	path		string										true	"Ticket ID"
// @Param			page		query		int											false	"page"
// @Param			limit		query		int											false	"limit"
// @Param			search		query		string										false	"search"
// @Success		200			{object}	dto.APIResponse[dto.PaginatedRegistrants]	"Successfully get registrants by ticket id"
// @Failure		500			{object}	dto.APIResponse[any]						"Failed to get registrants by ticket id"
// @Router			/registrant/{ticketID} [get]
func (t *registrantController) GetByTicketID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	ticketID := c.Params("ticketID")
	userID := c.Locals("userID").(uint)

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	search := c.Query("search")

	result, err := t.registrantService.GetByTicketID(ctx, ticketID, userID, page, limit, search)

	if err != nil {
		t.logger.Error("failed to get registrants by ticket id", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to get registrants by ticket id",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[*dto.PaginatedRegistrants]{
		Data: result,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully get registrants by ticket id",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Create
// @Description	Create new registrant
// @Tags			registrant
// @Produce		json
// @Param			ticketID	path		string					true	"Ticket ID"
// @Param			state		query		string					false	"state"
// @Success		201			{object}	dto.APIResponse[any]	"Successfully create registrant"
// @Failure		500			{object}	dto.APIResponse[any]	"Failed to create registrant"
// @Router			/registrant/{ticketID} [post]
func (t *registrantController) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	ticketID := c.Params("ticketID")
	state := c.Query("state")
	userID := c.Locals("userID").(uint)

	if err := t.registrantService.Create(ctx, ticketID, state, userID); err != nil {
		t.logger.Error("failed to create registrant", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to create registrant",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse[any]{
		Data: nil,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully create registrant",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Verify
// @Description	Verify registrant by current user id and ticket id
// @Tags			registrant
// @Produce		json
// @Param			ticketID	path	string	true	"Ticket ID"
// @Param			id			path	string	true	"Registrant ID"
// @Success		303			{}		string	"Redirected to /verify?err=1?ticket="
// @Failure		303			{}		string	"Redirected to /verify?err=0?ticket="
// @Router			/registrant/verify/{id}/{ticketID} [get]
func (t *registrantController) Verify(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	id := c.Params("id")
	ticketID := c.Params("ticketID")
	userID := c.Locals("userID").(uint)

	if err := t.registrantService.Verify(ctx, id, ticketID, userID); err != nil {
		t.logger.Error("failed to verify registrant", zap.Error(err))
		return c.Redirect(os.Getenv("FE_DOMAIN")+"/verify?err=1&ticket="+ticketID, fiber.StatusSeeOther)

	}

	return c.Redirect(os.Getenv("FE_DOMAIN")+"/verify?err=0&ticket="+ticketID, fiber.StatusSeeOther)
}

// @Summary		Generate
// @Description	Generate unique id for private mode
// @Tags			registrant
// @Produce		json
// @Param			ticketID	path	string					true	"Ticket ID"
// @Success		201			{}		dto.APIResponse[string]	"Successfully generating unique id"
// @Failure		500			{}		dto.APIResponse[any]	"Failed generating unique id"
// @Router			/registrant/gen/{ticketID} [get]
func (t *registrantController) Generate(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	ticketID := c.Params("ticketID")
	userID := c.Locals("userID").(uint)

	id, err := t.registrantService.Generate(ctx, ticketID, userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   true,
				Message:   "Failed generating unique id",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse[string]{
		Data: id,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully generating unique id",
			RequestID: c.Locals("request_id").(string),
		},
	})
}
