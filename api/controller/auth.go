package controller

import (
	"fybe/config"
	"fybe/helper"
	"fybe/model/dto"
	"fybe/service"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthController interface {
	Login(c *fiber.Ctx) error
	Callback(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type authController struct {
	service service.AuthService
	logger  *zap.Logger
}

func NewAuthController(
	service service.AuthService,
	logger *zap.Logger,
) AuthController {
	return &authController{
		service: service,
		logger:  logger,
	}
}

// @Summary		Google OAuth2 Login
// @Description	Oauth2 Google login entry
// @Tags			login
// @Param			redirect	query		string					false	"Redirect URL after login"	default(/dash)
// @Success		303			{}			string					"Redirected to Google login page"
// @Failure		500			{object}	dto.APIResponse[any]	"Failed to login"
// @Router			/login/google [get]
func (t *authController) Login(c *fiber.Ctx) error {
	redirect := c.Query("redirect")
	if redirect == "" {
		redirect = "/dash"
	}

	u, err := url.Parse(redirect)

	if err != nil || u.Host != "" || u.IsAbs() || !strings.HasPrefix(u.Path, "/") {
		redirect = "/dash"
	}

	log.Println(u)

	state, err := t.service.SaveState(redirect)

	if err != nil {
		t.logger.Error("failed to save state", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to login",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Redirect(config.GoogleOAuthConfig.AuthCodeURL(state), fiber.StatusSeeOther)
}

// @Summary		Google OAuth2 Callback
// @Description	Route to handle Google callback
// @Tags			login
// @Success		303	{}			string					"Redirected to frontend with desired redirect URL path"
// @Failure		500	{object}	dto.APIResponse[any]	"Failed to login"
// @Router			/login/google/callback [get]
func (t *authController) Callback(c *fiber.Ctx) error {
	redirect, err := t.service.Callback(c)

	if err != nil {
		t.logger.Error("failed to handle callback", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to login",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Redirect(os.Getenv("FE_DOMAIN")+redirect, fiber.StatusSeeOther)
}

// @Summary		Refresh token
// @Description	Route to generate new access token
// @Tags			auth
// @Success		201	{object}	dto.APIResponse[any]	"Successfully regenerate token"
// @Failure		401	{object}	dto.APIResponse[any]	"Failed to regenerate token"
// @Router			/auth/refresh [post]
func (t *authController) Refresh(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	if err := t.service.Refresh(c, userID); err != nil {
		t.logger.Error("failed to refresh token", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse[any]{
			Data: nil,
			Meta: dto.ApiMeta{
				Success:   false,
				Message:   "Failed to regenerate token",
				RequestID: c.Locals("request_id").(string),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse[any]{
		Data: nil,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully regenerate token",
			RequestID: c.Locals("request_id").(string),
		},
	})
}

// @Summary		Logout
// @Description	Remove token from cookies
// @Tags			auth
// @Success		200	{object}	dto.APIResponse[any]	"Successfully logout"
// @Router			/auth/logout [post]
func (t *authController) Logout(c *fiber.Ctx) error {
	helper.NewCookieHelper(c, "access_token").Delete()
	helper.NewCookieHelper(c, "refresh_token").Delete()

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse[any]{
		Data: nil,
		Meta: dto.ApiMeta{
			Success:   true,
			Message:   "Successfully logout",
			RequestID: c.Locals("request_id").(string),
		},
	})
}
